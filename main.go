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
	http.HandleFunc(root+"/api/igcget", getigcHandler)
	http.HandleFunc(root+"/api/igcpost", postigcHandler)

	appengine.Main()
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "API RESPONSE")
}

func getigcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get IGC RESPONSE")
}

func postigcHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post IGC RESPONSE")
}
