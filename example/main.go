package main

import (
	"fmt"

	"github.com/szabado/multiglob"
)

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
