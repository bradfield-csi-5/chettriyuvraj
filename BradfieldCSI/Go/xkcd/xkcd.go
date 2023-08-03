package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"xkcd/fetch"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No search term provided")
	}
	searchTerm := os.Args[1]

	lastComic, err := fetch.FetchLastComic()
	if err != nil {
		log.Fatal(err)
	}

	go fetch.FetchAllComics(lastComic.Num)

	err = fetch.WriteAllComics()
	if err != nil {
		log.Fatalf("error in writing comics to file %w", err)
	}

	err = searchIndex("xkcd.json", searchTerm)
	if err != nil {
		log.Fatal(err)
	}

}

func searchIndex(filename string, searchTerm string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var comic fetch.Comic
		err = decoder.Decode(&comic)
		if err != nil {
			return fmt.Errorf("error reading json file %w", err)
		}

		if strings.Contains(comic.Title, searchTerm) {
			fmt.Println("\n\n\n-------------x--------------x--------------x--------------x--------------x--------------\n\n\n")
			fmt.Printf("Found search term %s!\n\n Comic Title: %s\n\n Comic URL: %s,\n\n Comic Transcript: %s\n", searchTerm, comic.Title, comic.Link, comic.Transcript)
			fmt.Println("\n\n\n-------------x--------------x--------------x--------------x--------------x--------------\n\n\n")
		}
	}

	fmt.Println("Unable to find any (more) matches")
	return nil
}
