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

    mg := mgb.Build()

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
BenchmarkRegex-12               	    3000	    385154 ns/op	      12 B/op	       0 allocs/op
BenchmarkGlob-12                	   50000	     29401 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiGlob-12           	50000000	        29.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkBuilderParseGlob-12    	    1000	   1120589 ns/op	 1664596 B/op	   20889 allocs/op
BenchmarkParseGlob-12           	     500	   2490767 ns/op	 1555203 B/op	   39288 allocs/op
PASS
ok  	github.com/szabado/multiglob	7.324s
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
