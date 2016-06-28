package relations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstPos(t *testing.T) {
	source := strings.NewReader(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)
	tree, _ := NewTree(source)

	assert.True(t, tree.rootFirst.equal(newSet(1, 2, 3)))
}

func TestFollowPos(t *testing.T) {
	source := strings.NewReader(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)
	tree, _ := NewTree(source)

	assert.True(t, tree.follow[1].equal(newSet(1, 2, 3)))
	assert.True(t, tree.follow[2].equal(newSet(1, 2, 3)))
	assert.True(t, tree.follow[3].equal(newSet(4)))
	assert.True(t, tree.follow[4].equal(newSet(5)))
	assert.True(t, tree.follow[5].equal(newSet(6)))

	_, isSet := tree.follow[6]
	assert.False(t, isSet)
}

func TestAlphabet(t *testing.T) {
	source := strings.NewReader(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)
	tree, _ := NewTree(source)

	assert.Equal(t, 2, len(tree.alphabet))
}

func TestSymbols(t *testing.T) {
	source := strings.NewReader(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)
	tree, _ := NewTree(source)

	assert.Equal(t, 6, len(tree.elements))
	assert.True(t, tree.elements[1].contain('a', ""))
	assert.True(t, tree.elements[2].contain('b', ""))
	assert.True(t, tree.elements[6].contain('!', ""))
}

func TestFinal(t *testing.T) {
	source := strings.NewReader(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)
	tree, _ := NewTree(source)

	assert.Equal(t, 6, tree.final)
}

func TestMulticharFollow(t *testing.T) {
	source := strings.NewReader(`<abc,xy>+<bca,zz>`)
	tree, _ := NewTree(source)

	assert.True(t, tree.follow[1].equal(newSet(2)))
	assert.True(t, tree.follow[2].equal(newSet(3)))
	assert.True(t, tree.follow[3].equal(newSet(7)))
	assert.True(t, tree.follow[4].equal(newSet(5)))
	assert.True(t, tree.follow[5].equal(newSet(6)))
	assert.True(t, tree.follow[6].equal(newSet(7)))

	_, isSet := tree.follow[7]
	assert.False(t, isSet)
}

func TestMulticharRootFirst(t *testing.T) {
	source := strings.NewReader(`<abc,xy>+<bca,zz>`)
	tree, _ := NewTree(source)

	assert.True(t, tree.rootFirst.equal(newSet(1, 4)))
}

func TestMulticharAlphabet(t *testing.T) {
	source := strings.NewReader(`<abc,xy>+<bca,zz>`)
	tree, _ := NewTree(source)

	expected := []struct {
		in  rune
		out string
	}{
		{'a', "xy"}, {'b', ""}, {'c', ""}, {'b', "zz"}, {'a', ""},
	}

	for i, o := range expected {
		assert.True(t, tree.alphabet[i].contain(o.in, o.out))
	}
}
