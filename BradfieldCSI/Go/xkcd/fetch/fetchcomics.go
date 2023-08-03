// package which uses json interface of xkcd comics to fetch each comic and create an offline index

package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var maxFetchGoroutines int = 100
var baseUrl string = "https://xkcd.com"
var urlSuffix string = "/info.0.json"

type Comic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"string"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Image      string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

var comicNums = make(chan int, maxFetchGoroutines)   // channel where comic numbers are sent and received by goroutines that fetch them
var comics = make(chan *Comic, maxFetchGoroutines*2) // channel to send unmarshaled goroutines to, to write
var errors = make(chan error, maxFetchGoroutines*2)  // channel to handle errors
var tokens = make(chan struct{}, maxFetchGoroutines) // channel to ensure no more than 'maxGoroutines' goroutines are active at a time

func fetchComic(fetchUrl string) (*Comic, error) {
	var comic *Comic = new(Comic)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Get(fetchUrl)
	if err != nil {
		return nil, err

	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err

	}

	err = json.Unmarshal(data, comic)
	if err != nil {
		return nil, err
	}

	return comic, nil
}

func fetchIthComic(comicNumber int) {
	fetchUrl := baseUrl + "/" + strconv.Itoa(comicNumber) + urlSuffix
	fmt.Printf("Fetching comic %d!\n", comicNumber)
	comic, err := fetchComic(fetchUrl)
	if err != nil {
		errors <- fmt.Errorf("comic num. %d: %w", comicNumber, err)
	}
	comics <- comic
}

func FetchLastComic() (*Comic, error) {
	lastComicUrl := baseUrl + urlSuffix
	comic, err := fetchComic(lastComicUrl)
	if err != nil {
		return nil, fmt.Errorf("error fetching last comic %w", err)
	}
	return comic, nil
}

func FetchAllComics(lastComicIndex int) {

	go func() {
		err := <-errors
		log.Println(err)
	}()

	var wg sync.WaitGroup

	for i := 1; i <= lastComicIndex; i++ { // send comic numbers on comicNum to start fetch
		tokens <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer func() { <-tokens }()
			defer wg.Done()
			fetchIthComic(i)
		}(i)
	}

	wg.Wait()

	close(comics)

}

func writeComic(comic *Comic, encoder *json.Encoder) error {
	if err := encoder.Encode(comic); err != nil {
		return fmt.Errorf("error while writing to file %v", err)
	}
	return nil
}

func WriteAllComics() error {
	f, err := os.OpenFile("xkcd.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("unable to open file to write comic %v", err)
	}

	encoder := json.NewEncoder(f)
	i := 0
	for comic := range comics {
		err = writeComic(comic, encoder)
		if err != nil {
			log.Println(err)
		}
		i++
	}

	return nil
}
