package relations

import (
	"bytes"
	"io"
	"sort"

	"github.com/oleiade/lane"
	"github.com/s2gatev/hcache"
)

// pair is used in the construction of the subsequential transducer.
// It contains a transducer state and the symbols that remain to be
// added to the output.
type pair struct {
	state     *tState
	remaining string
}

func (p *pair) equal(o *pair) bool {
	return p.state.index == o.state.index && p.remaining == o.remaining
}

// pairs is a slice compatible with the hcache structure.
// The underlying values are of type *pair.
type pairs []hcache.Key

func (ps pairs) Len() int {
	return len(ps)
}

func (ps pairs) Less(i, j int) bool {
	p1 := ps[i].(*pair)
	p2 := ps[j].(*pair)

	if p1.state.index == p2.state.index {
		return p1.remaining < p2.remaining
	}

	return p1.state.index < p2.state.index
}

func (ps pairs) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// sState is a state in a subsequential transducer.
type sState struct {
	remainingPairs pairs
	next           map[rune]*sState
	out            map[rune]string
	final          bool
	finalOut       []string
	isVisited      bool
}

func newSState() *sState {
	return &sState{next: make(map[rune]*sState), out: make(map[rune]string)}
}

// getFinalOut gets final outputs if pair has final state.
func (ss *sState) getFinalOut() []string {
	var finalRemaining []string
	for _, p := range ss.remainingPairs {
		p := p.(*pair)
		if p.state.final {
			finalRemaining = append(finalRemaining, p.remaining)
		}
	}
	return finalRemaining
}

// lcp calculates the longest common prefix of the input strings.
func lcp(strs [][]rune) string {
	if len(strs) == 0 {
		return ""
	}

	for i, c := range strs[0] {
		for _, s := range strs[1:] {
			if i == len(s) || s[i] != c {
				return string(strs[0][:i])
			}
		}
	}
	return string(strs[0])
}

// RegularRelation is a struct containing the initial state of the
// subsequential transducer that recognizes the input regular relation.
type RegularRelation struct {
	start *sState
}

// Transduce feeds the input string into the RegularRelation transducer
// and returns all possible results from the output transducer tape.
func (s *RegularRelation) Transduce(input string) ([]string, bool) {
	node := s.start
	var output string

	for _, symbol := range input {
		nextnode, ok := node.next[symbol]
		if !ok {
			return nil, false
		}
		output += node.out[symbol]
		node = nextnode
	}

	if !node.final {
		return nil, false
	}

	var result []string
	for _, o := range node.finalOut {
		result = append(result, output+o)
	}

	return result, true
}

// Build builds a RegularRelation subsequential transducer from the
// input regular relation expression.
// TODO: define regular relation grammar.
func Build(source io.Reader) (*RegularRelation, error) {
	tr, err := newTransducer(source)
	if err != nil {
		return nil, err
	}

	stateQueue := lane.NewQueue()
	sc := hcache.New()

	initPair := &pair{state: tr.root, remaining: ""}
	start := sc.GetOrInsert(newSState(), initPair).(*sState)
	start.remainingPairs = append(start.remainingPairs, initPair)
	stateQueue.Enqueue(start)

	for stateQueue.Size() != 0 {
		state := stateQueue.Dequeue().(*sState)

		// Check if state should be final and add outputs to final output.
		if final := state.getFinalOut(); len(final) != 0 {
			state.final = true
			state.finalOut = append(state.finalOut, final...)
		}

		// Get groups of pairs that have states with same input symbol.
		withInput := make(map[rune]pairs)
		for _, p := range state.remainingPairs {
			p := p.(*pair)
			for in := range p.state.next {
				withInput[in] = append(withInput[in], p)
			}
		}

		for in, ps := range withInput {
			// Get all remaining+out strings from states with given input
			// and map them to corresponding next state.
			var outputs [][]rune
			nextStates := make(map[int]*tState)
			for _, p := range ps {
				p := p.(*pair)
				remaining := bytes.Runes([]byte(p.remaining))
				for _, o := range p.state.next[in] {
					out := append(remaining, bytes.Runes([]byte(o.out))...)
					outputs = append(outputs, out)
					nextStates[len(outputs)-1] = o.state
				}
			}

			// Calculate longest common prefix.
			state.out[in] = lcp(outputs)

			// Create new pairs by removing the longest common prefix from
			// the outputs.
			var newPairs pairs
			for i, out := range outputs {
				newPairs = append(newPairs, &pair{
					state:     nextStates[i],
					remaining: string(out[len(state.out[in]):]),
				})
			}
			sort.Sort(newPairs)

			// Check if state with such state pairs exists...
			nextState := sc.GetOrInsert(newSState(), newPairs...).(*sState)

			// ...and populate the state with the new pairs if necessary.
			if !nextState.isVisited {
				nextState.isVisited = true
				nextState.remainingPairs = newPairs
				stateQueue.Enqueue(nextState)
			}

			state.next[in] = nextState
		}
	}

	return &RegularRelation{start: start}, nil
}
