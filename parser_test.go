package relations

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (p rule) contain(in rune, out string) bool {
	return p.in == in && p.out == out
}

func testRegExpMetadata(regexp string, test func(*parserMeta)) {
	source := strings.NewReader(regexp)
	meta, _ := ComputeParserMeta(source)
	test(meta)
}

func TestFirstPos(t *testing.T) {
	testRegExpMetadata(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`, func(meta *parserMeta) {
		assert.True(t, meta.rootFirst.equal(newSet(1, 2, 3)))
	})
}

func TestFollowPos(t *testing.T) {
	testRegExpMetadata(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`, func(meta *parserMeta) {
		assert.True(t, meta.follow[1].equal(newSet(1, 2, 3)))
		assert.True(t, meta.follow[2].equal(newSet(1, 2, 3)))
		assert.True(t, meta.follow[3].equal(newSet(4)))
		assert.True(t, meta.follow[4].equal(newSet(5)))
		assert.True(t, meta.follow[5].equal(newSet(6)))

		_, isSet := meta.follow[6]
		assert.False(t, isSet)
	})
}

func TestSymbols(t *testing.T) {
	testRegExpMetadata(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`, func(meta *parserMeta) {
		assert.Equal(t, 6, len(meta.rules))
		assert.True(t, meta.rules[1].contain('a', ""))
		assert.True(t, meta.rules[2].contain('b', ""))
		assert.True(t, meta.rules[6].contain('!', ""))
	})
}

func TestFinal(t *testing.T) {
	testRegExpMetadata(`(<a,>+<b,>)*.<a,>.<b,>.<b,>`, func(meta *parserMeta) {
		assert.Equal(t, 6, meta.finalIndex)
	})
}

func TestMulticharFollow(t *testing.T) {
	testRegExpMetadata(`<abc,xy>+<bca,zz>`, func(meta *parserMeta) {
		assert.True(t, meta.follow[1].equal(newSet(2)))
		assert.True(t, meta.follow[2].equal(newSet(3)))
		assert.True(t, meta.follow[3].equal(newSet(7)))
		assert.True(t, meta.follow[4].equal(newSet(5)))
		assert.True(t, meta.follow[5].equal(newSet(6)))
		assert.True(t, meta.follow[6].equal(newSet(7)))

		_, isSet := meta.follow[7]
		assert.False(t, isSet)
	})
}

func TestMulticharRootFirst(t *testing.T) {
	testRegExpMetadata(`<abc,xy>+<bca,zz>`, func(meta *parserMeta) {
		assert.True(t, meta.rootFirst.equal(newSet(1, 4)))
	})
}
