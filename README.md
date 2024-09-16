# multiglob: matching multiple wildcard based patterns 

[![GoDoc](https://godoc.org/github.com/szabado/multiglob?status.svg)](https://godoc.org/github.com/szabado/multiglob)
[![codecov](https://codecov.io/gh/szabado/multiglob/branch/master/graph/badge.svg)](https://codecov.io/gh/szabado/multiglob)
[![Go Report Card](https://goreportcard.com/badge/github.com/szabado/multiglob)](https://goreportcard.com/report/github.com/szabado/multiglob)

`multiglob` is a library that for figuring out which glob patterns match a given string.

It uses a Radix tree under the hood and is optimized for speed.

## Example

```
func main() {
	mgb := multiglob.New()
	mgb.MustAddPattern("foo", "foo*")
	mgb.MustAddPattern("bar", "bar*")
	mgb.MustAddPattern("eyyyy!", "*ey*")

	mg := mgb.MustCompile()

	if mg.Match("football") {
		fmt.Println("I matched!")
	}

	matches := mg.FindAllPatterns("barney stinson")
	if matches == nil {
		fmt.Println("Oh no, I didn't match any pattern")
		return
	}

	for _, match := range matches {
		fmt.Println("I matched: ", match)
	}
}
```

## Performance Comparison

There's a few alternatives to using multiglob:
- Looping over the patterns using [github.com/gobwas/glob](https://github.com/gobwas/glob), testing which ones match.
- Looping over the patterns using `regexp`, testing which ones match.

The benchmarks in `comparision_test.go` aim to compare these different possible implementations, with the following results:

```
$ go test . -bench=.
goos: linux
goarch: amd64
pkg: github.com/szabado/multiglob
BenchmarkMultiMatchRegex-4            	    1000	   1260832 ns/op	      37 B/op	       0 allocs/op
BenchmarkMultiMatchGlob-4             	   20000	     61867 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiMatchMultiGlob-4        	 1000000	      4405 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiNotMatchRegex-4         	    3000	    544801 ns/op	      12 B/op	       0 allocs/op
BenchmarkMultiNotMatchGlob-4          	   30000	     41110 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiNotMatchMultiGlob-4     	 3000000	       586 ns/op	       0 B/op	       0 allocs/op
BenchmarkSingleMatchRegex-4           	   50000	     32161 ns/op	       0 B/op	       0 allocs/op
BenchmarkSingleMatchGlob-4            	 5000000	       323 ns/op	       0 B/op	       0 allocs/op
BenchmarkSingleMatchMultiGlob-4       	50000000	        29.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseRegex-4                 	     500	   3645894 ns/op	 3188164 B/op	   24300 allocs/op
BenchmarkParseGlob-4                  	     500	   3609843 ns/op	 1555202 B/op	   39288 allocs/op
BenchmarkParseMultiGlob-4             	    1000	   2039323 ns/op	 1870861 B/op	   22935 allocs/op
BenchmarkMultiGlobFindAllPatterns-4   	  300000	      5540 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiGlobFindPattern-4       	  300000	      4826 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiGlobFindAllGlobs-4      	  200000	      6327 ns/op	     528 B/op	       7 allocs/op
BenchmarkMultiGlobFindGlobs-4         	  300000	      4838 ns/op	     128 B/op	       5 allocs/op
PASS
ok  	github.com/szabado/multiglob	31.539s
```

All the `*Multi*` benchmarks are using 720 patterns, and all the `*Single*` tests only include one pattern.

Multiglob is ~10x faster than `glob` when it has a single pattern loaded, and about 14x faster than `glob` when there are all 720 patterns loaded. In both cases, Multiglob and glob are orders of magnitude faster than the `regexp` based solution.

*Note: these benchmarks were run in early 2019 and might not reflect the current performance of `glob` or `regexp`.*

## Limitations

This only supports wildcards (`*`) and character ranges (`[ab]`, `[^cd]`, `[e-h]`, etc.).
