package relations

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestTransducer(t *testing.T) {
	tr := buildTransducer(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)

	a := tr.start
	assert.True(t, a.positions.equal(newSet(1, 2, 3)))
	assert.True(t, a.next["a"][0].positions.equal(newSet(1, 2, 3, 4)))
	assert.True(t, a.next["b"][0].positions.equal(newSet(1, 2, 3)))

	b := a.next["a"][0]
	assert.True(t, b.next["a"][0].positions.equal(newSet(1, 2, 3, 4)))
	assert.True(t, b.next["b"][0].positions.equal(newSet(1, 2, 3, 5)))

	c := b.next["b"][0]
	assert.True(t, c.next["a"][0].positions.equal(newSet(1, 2, 3, 4)))
	assert.True(t, c.next["b"][0].positions.equal(newSet(1, 2, 3, 6)))

	d := c.next["b"][0]
	assert.True(t, d.next["a"][0].positions.equal(newSet(1, 2, 3, 4)))
	assert.True(t, d.next["b"][0].positions.equal(newSet(1, 2, 3)))
}

func TestTransducerFinal(t *testing.T) {
	tr := buildTransducer(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)

	a := tr.start
	assert.False(t, a.final)

	b := a.next["a"][0]
	assert.False(t, b.final)

	c := b.next["b"][0]
	assert.False(t, c.final)

	d := c.next["b"][0]
	assert.True(t, d.final)
}

func TestSameInputTape(t *testing.T) {
	tr := buildTransducer(`<a,b>+<a,c>+<a,d>`)

	a := tr.start
	assert.True(t, a.positions.equal(newSet(1, 2, 3)))
	assert.Equal(t, 1, len(a.out))
	assert.Equal(t, map[string][]string{"a": []string{"b", "c", "d"}}, a.out)

	assert.Equal(t, 1, len(a.next))
	assert.Equal(t, 3, len(a.next["a"]))

	b := a.next["a"]
	assert.True(t, b[0] == b[1] && b[1] == b[2])
	assert.True(t, b[0].positions.equal(newSet(4)))
}
