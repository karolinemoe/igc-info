// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package main
import (
    "fmt"
	"net/http"
    "google.golang.org/appengine"
	"github.com/marni/goigc"

)
func main() {
        http.HandleFunc("/", handle)
        appengine.Main()


	s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(s)
	if err != nil {
		fmt.Errorf("Problem reading the track", err)
	}

	fmt.Printf("Pilot: %s, gliderType: %s, date: %s",
		track.Pilot, track.GliderType, track.Date.String())
}
func handle(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, Application!")

}
