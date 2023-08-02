/* Add additional handlers so that clients can create, read, update, and delete database entries.
For example, a request of the form /update?item=socks&price=6 will update the price of an 
item in the inventory and report an error if the item does not exist or if the price is invalid. 
(Warning: this change introduces concurrent variable updates.) */


package main

import (
	"fmt"
	"net/http"
	"log"
	"strconv"
)

type database map[string]int


func (db database) list(w http.ResponseWriter, req *http.Request) {
    for item, price := range db {
        fmt.Fprintf(w, "%s: %d\n", item, price)
	} 
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%d\n", price)
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	_, ok := (*db)[item]
	if !ok { // if old price doesnt exist
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "No such item %q\n", item)
		return
	}

	newPrice, err := strconv.Atoi(req.URL.Query().Get("price"))
	if err != nil { // new price empty
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "Error in price conversion %q\n", item)
		return
	}

	(*db)[item] = newPrice
	fmt.Fprintf(w, "Successfully update price of %q to %d!\n", item, newPrice)
}

func main(){
	db := database{"shoes": 50, "coat": 10, "socks": 5}
	http.HandleFunc("/list", db.list) //registers with global servemux
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/update", db.update)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}