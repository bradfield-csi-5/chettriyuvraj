/* Use array pointers instead of slice in reverse function */

package main

import "fmt"

func main() {
	var a [5]int = [5]int{1, 2, 3, 4, 5}
	reverse(&a)
	fmt.Println(a)
}

func reverse(s *[5]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
