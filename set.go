package relations

type set map[int]struct{}

func newSet(elements ...int) *set {
	s := make(set)
	for _, element := range elements {
		s.add(element)
	}
	return &s
}

func (s *set) add(element int) {
	(*s)[element] = struct{}{}
}

func (s *set) cardinality() int {
	return len(*s)
}

func (s *set) contains(element int) bool {
	_, ok := (*s)[element]
	return ok
}

func (s *set) union(other *set) *set {
	unionSet := s.clone()
	if other != nil {
		for element := range (*other) {
			unionSet.add(element)
		}
	}
	return unionSet
}

func (s *set) clone() *set {
	newClone := newSet()
	for element := range (*s) {
		newClone.add(element)
	}
	return newClone
}

func (s *set) equal(other *set) bool {
	if s.cardinality() != other.cardinality() {
		return false
	}

	for element := range (*s) {
		if !other.contains(element) {
			return false
		}
	}

	return true
}
