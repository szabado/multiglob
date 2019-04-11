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

The lazy way you can do what multiglob does is loop over the patterns. That's _sloooow_,
so I wanted to make a better way. On my laptop, running `comparison_test.go` outputs this:

```
$ go test . -bench=.                                                                                                                                                                                                        [130]
goos: linux
goarch: amd64
pkg: github.com/szabado/multiglob
BenchmarkMatchRegex-4       	    2000	    557203 ns/op	      18 B/op	       0 allocs/op
BenchmarkMatchGlob-4        	   30000	     41237 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatchMultiGlob-4   	 3000000	       581 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseGlob-4        	     500	   3595179 ns/op	 1555203 B/op	   39288 allocs/op
BenchmarkParseMultiGlob-4   	    1000	   1912159 ns/op	 1854498 B/op	   22924 allocs/op
PASS
ok  	github.com/szabado/multiglob	9.512s
```

Glob is referring to [github.com/gobwas/glob](https://github.com/gobwas/glob), which I
took some inspiration from.

Glob is already _way_ faster than using a Regex. MultiGlob is _way_ faster than doing
using Glob naively.

## Isn't this basically an http router??

Yep! But I didn't want the overhead of `http` and I wanted to write this for fun.

# Limitations

This only supports wildcards `*`.
