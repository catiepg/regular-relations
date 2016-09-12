package relations

import (
	"bytes"
	"math/rand"
	"testing"
)

// writeRandomTo writes a random rune to the given buffer.
func writeRandomTo(b *bytes.Buffer) {
	var alphabet = []rune("abcdefghipqrstuvwxyz")
	n := rand.Intn(10) + 5
	for i := 0; i < n; i++ {
		b.WriteRune(alphabet[rand.Intn(len(alphabet))])
	}
}

// regularRelationExpr generates a random regular relation expression
// according to the following template:
// 		<string, string>+<string, string>+<string, string>+...
// The parameter specifies the number of pairs to be generated.
func regularRelationExpr(pairs int) *bytes.Buffer {
	var b bytes.Buffer

	for i := 1; i <= pairs; i++ {
		b.WriteRune('<')
		writeRandomTo(&b)
		b.WriteRune(',')
		writeRandomTo(&b)
		b.WriteRune('>')

		if i != pairs {
			b.WriteRune('+')
		}
	}
	return &b
}

func benchmarkRelations(b *testing.B, pairs int) {
	for n := 0; n < b.N; n++ {
		Build(regularRelationExpr(pairs))
	}
}

func BenchmarkRelations10(b *testing.B)    { benchmarkRelations(b, 10) }
func BenchmarkRelations100(b *testing.B)   { benchmarkRelations(b, 100) }
func BenchmarkRelations1000(b *testing.B)  { benchmarkRelations(b, 1000) }
func BenchmarkRelations10000(b *testing.B) { benchmarkRelations(b, 10000) }
