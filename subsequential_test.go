package relations

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestSubsequential(t *testing.T) {
	s := buildSubsequential(`<abc,xyz>+<acc,qwe>`)

	out, ok := s.get("abc")
	assert.Equal(t, 1, len(out))
	assert.Equal(t, "xyz", out[0])
	assert.True(t, ok)
}

func TestSubsequentialStructure(t *testing.T) {
	s := buildSubsequential(`<abc,xyz>+<acc,qwe>`)

	assert.Equal(t, 1, len(s.start.next))
	assert.Equal(t, "", s.start.out["a"])

	state2 := s.start.next["a"]
	assert.Equal(t, 2, len(state2.next))
	assert.Equal(t, "xyz", state2.out["b"])
	assert.Equal(t, "qwe", state2.out["c"])

	state3 := state2.next["b"]
	state4 := state2.next["c"]
	assert.Equal(t, 1, len(state3.next))
	assert.Equal(t, 1, len(state4.next))
	assert.True(t, state3.next["c"] == state4.next["c"])
}
