// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"strings"
// )

// const apiKey = "bcebdc6d"
// const baseURL = "https://omdbapi.com/"
// const pageLimit = 3 // max number of data pages to be fetched

// type Movie struct {
// 	Title  string
// 	Type   string
// 	Poster string // Poster URL - can use custom marshaling to cast fields into more clear names
// 	Year   int
// }

// type OMDBResponse struct {
// 	Movies []Movie
// 	Error  string
// }

// func FetchMovieData(params map[string]string) []Movie {
// 	queryURL := baseURL + "?" + apiKey
// 	queryURL = ConstructURL(queryURL, params) // add remaining params
// 	pageNo := 1
// 	apiResponse := OMDBResponse{make([]Movie, 0), "Movie found"} // dummy response to start loop
// 	movieData := []Movie{}

// 	for apiResponse.Error != "" && pageNo < pageLimit {
// 		queryURL = queryURL[:strings.LastIndex(queryURL, "&")+1] + strconv.Itoa(pageNo) // updating page number at each iteration

// 		response, queryErr := http.Get(queryURL) //fetch data from URL
// 		fmt.Println(queryURL)
// 		if queryErr != nil {
// 			log.Fatal(queryErr)
// 		}
// 		data, readErr := io.ReadAll(response.Body) // read data on this page
// 		response.Body.Close()
// 		if readErr != nil {
// 			fmt.Println("Error in reading page number" + strconv.Itoa(pageNo))
// 			continue
// 		}

// 		if err := json.Unmarshal(data, &apiResponse); err != nil { //unmarshal the data into response
// 			log.Fatal(err)
// 		}

// 		if len(apiResponse.Movies) > 0 {
// 			movieData = append(movieData, apiResponse.Movies...)
// 		}
// 	}

// 	return movieData

// }

// func ConstructURL(URL string, params map[string]string) string {
// 	for key, val := range params {
// 		URL += "&" + key + "=" + url.QueryEscape(val)
// 	}
// 	fmt.Println(URL)
// 	return URL
// }
