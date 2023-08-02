/* The JSON-based web service of the Open Movie Database lets you search https://omdbapi.com/ for a movie by name and download its poster image.
Write a tool poster that downloads the poster image for the movie named on the command line */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const apiKey = "bcebdc6d"
const baseURL = "https://omdbapi.com/"
const pageLimit = 3 // max number of data pages to be fetched

type Movie struct {
	Title     string
	Type      string
	PosterURL string `json:"Poster"`
}

type OMDBResponse struct {
	Movie []Movie `json:"Search"`
	Error string
}

func main() {
	params := map[string]string{}
	for _, movieName := range os.Args[1:] {
		params["s"] = movieName
		movieData := fetchMovieData(params)
		fetchPosters(movieData)
	}

}

func fetchPosters(movieData []Movie) {
	for i, movie := range movieData {
		response, fetchErr := http.Get(movie.PosterURL) //fetch data from poster URL
		if fetchErr != nil {
			log.Fatal(fetchErr)
		}

		body, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		response.Body.Close()

		fileName := movieData[0].Title + strconv.Itoa(i) + movie.PosterURL[strings.LastIndex(movie.PosterURL, "."):] // title + index + extension is the file name
		os.WriteFile(fileName, body, 0666)
	}

}

func fetchMovieData(params map[string]string) []Movie {
	queryURL := baseURL + "?apiKey=" + apiKey
	queryURL = constructURL(queryURL, params) + "&page=1" // add remaining params
	pageNo := 1
	apiResponse := OMDBResponse{make([]Movie, 0), "Movie found"} // dummy response to start loop
	movieData := []Movie{}

	for apiResponse.Error != "" && pageNo < pageLimit {
		queryURL = queryURL[:strings.LastIndex(queryURL, "&")+1] + "page=" + strconv.Itoa(pageNo) // updating page number at each iteration

		response, queryErr := http.Get(queryURL) //fetch data from URL
		// fmt.Println(queryURL)
		if queryErr != nil {
			log.Fatal(queryErr)
		}
		data, readErr := io.ReadAll(response.Body) // read data on this page
		response.Body.Close()
		if readErr != nil {
			fmt.Println("Error in reading page number" + strconv.Itoa(pageNo))
			continue
		}

		if err := json.Unmarshal(data, &apiResponse); err != nil { //unmarshal the data into response
			log.Fatal(err)
		}

		if len(apiResponse.Movie) > 0 {
			movieData = append(movieData, apiResponse.Movie...)
		}

		// fmt.Println(string(data))

		pageNo++
	}

	return movieData

}

func constructURL(URL string, params map[string]string) string {
	for key, val := range params {
		URL += "&" + key + "=" + url.QueryEscape(val)
	}
	// fmt.Println(URL)
	return URL
}
