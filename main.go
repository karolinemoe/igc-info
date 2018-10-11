// Assignment 1 Cloud technologies

package main

import (
    "fmt"
	"net/http"
    "google.golang.org/appengine"

func main() {

	http.HandleFunc("/", homeHandler)

	appengine.Main()

}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, r, http.StatusNotFound)

}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404")
	}
}
