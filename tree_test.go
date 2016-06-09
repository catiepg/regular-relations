package relations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO: reverse arguments
func TestParseTreeBuilder(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root := buildTree(regex).root

	assert.Equal(t, root.data.operator, ".")
	assert.Equal(t, root.left.data.operator, "*")
	assert.Equal(t, root.right.data.operator, ".")

	assert.Equal(t, root.left.left.data.operator, "+")
	assert.Equal(t, root.left.left.left.data.in, "a")
	assert.Equal(t, root.left.left.left.data.out, "")
	assert.Equal(t, root.left.left.right.data.in, "b")
	assert.Equal(t, root.left.left.right.data.out, "")

	assert.Equal(t, root.right.left.data.in, "a")
	assert.Equal(t, root.right.left.data.out, "")
	assert.Equal(t, root.right.right.data.operator, ".")

	assert.Equal(t, root.right.right.left.data.in, "b")
	assert.Equal(t, root.right.right.left.data.out, "")
	assert.Equal(t, root.right.right.right.data.operator, ".")

	assert.Equal(t, root.right.right.right.left.data.in, "b")
	assert.Equal(t, root.right.right.right.left.data.out, "")
	assert.Equal(t, root.right.right.right.right.data.in, "!")
}

func TestParseTreeIndexing(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root := buildTree(regex).root

	assert.Equal(t, root.left.left.left.index, 1)
	assert.Equal(t, root.left.left.right.index, 2)
	assert.Equal(t, root.right.left.index, 3)
	assert.Equal(t, root.right.right.left.index, 4)
	assert.Equal(t, root.right.right.right.left.index, 5)
	assert.Equal(t, root.right.right.right.right.index, 6)
}

func TestTreeNodeNullable(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root := buildTree(regex).root

	assert.True(t, root.left.nullable)
	assert.False(t, root.nullable)
	assert.False(t, root.right.nullable)
}

func TestFirstPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root := buildTree(regex).root

	assert.True(t, root.first.equal(newSet(1, 2, 3)))
	assert.True(t, root.left.first.equal(newSet(1, 2)))
	assert.True(t, root.right.first.equal(newSet(3)))

	assert.True(t, root.left.left.first.equal(newSet(1, 2)))
	assert.True(t, root.left.left.left.first.equal(newSet(1)))
	assert.True(t, root.left.left.right.first.equal(newSet(2)))

	assert.True(t, root.right.left.first.equal(newSet(3)))
	assert.True(t, root.right.right.first.equal(newSet(4)))

	assert.True(t, root.right.right.left.first.equal(newSet(4)))
	assert.True(t, root.right.right.right.first.equal(newSet(5)))

	assert.True(t, root.right.right.right.left.first.equal(newSet(5)))
	assert.True(t, root.right.right.right.right.first.equal(newSet(6)))
}

func TestLastPos(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root := buildTree(regex).root

	assert.True(t, root.last.equal(newSet(6)))
	assert.True(t, root.left.last.equal(newSet(1, 2)))
	assert.True(t, root.right.last.equal(newSet(6)))

	assert.True(t, root.left.left.last.equal(newSet(1, 2)))
	assert.True(t, root.left.left.left.last.equal(newSet(1)))
	assert.True(t, root.left.left.right.last.equal(newSet(2)))

	assert.True(t, root.right.left.last.equal(newSet(3)))
	assert.True(t, root.right.right.last.equal(newSet(6)))

	assert.True(t, root.right.right.left.last.equal(newSet(4)))
	assert.True(t, root.right.right.right.last.equal(newSet(6)))

	assert.True(t, root.right.right.right.left.last.equal(newSet(5)))
	assert.True(t, root.right.right.right.right.last.equal(newSet(6)))
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
