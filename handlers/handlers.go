package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/marni/goigc"
	"net/http"
	"paragliding/database"
	"paragliding/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*   Global vars   */
var Start time.Time
var id int
var GlobalDB database.MongoDB

/*
Redirects req from paragliding/ to paragliding/api/
 */
func Redirect(w http.ResponseWriter, r *http.Request){

	reg := regexp.MustCompile("^/(paragliding/)$")
	parts := reg.FindStringSubmatch(r.URL.Path)
	if parts != nil {
		http.Redirect(w, r, r.URL.Path +"/api/", 301)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

/*
** Basepoint of the API. Gives basic info about the API
 */
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reg := regexp.MustCompile("^/paragliding/(api/)$")
	parts := reg.FindStringSubmatch(r.URL.Path)

	uptime := model.GetUptime(Start)

	if parts != nil {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(database.Information{
				Uptime: uptime, Info: "Service for IGC tracks.", Version: "version 1.0",
			})
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

}

/*
** Accepts POST or GET request
** Restores a track when the right igc- url is sent with POST
** Shows slices of IDs of tracks restored in the memory when GET are used
 */
func RegAndShowTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch r.Method {
	case "POST":
		var url database.URL
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			http.Error(w, "Malformed igc- url. The api only accepts JSON", 404)
		}

		track, err := igc.ParseLocation(url.Url)
		if err != nil {
			http.Error(w, "Empty or wrong igc url provided", 404)
		} else {
			var trackLen float64
			for i := 0; i < len(track.Points)-1; i++ {
				trackLen += track.Points[i].Distance(track.Points[i+1])
			}

			if err != nil {
				http.Error(w, "Error connecting to db", 404)
			}
			count := GlobalDB.Count()
			id = count + 1

			newTrack := database.Track{
				TrackID: id,
				HDate: track.Date,
				Pilot: track.Pilot, Glider: track.GliderType,
				GliderId: track.GliderID, TrackLength: trackLen,
				Track_src_url: url.Url,
			}
			err = GlobalDB.Add(newTrack)
			if err != nil {
				http.Error(w, "Error inserting track into db", 404)
			}
			json.NewEncoder(w).Encode(id)
		}

	case "GET":
		allTracks, ok := GlobalDB.GetAll()
		if !ok {
			http.Error(w, "Error getting all the tracks", 404)
		}
		fmt.Println(GlobalDB.Count())
		var allIDs = make([]int, GlobalDB.Count())
		for index, tr:= range allTracks {
			allIDs[index] = tr.TrackID
		}
		json.NewEncoder(w).Encode(allIDs)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

/*
** Retrieves a track by its id, Accepts only GET
 */

func ShowTrackInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var field string
	path := strings.Split(r.URL.Path, "/")
	//fmt.Println(len(path), path,  "\n")
	if len(path) < 3 || len(path) > 6 {
		http.Error(w, "Too many args", http.StatusNotFound)
		return
	}
	ID, conErr := strconv.Atoi(path[4])
	if conErr != nil {
		http.Error(w, "Wrong or empty id provided!", http.StatusNotFound)
		return
	}
	if len(path) == 6 {
		field = path[5]
	}

	if r.Method == "GET" {
		track := database.Track{}
		track, ok := GlobalDB.Get(ID)
		if !ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if field == "" {
			json.NewEncoder(w).Encode(track)
		} else if field != "" {
			ShowTrackField(w, r, track, field)
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

/*
** Retrieves the track field, Accepts only GET
*/
func ShowTrackField(w http.ResponseWriter, r *http.Request, obj database.Track, field string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	switch field {
	case "pilot":
		fmt.Fprint(w, obj.Pilot)
	case "glider":
		fmt.Fprint(w, obj.Glider)
	case "glider_id":
		fmt.Fprint(w, obj.GliderId)
	case "calculated total track length":
		fmt.Fprint(w, obj.TrackLength)
	case "H_date":
		fmt.Fprint(w, obj.HDate)
	default:
		http.Error(w, "Wrong field provided", http.StatusNotFound)
	}
}
