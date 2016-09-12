package relations

type set map[int]struct{}

func newSet(elements ...int) set {
	s := make(set)
	for _, e := range elements {
		s.add(e)
	}
	return s
}

func (s set) add(e int) {
	s[e] = struct{}{}
}

func (s set) cardinality() int {
	return len(s)
}

func (s set) contains(e int) bool {
	_, ok := s[e]
	return ok
}

func (s set) union(other set) set {
	unionSet := s.clone()
	if other != nil {
		for e := range other {
			unionSet.add(e)
		}
	}
	return unionSet
}

func (s set) clone() set {
	newClone := newSet()
	for e := range s {
		newClone.add(e)
	}
	return newClone
}

func (s set) equal(other set) bool {
	if s.cardinality() != other.cardinality() {
		return false
	}

	for e := range s {
		if !other.contains(e) {
			return false
		}
	}

	return true
}
