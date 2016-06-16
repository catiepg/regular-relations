package relations

// TODO: final
// TODO: case <a,b> and <a,c>
type tState struct {
	positions *set
	marked    bool
	next      map[string][]*tState
	out       map[string][]string
}

func newTState(positions *set, marked bool) *tState {
	return &tState{positions: positions, marked: marked,
		next: make(map[string][]*tState), out: make(map[string][]string)}
}

type tStates []*tState

func (ss *tStates) getUnmarked() *tState {
	for _, state := range *ss {
		if state.marked == false {
			return state
		}
	}
	return nil
}

func (ss *tStates) hasWithPositions(positions *set) *tState {
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
	newTransducer := &transducer{}
	var states tStates

	start := newTState(t.rootFirst, false)
	newTransducer.start = start
	states = append(states, start)

	for {
		// check for unmarked
		state := states.getUnmarked()

		// no unmarked found
		if state == nil {
			break
		}

		state.marked = true

		for _, symb := range t.alphabet {
			u := newSet()
			for position := range *state.positions {
				if t.symbols[position].equal(symb) {
					u = u.union(t.follow[position])
				}
			}

			newState := states.hasWithPositions(u)
			if u.cardinality() != 0 && newState == nil {
				newState = newTState(u, false)
				states = append(states, newState)
			}

			state.next[symb.in] = append(state.next[symb.in], newState)
			state.out[symb.in] = append(state.out[symb.in], symb.out)
		}
	}

	return newTransducer
}
