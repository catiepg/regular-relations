package relations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTreeBuilder(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root, _ := parseTree(regex)

	assert.Equal(t, root.data, ".")
	assert.Equal(t, root.left.data, "*")
	assert.Equal(t, root.right.data, ".")

	assert.Equal(t, root.left.left.data, "+")
	assert.Equal(t, root.left.left.left.data, "<a,>")
	assert.Equal(t, root.left.left.right.data, "<b,>")

	assert.Equal(t, root.right.left.data, "<a,>")
	assert.Equal(t, root.right.right.data, ".")

	assert.Equal(t, root.right.right.left.data, "<b,>")
	assert.Equal(t, root.right.right.right.data, ".")

	assert.Equal(t, root.right.right.right.left.data, "<b,>")
	assert.Equal(t, root.right.right.right.right.data, "!")
}

func TestParseTreeIndexing(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root, _ := parseTree(regex)

	assert.Equal(t, root.left.left.left.index, 1)
	assert.Equal(t, root.left.left.right.index, 2)
	assert.Equal(t, root.right.left.index, 3)
	assert.Equal(t, root.right.right.left.index, 4)
	assert.Equal(t, root.right.right.right.left.index, 5)
	assert.Equal(t, root.right.right.right.right.index, 6)
}

func TestTreeNodeNullable(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root, _ := parseTree(regex)

	assert.True(t, root.left.nullable)
	assert.False(t, root.nullable)
	assert.False(t, root.right.nullable)
}

func TestFirstPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root, _ := parseTree(regex)

	assert.True(t, root.firstPos.equal(newSet(1, 2, 3)))
	assert.True(t, root.left.firstPos.equal(newSet(1, 2)))
	assert.True(t, root.right.firstPos.equal(newSet(3)))

	assert.True(t, root.left.left.firstPos.equal(newSet(1, 2)))
	assert.True(t, root.left.left.left.firstPos.equal(newSet(1)))
	assert.True(t, root.left.left.right.firstPos.equal(newSet(2)))

	assert.True(t, root.right.left.firstPos.equal(newSet(3)))
	assert.True(t, root.right.right.firstPos.equal(newSet(4)))

	assert.True(t, root.right.right.left.firstPos.equal(newSet(4)))
	assert.True(t, root.right.right.right.firstPos.equal(newSet(5)))

	assert.True(t, root.right.right.right.left.firstPos.equal(newSet(5)))
	assert.True(t, root.right.right.right.right.firstPos.equal(newSet(6)))
}

func TestLastPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root, _ := parseTree(regex)

	assert.True(t, root.lastPos.equal(newSet(6)))
	assert.True(t, root.left.lastPos.equal(newSet(1, 2)))
	assert.True(t, root.right.lastPos.equal(newSet(6)))

	assert.True(t, root.left.left.lastPos.equal(newSet(1, 2)))
	assert.True(t, root.left.left.left.lastPos.equal(newSet(1)))
	assert.True(t, root.left.left.right.lastPos.equal(newSet(2)))

	assert.True(t, root.right.left.lastPos.equal(newSet(3)))
	assert.True(t, root.right.right.lastPos.equal(newSet(6)))

	assert.True(t, root.right.right.left.lastPos.equal(newSet(4)))
	assert.True(t, root.right.right.right.lastPos.equal(newSet(6)))

	assert.True(t, root.right.right.right.left.lastPos.equal(newSet(5)))
	assert.True(t, root.right.right.right.right.lastPos.equal(newSet(6)))
}

func TestFollowPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	_, followPos := parseTree(regex)

	assert.True(t, followPos[1].equal(newSet(1, 2, 3)))
	assert.True(t, followPos[2].equal(newSet(1, 2, 3)))
	assert.True(t, followPos[3].equal(newSet(4)))
	assert.True(t, followPos[4].equal(newSet(5)))
	assert.True(t, followPos[5].equal(newSet(6)))

	_, isSet := followPos[6]
	assert.False(t, isSet)
}
