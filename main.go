// Assignment 1 Cloud technologies

package main

import (
    "fmt"
	"net/http"

)
func main() {
        http.HandleFunc("/", handle)

}
func handle(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello, Application!")

}
