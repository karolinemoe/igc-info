// Assignment 1 Cloud technologies

package main

import (
	"encoding/json"
	"fmt"
	"time"
	"net/http"
    "google.golang.org/appengine"
)

// store boot time
var startTime = time.Now()

func main() {

	root := "/igcinfo"

	http.HandleFunc(root+"/api", apiHandler)
	http.HandleFunc(root+"/api/igc", igcHandler)

	appengine.Main()
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "API RESPONSE")

	type Info struct {
		Uptime string `json:"time"`
		Info string `json:"info"`
		Version string `json:"version"`
	}

	info := &Info{
		Uptime: uptime(),
		Info: "Service for IGC tracks.",
		Version: "v1",
	}

	i, _ := json.Marshal(info)

	fmt.Println(string(i))
	fmt.Fprint(w, string(i))

}

func uptime() string{
	// calculate time since boot in UNIX time (seconds) and convert that into a time.Time object
	uptime := time.Now().Unix() - startTime.Unix()
	ut := time.Unix(uptime, 0)

	// calculate each component of the ISO8601 format
	years, months, days 	:= ut.Date()
	hours, minutes, seconds := ut.Clock()

	// subtract 1 as all three is not 0-indexed
	months, hours, days = months-1, hours-1, days-1

	// construct final string with all components
	return fmt.Sprintf("%s%d%s%d%s%d%s%d%s%d%s%d%s", "P", absInt(int64(years - 1970)), "Y", months, "M", days, "DT", hours, "H", minutes, "M", seconds, "S")
}

// take absolute value of n
func absInt(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
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
}

