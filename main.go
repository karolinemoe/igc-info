// Assignment 1 Cloud technologies

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/hashstructure"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"google.golang.org/appengine"
)

// IGCTrack struct for the track data
type IGCTrack struct {
	HDate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	ID          string    `json:"track_id"`
	TrackLength float64   `json:"track_length"`
	Data        igc.Track `json:"-"`
}

// store boot time
var startTime = time.Now()
var igcTracks []IGCTrack

//var currentID = 0


func main() {

	root := "/igcinfo"
	r := mux.NewRouter()
	route := r.PathPrefix("/igcinfo/api").Subrouter()

	http.HandleFunc(root+"/api", apiHandler)
	http.HandleFunc(root+"/api/igc", igcHandler)
	route.HandleFunc("/igc/{id}", getIgc).Methods("GET")
	route.HandleFunc("/igc/{id}/", getIgc).Methods("GET")
	route.HandleFunc("/igc/{id}/{field}", getIgc).Methods("GET")

	http.Handle("/", r)


	appengine.Main()
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	type Info struct {
		Uptime string `json:"uptime"`
		Info string `json:"info"`
		Version string `json:"version"`
	}

	info := &Info{
		Uptime: uptime(),
		Info: "Service for IGC tracks.",
		Version: "v1",
	}

	i, _ := json.MarshalIndent(info, "", " ")

	fmt.Fprint(w, string(i))
}

// calculates time since boot un UNIX time and converts it into a Time object, afterwards converst it to a final string with all components
func uptime() string{

	uptime := time.Now().Unix() - startTime.Unix()
	ut := time.Unix(uptime, 0)

	years, months, days 	:= ut.Date()
	hours, minutes, seconds := ut.Clock()
	months, days = months-1, days-1

	return fmt.Sprintf("%s%d%s%d%s%d%s%d%s%d%s%d%s", "P", absInt(int64(years - 1970)), "Y", months, "M", days, "DT", hours, "H", minutes, "M", seconds, "S")
}

// absolute value of n
func absInt(n int64) int64 {

	if n < 0 {
		return -n
	}
	return n
}

func igcHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	// Case GET:
	case http.MethodGet:
		ids := []string{}
		for _, track := range igcTracks { ids = append(ids, track.ID) }
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ids)

	// Case POST:
	case http.MethodPost:

		var body struct { URL string }

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			http.Error(w, err.Error(), 400); return
		}

		// if object has no 'URL' property attached, reject furter processing
		if body.URL == "" {
			http.Error(w, "Request does not have 'URL' property", 400); return
		}


		/*body = make(map[string]interface{})
		_ = json.NewDecoder(r.Body).Decode(&body)*/

		newID := newTrack(body.URL, w)
		if newID == "" {
			http.Error(w, "Not able to process the URL", http.StatusBadRequest)
			return
		}

		default:
		fmt.Fprint(w, "Error message")
	}
}

func newTrack(url string, w http.ResponseWriter) string {
	fmt.Fprint(w, url + "       ")
	igcData, err := igc.ParseLocation(url)

	if err != nil {
		fmt.Println("Problem reading the track")
		return ""
	}

	checksum, err := hashstructure.Hash(igcData, nil)
	if err != nil {
		http.Error(w, "Problem generating checksum", 400)
	}
	trackID := strconv.Itoa(int(checksum))

	//trackID := strconv.Itoa(currentID+1)
	//currentID = currentID+1

	// add track to memory if it doesn't exist
	if !trackExist(trackID) {

		trackMetaData := IGCTrack{
			HDate:       igcData.Date,
			Pilot:       igcData.Pilot,
			Glider:      igcData.GliderType,
			GliderID:    igcData.GliderID,
			TrackLength: calcTrackLength(igcData.Points),
			ID:          trackID,
			Data:        igcData,
		}
		igcTracks = append(igcTracks, trackMetaData)
	}

	type IGCid struct {
		ID string `json:"id"`
	}

	json.NewEncoder(w).Encode(IGCid{ID: trackID})
	return trackID
}

func getIgc(w http.ResponseWriter, r *http.Request) {

	// extract query parameters
	params := mux.Vars(r)

	// look for track
	track, err := findTrackWithID(params["id"])
	if err != nil {
		http.Error(w, "No Content", 204); return
	}

	// if no field entered return the whole element
	_, ok :=  params["field"]
	if !ok {
		w.Header().Set("Content-Type", "application/jsons")
		json.NewEncoder(w).Encode(track); return
	}

	// use a map to match query with track metadata field
	trackMetaDataMap := map[string]interface{}{
		"pilot":        track.Pilot,
		"glider":       track.Glider,
		"glider_id":    track.GliderID,
		"track_length": track.TrackLength,
		"H_date":       track.HDate.String(),
	}

	if trackMetaDataMap[params["field"]] == nil {
		http.Error(w, "No proper field specified", 400); return
	}

	w.Header().Set("Content-Type", "application/text")
	// fmt.Fprint(w, trackMetaDataMap[params["field"]])
	json.NewEncoder(w).Encode(trackMetaDataMap[params["field"]])
}

// search through the tracks for specific ID and return that, else return an error
func findTrackWithID(ID string) (IGCTrack, error) {

	var t IGCTrack
	for i := 0; i < len(igcTracks); i++ {
		if igcTracks[i].ID == ID {
			return igcTracks[i], nil
		}
	}
	return t, errors.New("track not found")
}

// check if track already exists in the memory
func trackExist(trackID string) bool {

	_, err := findTrackWithID(trackID)
	return err == nil
}


// calculate the total lenth of a track based on its waypoints
func calcTrackLength(points []igc.Point) float64 {

	var tl float64
	for i := 0; i < len(points)-1; i++ {
		tl += points[i].Distance(points[i+1])
	}
	return tl
}