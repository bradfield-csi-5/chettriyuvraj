package main

import "fmt"

func main() {

	var a []string = []string{"abc", "abc", "bbc", "b", "c", "d", "e", "e", "fc", "fc", "f"}
	a = removeAdjDupl(a)
	fmt.Println(a)
}

func removeAdjDupl(s []string) []string {
	for i := 0; i < len(s)-1; i++ {
		if s[i+1] == s[i] { // dupl found
			s = remove(s, i+1)
			fmt.Println(s) // prints at each removal
			i--            // keep index rooted here at next iteration to remove other possible duplicates
		}

	}
	return s
}

func remove(s []string, i int) []string {
	if i < len(s) {
		copy(s[i:], s[i+1:])
		return s[:len(s)-1]
	} else {
		return s
	}
}
