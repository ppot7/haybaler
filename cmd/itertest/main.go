package main

import (
	"bufio"
	"fmt"
	"iter"
	"strings"
)

var values = `zero
one
two
three
four
five
six
seven
eight
nine
`

func Stream() iter.Seq[string] {
	return func(yield func(string) bool) {
		strScanner := bufio.NewScanner(strings.NewReader(values))

		for strScanner.Scan() {
			if !yield(strScanner.Text()) {
				return
			}
		}

	}
}

func main() {
	fmt.Println("Starting iteration test")

	recordSet := Stream()

	for record := range recordSet {
		fmt.Println(record)
	}

	fmt.Println("Stopping iteration test")
}
