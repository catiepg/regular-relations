package relations

type tState struct {
	positions *set
	next      map[string][]*tState
	out       map[string][]string
	final     bool
}

func newTState(positions *set) *tState {
	return &tState{
		positions: positions,
		next:      make(map[string][]*tState),
		out:       make(map[string][]string),
	}
}

type tStates []*tState

func (ss *tStates) get(positions *set) *tState {
	for _, state := range *ss {
		if state.positions.equal(positions) {
			return state
		}
	}
	return nil
}

type transducer struct {
	start *tState
}

func buildTransducer(raw string) *transducer {
	t := buildTree(raw)
	start := newTState(t.rootFirst)

	newTransducer := &transducer{}
	newTransducer.start = start

	var states, unmarkedStates tStates
	unmarkedStates = append(unmarkedStates, start)

	for {
		if len(unmarkedStates) == 0 {
			break
		}

		// Get next unmarked state and add it to the set of marked states
		state := unmarkedStates[0]
		unmarkedStates = unmarkedStates[1:]
		states = append(states, state)

		for _, symb := range t.alphabet {
			// Get union of follow for all positions with current symbol
			u := newSet()
			for position := range *state.positions {
				if t.symbols[position].equal(symb) {
					u = u.union(t.follow[position])
				}
			}

			// Check if state with these positions already exists,
			newState := states.get(u)
			if newState == nil {
				newState = unmarkedStates.get(u)
			}

			// otherwise create new
			if u.cardinality() != 0 && newState == nil {
				newState = newTState(u)
				if newState.positions.contains(t.final) {
					newState.final = true
				}
				unmarkedStates = append(unmarkedStates, newState)
			}

			// Add transitions
			state.next[symb.in] = append(state.next[symb.in], newState)
			state.out[symb.in] = append(state.out[symb.in], symb.out)
		}
	}

	return newTransducer
}
