// Assignment 1 Cloud technologies

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"errors"
	"time"
	"net/http"

    "google.golang.org/appengine"
	"github.com/marni/goigc"
	"github.com/mitchellh/hashstructure"


)

// IGCTrack struct for the track data
type IGCTrack struct {
	HDate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	ID          uint64    `json:"-"`
	TrackLength float64   `json:"track_length"`
	Data        igc.Track `json:"-"`
}

// store boot time
var startTime = time.Now()
var igcTracks []IGCTrack


func main() {

	root := "/igcinfo"

	http.HandleFunc(root+"/api", apiHandler)
	http.HandleFunc(root+"/api/igc", igcHandler)

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
		for _, track := range igcTracks { ids = append(ids, strconv.Itoa(int(track.ID))) }
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ids)

	// Case POST:
	case http.MethodPost:
		/*var igcLink struct { URL string }
		err := json.NewDecoder(r.Body).Decode(&igcLink)


		if igcLink.URL == "" {
			http.Error(w, "Request missing the URL", 400)
			return
		}

		igcData, err := igc.ParseLocation(igcLink.URL)

		if err != nil {
			fmt.Fprint(w, err)
			http.Error(w, "Problem reading the track", 400)
			return
		}*/

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		var payload map[string] interface{}
		json.Unmarshal(body, &payload)

		if value, exists := payload["url"]; exists {
			igcData, err := igc. ParseLocation(value.(string))

			if err == nil {
				trackID, err := hashstructure.Hash(igcData, nil)
				if err != nil {
					http.Error(w, "Problem generating checksum", 400)
				}

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

				trackIDStr := strconv.FormatUint(trackID, 10)
				json.NewEncoder(w).Encode(IGCid{ID: trackIDStr})
				fmt.Fprint(w, "POST")
			}

		} else {
			http.Error(w, "Request missing the URL", 400)
			return
		}

		//igcData, err := igc.ParseLocation(payload["url"].(string))

		if err != nil {
			fmt.Fprint(w, err)
			http.Error(w, "Problem reading the track", 400)
			return
		}




		default:
		fmt.Fprint(w, "Error message")
	}
}

// search through the tracks for specific ID and return that, else return an error
func findTrackWithID(ID uint64) (IGCTrack, error) {

	var t IGCTrack
	for i := 0; i < len(igcTracks); i++ {
		if igcTracks[i].ID == ID {
			return igcTracks[i], nil
		}
	}
	return t, errors.New("Track not found")
}

// check if track already exists in the memory
func trackExist(trackID uint64) bool {

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

