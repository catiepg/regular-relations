package relations

import "io"

// TODO: handle epsilon and empty set as input - 0 and 1
type tTransition struct {
	state *tState
	out   string
}

type tState struct {
	index int
	next  map[rune][]*tTransition
	final bool
}

type tStates struct {
	all       map[int]*tState // index -> state
	positions map[int]*set    // index -> positions
	marked    []int
	unmarked  []int
	index     int
}

func (ss *tStates) add(positions *set) *tState {
	ss.index += 1
	state := &tState{
		index: ss.index,
		next:  make(map[rune][]*tTransition),
	}
	ss.all[ss.index] = state
	ss.positions[ss.index] = positions
	ss.unmarked = append(ss.unmarked, ss.index)
	return state
}

func (ss *tStates) get(positions *set) *tState {
	for i, p := range ss.positions {
		if p.equal(positions) {
			return ss.all[i]
		}
	}
	return nil
}

type transducer struct {
	start *tState
}

func NewTransducer(source io.Reader) (*transducer, error) {
	tree, err := NewTree(source)
	if err != nil {
		return nil, err
	}

	states := tStates{all: make(map[int]*tState), positions: make(map[int]*set)}
	start := states.add(tree.rootFirst)

	for {
		if len(states.unmarked) == 0 {
			break
		}

		// Get next unmarked state and add it to the set of marked states
		stateIndex := states.unmarked[0]
		states.unmarked = states.unmarked[1:]
		states.marked = append(states.marked, stateIndex)
		state := states.all[stateIndex]

		for _, symb := range tree.alphabet {
			// Get union of follow for all positions with current symbol
			u := newSet()
			for position := range *states.positions[stateIndex] {
				if tree.elements[position].contain(symb.in, symb.out) {
					u = u.union(tree.follow[position])
				}
			}

			// Check if state with these positions already exists...
			nextState := states.get(u)

			// ...otherwise create new state
			if u.cardinality() != 0 && nextState == nil {
				nextState = states.add(u)
				if u.contains(tree.final) {
					nextState.final = true
				}
			}

			// Add transitions
			if nextState != nil {
				o := &tTransition{state: nextState, out: symb.out}
				state.next[symb.in] = append(state.next[symb.in], o)
			}
		}
	}

	return &transducer{start: start}, nil
}
