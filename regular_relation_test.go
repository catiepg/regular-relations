package relations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRegularRelation(regexp string, test func(*RegularRelation)) {
	source := strings.NewReader(regexp)
	rr, _ := Build(source)
	test(rr)
}

func TestBuildRegularRelation(t *testing.T) {
	testRegularRelation(`<abc,xyz>+<acc,qwe>`, func(rr *RegularRelation) {
		out, ok := rr.Transduce("abc")
		assert.Equal(t, 1, len(out))
		assert.Equal(t, "xyz", out[0])
		assert.True(t, ok)
	})
}

func TestSubsequentialTransducerStructure(t *testing.T) {
	testRegularRelation(`<abc,xyz>+<acc,qwe>`, func(rr *RegularRelation) {
		assert.Equal(t, 1, len(rr.start.next))
		assert.Equal(t, "", rr.start.out['a'])

		state2 := rr.start.next['a']
		assert.Equal(t, 2, len(state2.next))
		assert.Equal(t, "xyz", state2.out['b'])
		assert.Equal(t, "qwe", state2.out['c'])

		state3 := state2.next['b']
		state4 := state2.next['c']
		assert.Equal(t, 1, len(state3.next))
		assert.Equal(t, 1, len(state4.next))
	})
}
