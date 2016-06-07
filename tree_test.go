package relations

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseTreeBuilder(t *testing.T) {
	regex := `(<a,>+<b,>)*.<a,>.<b,>.<b,>`
	root := parseTree(regex)

	assert.Equal(t, root.data, ".")
	assert.Equal(t, root.left.data, "*")
	assert.Equal(t, root.right.data, ".")

	assert.Equal(t, root.left.left.data, "+")
	assert.Equal(t, root.left.left.left.data, "<a,>")
	assert.Equal(t, root.left.left.right.data, "<b,>")

	assert.Equal(t, root.right.left.data, "<a,>")
	assert.Equal(t, root.right.right.data, ".")

	assert.Equal(t, root.right.right.left.data, "<b,>")
	assert.Equal(t, root.right.right.right.data, "<b,>")
}
