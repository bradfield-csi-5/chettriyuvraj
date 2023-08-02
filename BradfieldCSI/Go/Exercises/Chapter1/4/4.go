package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	hashMap := make(map[string]string) // Map duplicate lines to file names -> Assuming argument always contains file names
	for _, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in opening file: %v", err)
		} else {
			mapLinesToFiles(f, filename, hashMap)
		}
	}

	fmt.Println(hashMap) // print raw hashmap rep directly

}

func mapLinesToFiles(f *os.File, filename string, hashMap map[string]string) {
	input := bufio.NewScanner(f)    // Read input
	curMap := make(map[string]bool) // To check if a line has occurred in current file or not
	for input.Scan() {
		text := input.Text()
		_, textInCurFile := curMap[text]
		if !textInCurFile {
			curMap[text] = true               // Record that line has occurred in current file
			_, textInAnyFile := hashMap[text] // check if text has occurred in any file previously
			sep := " "                        // Separator to separate file names
			if !textInAnyFile {
				sep = ""
			}
			hashMap[text] += sep + filename
		}
	}
}
