package relations

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestTransducer(t *testing.T) {
	tr := buildTransducer(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)

	state1 := tr.start
	assert.Equal(t, 1, state1.index)
	assert.Equal(t, 2, state1.next["a"][0].state.index)
	assert.Equal(t, 1, state1.next["b"][0].state.index)
	assert.Equal(t, 1, len(state1.next["a"]))
	assert.Equal(t, 1, len(state1.next["b"]))

	state2 := state1.next["a"][0].state
	assert.Equal(t, 2, state2.next["a"][0].state.index)
	assert.Equal(t, 3, state2.next["b"][0].state.index)
	assert.Equal(t, 1, len(state2.next["a"]))
	assert.Equal(t, 1, len(state2.next["b"]))

	state3 := state2.next["b"][0].state
	assert.Equal(t, 2, state3.next["a"][0].state.index)
	assert.Equal(t, 4, state3.next["b"][0].state.index)
	assert.Equal(t, 1, len(state3.next["a"]))
	assert.Equal(t, 1, len(state3.next["b"]))

	state4 := state3.next["b"][0].state
	assert.Equal(t, 2, state4.next["a"][0].state.index)
	assert.Equal(t, 1, len(state4.next["a"]))
	assert.Equal(t, 1, len(state4.next["b"]))
	assert.Equal(t, 1, state4.next["b"][0].state.index)
}

func TestTransducerFinal(t *testing.T) {
	tr := buildTransducer(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)

	a := tr.start
	assert.False(t, a.final)

	b := a.next["a"][0].state
	assert.False(t, b.final)

	c := b.next["b"][0].state
	assert.False(t, c.final)

	d := c.next["b"][0].state
	assert.True(t, d.final)
}

func TestSameInputTape(t *testing.T) {
	tr := buildTransducer(`<a,b>+<a,c>+<a,d>`)

	state1 := tr.start
	assert.Equal(t, 1, state1.index)

	assert.Equal(t, 1, len(state1.next))
	assert.Equal(t, 3, len(state1.next["a"]))

	next := state1.next["a"]
	assert.True(t, next[0].state == next[1].state && next[1].state == next[2].state)
	assert.Equal(t, 2, next[0].state.index)
	assert.True(t, next[0].state.final)
}

func TestMulticharInputTransducer(t *testing.T) {
	tr := buildTransducer(`<abc,xy>+<aca,zz>`)

	state1 := tr.start
	assert.Equal(t, 1, state1.index)
	assert.Equal(t, 2, len(state1.next["a"]))

	assert.Equal(t, 2, state1.next["a"][0].state.index)
	assert.Equal(t, "xy", state1.next["a"][0].out)

	assert.Equal(t, 3, state1.next["a"][1].state.index)
	assert.Equal(t, "zz", state1.next["a"][1].out)

	state2 := state1.next["a"][0].state
	assert.Equal(t, 0, len(state2.next["a"]))
	assert.Equal(t, 0, len(state2.next["c"]))
	assert.Equal(t, 4, state2.next["b"][0].state.index)
	assert.Equal(t, "", state2.next["b"][0].out)

	state3 := state1.next["a"][1].state
	assert.Equal(t, 0, len(state3.next["a"]))
	assert.Equal(t, 0, len(state3.next["b"]))
	assert.Equal(t, 5, state3.next["c"][0].state.index)
	assert.Equal(t, "", state3.next["c"][0].out)

	state4 := state2.next["b"][0].state
	assert.Equal(t, 0, len(state4.next["a"]))
	assert.Equal(t, 0, len(state4.next["b"]))
	assert.Equal(t, 6, state4.next["c"][0].state.index)
	assert.Equal(t, "", state4.next["c"][0].out)

	state5 := state3.next["c"][0].state
	assert.Equal(t, 0, len(state5.next["b"]))
	assert.Equal(t, 0, len(state5.next["c"]))
	assert.Equal(t, 6, state5.next["a"][0].state.index)
	assert.Equal(t, "", state5.next["a"][0].out)

	state6 := state5.next["a"][0].state
	assert.True(t, state6.final)
}
