package relations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransducer(t *testing.T) {
	source := strings.NewReader(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)
	tr, _ := NewTransducer(source)

	state1 := tr.root
	assert.Equal(t, 1, state1.index)
	assert.Equal(t, 2, state1.next['a'][0].state.index)
	assert.Equal(t, 1, state1.next['b'][0].state.index)
	assert.Equal(t, 1, len(state1.next['a']))
	assert.Equal(t, 1, len(state1.next['b']))

	state2 := state1.next['a'][0].state
	assert.Equal(t, 2, state2.next['a'][0].state.index)
	assert.Equal(t, 3, state2.next['b'][0].state.index)
	assert.Equal(t, 1, len(state2.next['a']))
	assert.Equal(t, 1, len(state2.next['b']))

	state3 := state2.next['b'][0].state
	assert.Equal(t, 2, state3.next['a'][0].state.index)
	assert.Equal(t, 4, state3.next['b'][0].state.index)
	assert.Equal(t, 1, len(state3.next['a']))
	assert.Equal(t, 1, len(state3.next['b']))

	state4 := state3.next['b'][0].state
	assert.Equal(t, 2, state4.next['a'][0].state.index)
	assert.Equal(t, 1, len(state4.next['a']))
	assert.Equal(t, 1, len(state4.next['b']))
	assert.Equal(t, 1, state4.next['b'][0].state.index)
}

func TestTransducerFinal(t *testing.T) {
	source := strings.NewReader(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`)
	tr, _ := NewTransducer(source)

	a := tr.root
	assert.False(t, a.final)

	b := a.next['a'][0].state
	assert.False(t, b.final)

	c := b.next['b'][0].state
	assert.False(t, c.final)

	d := c.next['b'][0].state
	assert.True(t, d.final)
}

func TestSameInputTape(t *testing.T) {
	source := strings.NewReader(`<a,b>+<a,c>+<a,d>`)
	tr, _ := NewTransducer(source)

	state1 := tr.root
	assert.Equal(t, 1, state1.index)

	assert.Equal(t, 1, len(state1.next))
	assert.Equal(t, 3, len(state1.next['a']))

	next := state1.next['a']
	assert.True(t, next[0].state == next[1].state && next[1].state == next[2].state)
	assert.Equal(t, 2, next[0].state.index)
	assert.True(t, next[0].state.final)
}

func TestMulticharInputTransducer(t *testing.T) {
	source := strings.NewReader(`<abc,xy>+<aca,zz>`)
	tr, _ := NewTransducer(source)

	state1 := tr.root
	assert.Equal(t, 1, state1.index)
	assert.Equal(t, 2, len(state1.next['a']))
}
