package relations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstPos(t *testing.T) {
	rootFirst := buildTree(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`).rootFirst

	assert.True(t, rootFirst.equal(newSet(1, 2, 3)))
}

func TestFollowPos(t *testing.T) {
	follow := buildTree(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`).follow

	assert.True(t, follow[1].equal(newSet(1, 2, 3)))
	assert.True(t, follow[2].equal(newSet(1, 2, 3)))
	assert.True(t, follow[3].equal(newSet(4)))
	assert.True(t, follow[4].equal(newSet(5)))
	assert.True(t, follow[5].equal(newSet(6)))

	_, isSet := follow[6]
	assert.False(t, isSet)
}

func TestAlphabet(t *testing.T) {
	alphabet := buildTree(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`).alphabet

	assert.Equal(t, 2, len(alphabet))
}

func TestSymbols(t *testing.T) {
	symbols := buildTree(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`).symbols

	assert.Equal(t, 6, len(symbols))
	assert.True(t, symbols[1].equal(&pair{in: "a", out: ""}))
	assert.True(t, symbols[2].equal(&pair{in: "b", out: ""}))
	assert.True(t, symbols[6].equal(&pair{in: "!", out: ""}))
}

func TestFinal(t *testing.T) {
	final := buildTree(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`).final

	assert.Equal(t, 6, final)
}

func TestMulticharFollow(t *testing.T) {
	follow := buildTree(`<abc,xy>+<bca,zz>`).follow

	assert.True(t, follow[1].equal(newSet(2)))
	assert.True(t, follow[2].equal(newSet(3)))
	assert.True(t, follow[3].equal(newSet(7)))
	assert.True(t, follow[4].equal(newSet(5)))
	assert.True(t, follow[5].equal(newSet(6)))
	assert.True(t, follow[6].equal(newSet(7)))

	_, isSet := follow[7]
	assert.False(t, isSet)
}

func TestMulticharRootFirst(t *testing.T) {
	rootFirst := buildTree(`<abc,xy>+<bca,zz>`).rootFirst

	assert.True(t, rootFirst.equal(newSet(1, 4)))
}

func TestMulticharAlphabet(t *testing.T) {
	alphabet := buildTree(`<abc,xy>+<bca,zz>`).alphabet

	expected := []*pair{
		&pair{"a", "xy"},
		&pair{"b", ""},
		&pair{"c", ""},
		&pair{"b", "zz"},
		&pair{"a", ""},
	}

	for i, charPair := range expected {
		assert.True(t, charPair.equal(alphabet[i]))
	}
}
