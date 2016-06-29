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
	all       map[int]*tState // state index -> state
	positions map[int]*set    // state index -> positions
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
		state := states.all[stateIndex]

		// Get union of follow for positions in the state than correspond
		// to the same element, instead of going through each element in the
		// alphabet
		followUnion := make(map[*element]*set)
		for position := range *states.positions[stateIndex] {
			elem := tree.elements[position]
			if _, ok := followUnion[elem]; ok {
				for p := range *tree.follow[position] {
					followUnion[elem].add(p)
				}
			} else {
				if tree.follow[position] != nil {
					followUnion[elem] = tree.follow[position].clone()
				}
			}
		}

		for symb, union := range followUnion {
			// Check if state with these positions already exists...
			nextState := states.get(union)

			// ...otherwise create new state
			if nextState == nil {
				nextState = states.add(union)
				if union.contains(tree.finalIndex) {
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
