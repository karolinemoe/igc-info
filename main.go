// Assignment 1 Cloud technologies

package main

import (
    "fmt"
    "encoding/json"
	//"time"
	"net/http"
    "google.golang.org/appengine"
	"github.com/gorilla/mux"
)

func main() {

	root := "/igcinfo"
	router := mux.NewRouter()

	//http.HandleFunc(root, homeHandler)
	router.HandleFunc(root + "/api", apiHandler)

	appengine.Main()

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

	fmt.Fprint(w, "API RESPONS")
}

