/* Write a program wordfreq to report the frequency of each word in an input text file.
Call input.Split(bufio.ScanWords)before the first call to Scan to break the input into words instead of lines. */

package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	text := "This is a source of       text  \n\n new lines and    spaces and \n all"
	textFile := strings.NewReader(text) // artificial text file

	hashMap := make(map[string]int) // hashmap for counting words

	scanner := bufio.NewScanner(textFile)
	scanner.Split(bufio.ScanWords) // splitting into words

	for scanner.Scan() { // count
		// fmt.Println(scanner.Text())
		hashMap[scanner.Text()]++
	}

	fmt.Println(hashMap)

}
