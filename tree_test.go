package relations

import (
	"testing"

	"github.com/deckarep/golang-set"
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

	assert.True(t, root.firstPos.Equal(mapset.NewSetWith(1, 2, 3)))
	assert.True(t, root.left.firstPos.Equal(mapset.NewSetWith(1, 2)))
	assert.True(t, root.right.firstPos.Equal(mapset.NewSetWith(3)))

	assert.True(t, root.left.left.firstPos.Equal(mapset.NewSetWith(1, 2)))
	assert.True(t, root.left.left.left.firstPos.Equal(mapset.NewSetWith(1)))
	assert.True(t, root.left.left.right.firstPos.Equal(mapset.NewSetWith(2)))

	assert.True(t, root.right.left.firstPos.Equal(mapset.NewSetWith(3)))
	assert.True(t, root.right.right.firstPos.Equal(mapset.NewSetWith(4)))

	assert.True(t, root.right.right.left.firstPos.Equal(mapset.NewSetWith(4)))
	assert.True(t, root.right.right.right.firstPos.Equal(mapset.NewSetWith(5)))

	assert.True(t, root.right.right.right.left.firstPos.Equal(mapset.NewSetWith(5)))
	assert.True(t, root.right.right.right.right.firstPos.Equal(mapset.NewSetWith(6)))
}

func TestLastPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root, _ := parseTree(regex)

	assert.True(t, root.lastPos.Equal(mapset.NewSetWith(6)))
	assert.True(t, root.left.lastPos.Equal(mapset.NewSetWith(1, 2)))
	assert.True(t, root.right.lastPos.Equal(mapset.NewSetWith(6)))

	assert.True(t, root.left.left.lastPos.Equal(mapset.NewSetWith(1, 2)))
	assert.True(t, root.left.left.left.lastPos.Equal(mapset.NewSetWith(1)))
	assert.True(t, root.left.left.right.lastPos.Equal(mapset.NewSetWith(2)))

	assert.True(t, root.right.left.lastPos.Equal(mapset.NewSetWith(3)))
	assert.True(t, root.right.right.lastPos.Equal(mapset.NewSetWith(6)))

	assert.True(t, root.right.right.left.lastPos.Equal(mapset.NewSetWith(4)))
	assert.True(t, root.right.right.right.lastPos.Equal(mapset.NewSetWith(6)))

	assert.True(t, root.right.right.right.left.lastPos.Equal(mapset.NewSetWith(5)))
	assert.True(t, root.right.right.right.right.lastPos.Equal(mapset.NewSetWith(6)))
}

func TestFollowPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	_, followPos := parseTree(regex)

	assert.True(t, followPos[1].Equal(mapset.NewSetWith(1, 2, 3)))
	assert.True(t, followPos[2].Equal(mapset.NewSetWith(1, 2, 3)))
	assert.True(t, followPos[3].Equal(mapset.NewSetWith(4)))
	assert.True(t, followPos[4].Equal(mapset.NewSetWith(5)))
	assert.True(t, followPos[5].Equal(mapset.NewSetWith(6)))

	_, isSet := followPos[6]
	assert.False(t, isSet)
}
