// Assignment 1 Cloud technologies

package main

import (
    "fmt"
	//"time"
	"net/http"
    "google.golang.org/appengine"
	"github.com/gorilla/mux"
)

func main() {

	root := "/igcinfo"
	router := mux.NewRouter()

	router.HandleFunc(root+"/api", apiHandler)

	appengine.Main()

}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "API RESPONS")
}

