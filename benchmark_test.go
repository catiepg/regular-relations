package relations

import (
	"bytes"
	"math/rand"
	"testing"
)

var alphabet = []rune("abcdefghipqrstuvwxyz")

func writeRandomTo(b *bytes.Buffer) {
	n := rand.Intn(10) + 5

	for i := 0; i < n; i++ {
		b.WriteRune(alphabet[rand.Intn(len(alphabet))])
	}
}

func newSubTransducer(pairCount int) *Subsequential {
	var b bytes.Buffer

	for i := 1; i <= pairCount; i++ {
		b.WriteRune('<')
		writeRandomTo(&b)
		b.WriteRune(',')
		writeRandomTo(&b)
		b.WriteRune('>')

		if i != pairCount {
			b.WriteRune('+')
		}
	}

	sub, _ := NewSubsequential(&b)
	return sub
}

func TestLargeInput(t *testing.T) {
	newSubTransducer(200)
}
