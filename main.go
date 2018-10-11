// Assignment 1 Cloud technologies

package main

import (
    "fmt"
	"log"
	"net/http"
    "google.golang.org/appengine"
	"github.com/gorilla/mux"
	"os"
)
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := mux.NewRouter()
	appBase := "/igcinfo"

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(":"+port, router))



	http.HandleFunc("/", handle)
        appengine.Main()

}
func handle(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, Application!")

}
