# multiglob: wildcard based pattern matching

Inspired by a problem I encountered, this matches a string against a list of patterns and tells you which one it
matched against!

This uses a Radix tree under the hood and aims pretty darn fast.

## Theorized usage:

```
func main() {
    mg := multiglob.New()
    mg.Add("foo", "foo*")
    mg.Add("bar", "bar*")
    mg.Build()

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

Glob is already _way_ faster than using a Regex. Let's set a moderate goal for multiglob: less than 20,000 ns/op,
and no allocs.

## Isn't this basically an http router??

Yep! But I didn't want the overhead of `http` and I wanted to write this for fun.