# regular relations

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/catiepg/regular-relations)
[![Go Report Card](https://goreportcard.com/badge/github.com/catiepg/regular-relations)](https://goreportcard.com/report/github.com/catiepg/regular-relations)

Construct a subsequential transducer from a regular expression.

## Background
Finite-state automata equate to regular languages and finite-state transducers equate to regular relations.

Regular relations are closed under _concatenation_, _union_ and _Kleene star (closure)_, here denoted as `.`, `+` and `*` respectively.

A finite-state transducer corresponds to a function from strings to strings. Therefore, to construct a subsequential transducer from a regular expression the base elements of that expression must be pairs of strings (input/output).

## Implementation
The subsequential transducer is constructed in two steps:
* Build a finite-state automata (effectively a finite-state non-deterministic transducer) from the regular expression using the Berry-Sethi construction by treating string pairs as distinct symbols.
* Construct an equivalent 2-subsequential transducer from the non-deterministic transducer as described in [_Finitely Subsequential Transducers_](http://www.cs.nyu.edu/~mohri/pub/finite.ps).

## Example
```go
  regexp := strings.NewReader(`<foo,bar>+<none,>`)
  transducer, _ := relations.Build(regexp)
  
  transducer.Transduce("foo")     // [bar], true
  transducer.Transduce("missing") // [], false
```

### Notes

The regular expression must represent a subsequential function. Otherwise the algorithm will never finish the construction.
