# multiglob: matching multiple wildcard based patterns 

Inspired by a problem I encountered at work, this matches a string against a list of patterns and tells you which
one it matched against!

This uses a Radix tree under the hood and aims pretty darn fast.

## Usage:

See under `/example` for the usage, but it's copied below.

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

The lazy way you can do what multiglob does is loop over the patterns. That's _sloooow_, but I wanted
to know how slow. I benchmarked it against the standard library `regexp` package as well as 
[github.com/gobwas/glob](https://github.com/gobwas/glob), which I took some inspiration from.
On my laptop, the benchmarks in `comparison_test.go` produce this:

```
$ go test . -bench=.                                                                                                                                                                                                          [0]
goos: linux
goarch: amd64
pkg: github.com/szabado/multiglob
BenchmarkMultiMatchRegex-4            	    1000	   1263544 ns/op	      37 B/op	       0 allocs/op
BenchmarkMultiMatchGlob-4             	   20000	     61713 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiMatchMultiGlob-4        	  300000	      4317 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiNotMatchRegex-4         	    3000	    542037 ns/op	      12 B/op	       0 allocs/op
BenchmarkMultiNotMatchGlob-4          	   30000	     41217 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiNotMatchMultiGlob-4     	 3000000	       603 ns/op	       0 B/op	       0 allocs/op
BenchmarkSingleMatchRegex-4           	   50000	     32057 ns/op	       0 B/op	       0 allocs/op
BenchmarkSingleMatchGlob-4            	 5000000	       323 ns/op	       0 B/op	       0 allocs/op
BenchmarkSingleMatchMultiGlob-4       	50000000	        29.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseRegex-4                 	     500	   3681185 ns/op	 3188165 B/op	   24300 allocs/op
BenchmarkParseGlob-4                  	     500	   3640196 ns/op	 1555202 B/op	   39288 allocs/op
BenchmarkParseMultiGlob-4             	    1000	   1917451 ns/op	 1854570 B/op	   22924 allocs/op
BenchmarkMultiGlobFindAllPatterns-4   	  300000	      5405 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/szabado/multiglob	23.991s
```

All the `Multi` benchmarks are matching one string across a bunch of patterns, and all the `Single` one are
one pattern. Basically? It's **fast**. It's more then 10 times faster than Glob, and that ratio gets better the
more patterns there are.

Based on the benchmark, it also has better performance growth than Glob.
The Multi tests have 720 patterns, and MultiGlob took ~150 times longer to execute the Multi tests compared to
Glob's ~190 times increase.

Glob is already _way_ faster than using a Regex. MultiGlob is _way_ faster than doing
using Glob naively; you just have to accept the reduced functionality.

### Open Questions


## Isn't this basically an http router??

Yep! But I didn't want the overhead of `http` and I wanted to write this for fun.

## Limitations

This only supports wildcards `*`. If you need more, I'd suggest checking out [glob](https://github.com/gobwas/glob).
