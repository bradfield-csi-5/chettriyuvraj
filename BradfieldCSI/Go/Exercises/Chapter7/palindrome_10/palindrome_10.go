package main

import (
	"fmt"
	"sort"
)

type StringSlice []string

func (s StringSlice) Len() int { return len(s)}
func (s StringSlice) Less(i,j int) bool { return s[i] < s[j]}
func (s StringSlice) Swap (i,j int) { s[i], s[j] = s[j], s[i]}

func main() {
	fmt.Println(IsPalindrome(StringSlice([]string{"a","b","c"})))
	fmt.Println(IsPalindrome(StringSlice([]string{"a","b","a"})))
}

func IsPalindrome(s sort.Interface) bool {
	isPalindrome := true
	for i, j := 0, s.Len()-1; i < j ; i,j = i+1, j-1 {
		if !(!s.Less(i,j) && !s.Less(j,i)) { // elements at i,j indices are not equal
			isPalindrome = false
			break
		}
	} 
	return isPalindrome
}