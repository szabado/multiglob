# multiglob: matching multiple wildcard based patterns 

Inspired by a problem I encountered at work, this matches a string against a list of patterns and tells you which one it
matched against! Or it will do, once I finish it.

This uses a Radix tree under the hood and aims pretty darn fast.

## Use Case



## Theorized usage:

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
dinky old 2015 laptop, running `comparison_test.go` gives me the following output:

```
$ go test . -bench=.                                                                                                                                                                                                          [0]
goos: linux
goarch: amd64
pkg: github.com/szabado/multiglob
BenchmarkRegex-4   	    2000	    622554 ns/op	   13855 B/op	       3 allocs/op
BenchmarkGlob-4    	   50000	     38622 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/szabado/multiglob	3.686s
```

Glob is referring to [github.com/gobwas/glob](https://github.com/gobwas/glob), which I took some inspiration from.

Glob is already _way_ faster than using a Regex. Let's set a modest goal for multiglob: less than 20,000 ns/op,
and no allocs.

### (Rough) Results

```
go test . -bench=.                                                                                                                                                                                                          [0]
goos: linux
goarch: amd64
pkg: github.com/szabado/multiglob
BenchmarkRegex-4              	    3000	    535890 ns/op	      12 B/op	       0 allocs/op
BenchmarkGlob-4               	   30000	     41002 ns/op	       0 B/op	       0 allocs/op
BenchmarkMultiGlob-4          	 1000000	      1735 ns/op	    2016 B/op	       6 allocs/op
BenchmarkBuilderParseGlob-4   	    1000	   1572638 ns/op	 1664586 B/op	   20889 allocs/op
BenchmarkParseGlob-4          	     500	   3601170 ns/op	 1555202 B/op	   39288 allocs/op
PASS
ok  	github.com/szabado/multiglob	9.036s
```

HAHA! That's quite the savings. Pretty quality work for 3am. Now to figure out if that actually works on
more complex test cases. Crazy what having way less functionality will do for you, innit? The allocs should
also be dealt with, memory thrashing ain't good.

## Isn't this basically an http router??

Yep! But I didn't want the overhead of `http` and I wanted to write this for fun.

# Limitations

This only supports wildcards `*`, and doesn't support any of `glob`'s more advanced functionality, let alone any
Regular Expressions.
