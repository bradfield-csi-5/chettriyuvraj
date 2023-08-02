/* Write an in-place function that squashes each run of adjacent Unicode spaces (see unicode.IsSpace)
in a UTF-8-encoded []byte slice into a single ASCII space */

package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {

	s := []byte("\n\n世界\t\tabc\n\n\n\tabbc世界 c\n\n")

	for i := 0; i < len(s); i++ {
		j := i
		for j < len(s) { //finds as many continuous unicode spaces as it can
			r, size := utf8.DecodeRune(s[j:]) // decoding rune at jth index
			if unicode.IsSpace(r) {
				j += size
			} else {
				break
			}
		}

		curR, _ := utf8.DecodeRune(s[i:])
		if unicode.IsSpace(curR) { //if cur char at i is a unicode space
			s[i] = 32                  //set i to ascii space
			s = removeN(s, i+1, j-i-1) // remove all bytes from i + 1 till j
		}

	}

	fmt.Println(string(s))

}

func removeN(s []byte, i int, n int) []byte { // remove n bytes from a byte slice starting from i
	if n > 0 && i+n < len(s) {
		copy(s[i:], s[i+n:])
		return s[:len(s)-n]
	} else {
		return s
	}

}
