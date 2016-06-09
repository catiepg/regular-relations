package relations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	rootFirst := buildTree(regex).rootFirst

	assert.True(t, rootFirst.equal(newSet(1, 2, 3)))
}

func TestFollowPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	follow := buildTree(regex).follow

	assert.True(t, follow[1].equal(newSet(1, 2, 3)))
	assert.True(t, follow[2].equal(newSet(1, 2, 3)))
	assert.True(t, follow[3].equal(newSet(4)))
	assert.True(t, follow[4].equal(newSet(5)))
	assert.True(t, follow[5].equal(newSet(6)))

	_, isSet := follow[6]
	assert.False(t, isSet)
}

func TestAlphabet(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	alphabet := buildTree(regex).alphabet

	assert.Equal(t, 2, len(alphabet))
}

func TestSymbols(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	symbols := buildTree(regex).symbols

	assert.Equal(t, 6, len(symbols))
	assert.True(t, symbols[1].equal(&pair{in: "a", out: ""}))
	assert.True(t, symbols[2].equal(&pair{in: "b", out: ""}))
	assert.True(t, symbols[6].equal(&pair{in: "!", out: ""}))
}
