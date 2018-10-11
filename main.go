// Assignment 1 Cloud technologies

package main

import (
    "fmt"
	//"time"
	"net/http"
    "google.golang.org/appengine"
	//"github.com/gorilla/mux"
)

func main() {

	root := "/igcinfo"

	http.HandleFunc(root+"/api", apiHandler)
	http.HandleFunc(root+"/api/igc", igcHandler)

	appengine.Main()
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API RESPONSE")

}

func igcHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "GET")
	case http.MethodPost:
		fmt.Fprint(w, "POST")
	default:
		fmt.Fprint(w, "Error message")
	}

	fmt.Fprint(w, "Get IGC RESPONSE")
}

