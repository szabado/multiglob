# multiglob: matching multiple wildcard based patterns 

Inspired by a problem I encountered at work, this matches a string against a list of patterns and tells you which one it
matched against! Or it will do, once I finish it.

This uses a Radix tree under the hood and aims pretty darn fast.

## Use Case


## Usage:

```
func main() {
    mgb := multiglob.New()
    mgb.AddPattern("foo", "foo*")
    mgb.AddPattern("bar", "bar*")

    mg := mgb.MustCompile()

    mg.Match("football", func(name string) {
        fmt.Printf("I matched: %s\n", name)
    })
}
```

## Performance Comparison

The lazy way you can do what multiglob does is loop over the patterns and do whatever 
you want when it matches. That's _sloooow_, so I wanted to make a better way. On my 
laptop, running `comparison_test.go` gives me the following output:

```
$ go test -bench=. .
goos: darwin
goarch: amd64
pkg: github.com/szabado/multiglob
BenchmarkMatchRegex-12        	    3000	    382638 ns/op	      12 B/op	       0 allocs/op
BenchmarkMatchGlob-12         	   50000	     28709 ns/op	       0 B/op	       0 allocs/op
BenchmarkMatchMultiGlob-12    	 3000000	       435 ns/op	       0 B/op	       0 allocs/op
BenchmarkParseGlob-12         	     500	   2509470 ns/op	 1555204 B/op	   39288 allocs/op
BenchmarkParseMultiGlob-12    	    1000	   1280321 ns/op	 1854440 B/op	   22924 allocs/op
PASS
ok  	github.com/szabado/multiglob	7.724s
```

Glob is referring to [github.com/gobwas/glob](https://github.com/gobwas/glob), which I
took some inspiration from.

Glob is already _way_ faster than using a Regex. MultiGlob is _way_ faster than doing
using glob naively.

## Isn't this basically an http router??

Yep! But I didn't want the overhead of `http` and I wanted to write this for fun.

# Limitations

This only supports wildcards `*`, and doesn't support any of `glob`'s more advanced functionality, let alone any
Regular Expressions.
