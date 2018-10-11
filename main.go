// Assignment 1 Cloud technologies

package main

import (
    "fmt"
    "encoding/json"
	"time"
	"net/http"
    "google.golang.org/appengine"
	"github.com/gorilla/mux"
)

func main() {

	root := "/igcinfo"
	router := mux.NewRouter()

	http.HandleFunc(root, homeHandler)
	router.HandleFunc(root + "/api", apiHandler).Methods("GET")

	appengine.Main()

}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	errorHandler(w, r, http.StatusNotFound)
}

func apiHandler(w http.ResponseWriter, r *http.Request){
	// API info object schema
	type APIInfo struct {
		//Uptime  string `json:"uptime"`
		Info    string `json:"info"`
		Version string `json:"version"`
		//BootTime string `json:"boot_time"`
	}

	// construct the current response
	info := APIInfo{
		//Uptime:  uptimeFormat(),
		//BootTime: startTime.UTC().Format(time.RFC3339),
		Info:    "Service for IGC tracks",
		Version: "v1",
	}

	json.NewEncoder(w).Encode(info)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 page not found")
	}
}
