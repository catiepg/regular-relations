package relations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSet(t *testing.T) {
	s := newSet(1, 2 ,3)

	assert.Equal(t, s.cardinality(), 3)
	assert.True(t, s.contains(1))
	assert.True(t, s.contains(2))
	assert.True(t, s.contains(3))
}

func TestNoDuplicates(t *testing.T) {
	s := newSet(1, 2 ,3)
	assert.Equal(t, s.cardinality(), 3)

	s.add(3)
	assert.Equal(t, s.cardinality(), 3)
}

func TestUnion(t *testing.T) {
	s := newSet(1, 2)
	o := newSet(2, 3)
	u := s.union(o)

	assert.Equal(t, u.cardinality(), 3)
	assert.True(t, u.contains(1))
	assert.True(t, u.contains(2))
	assert.True(t, u.contains(3))
}

func TestEqual(t *testing.T) {
	s := newSet(1, 2)
	o := newSet(1, 2)
	assert.True(t, s.equal(o))

	s.add(3)
	assert.False(t, s.equal(o))
}
