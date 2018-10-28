package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"paragliding/database"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/marni/goigc"
)

// Start time of the api
var Start time.Time

//GlobalDB ...
var GlobalDB database.MongoDB

//ticker obj
var ticker database.Ticker

// Connect gets connection and initialize global vars
func Connect() {
	Start = time.Now()

	dbURL, ok := os.LookupEnv("DBURL")
	dbName, ok2 := os.LookupEnv("DBNAME")
	dbCollection, ok3 := os.LookupEnv("DBCOLLECTION")
	if !ok || !ok2 || !ok3 {
		dbURL = "mongodb://test:test123@ds221003.mlab.com:21003/mongo_db"
		dbName = "mongo_db"
		dbCollection = "track"
	}
	GlobalDB = database.MongoDB{DatabaseURL: dbURL, DatabaseName: dbName, CollectionName: dbCollection}
	GlobalDB.Init()
}

// Redirect req from paragliding/ to paragliding/api/
func Redirect(w http.ResponseWriter, r *http.Request) {

	reg := regexp.MustCompile("^/(paragliding/)$")
	parts := reg.FindStringSubmatch(r.URL.Path)
	if parts != nil {
		http.Redirect(w, r, r.URL.Path+"/api/", 301)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

/*
** Index Basepoint of the API. Gives basic info about the API
 */
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	reg := regexp.MustCompile("^/paragliding/(api/)$")
	parts := reg.FindStringSubmatch(r.URL.Path)

	uptime := GetUptime(Start)

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
** RegAndShowTrack Accepts POST or GET request
** Restores a track when the right igc- url is sent with POST
** Shows slices of IDs of tracks restored in the memory when GET are used
 */
func RegAndShowTrack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch r.Method {
	case "POST":
		HandleTrackPost(w, r)

	case "GET":
		allTracks, ok := GlobalDB.GetAll()
		if !ok {
			http.Error(w, http.StatusText(404), 404)
		}
		var allIDs = make([]int, GlobalDB.Count())
		for index, tr := range allTracks {
			allIDs[index] = tr.TrackID
		}
		json.NewEncoder(w).Encode(allIDs)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// HandleTrackPost Handles track registration and ticker update
func HandleTrackPost(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, http.StatusText(404), 404)
		}
		count := GlobalDB.Count()
		id := count + 1

		newTrack := database.Track{
			TrackID: id,
			HDate:   track.Date,
			Pilot:   track.Pilot, Glider: track.GliderType,
			GliderId: track.GliderID, TrackLength: trackLen,
			Track_src_url: url.Url,
			Timestamp:     time.Now(),
		}
		err = GlobalDB.Add(newTrack)
		TriggerWebhook()
		if err != nil {
			http.Error(w, http.StatusText(404), 404)
		}
		json.NewEncoder(w).Encode(id)
	}
}

/*
** ShowTrackInfo Retrieves a track by its id, Accepts only GET
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
** ShowTrackField Retrieves the track field, Accepts only GET
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
	case "track_length":
		fmt.Fprint(w, obj.TrackLength)
	case "H_date":
		fmt.Fprint(w, obj.HDate)
	case "track_src_url":
		fmt.Fprint(w, obj.Track_src_url)
	default:
		http.Error(w, "Wrong field provided", http.StatusNotFound)
	}
}

/****		TICKER HANDLERS		 ***/

//GetLatestTicker the latest timestamp
func GetLatestTicker(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
	process := time.Now()
	PopulateTickerInfo(3, false, "")
	t := time.Now()
	ticker.Process = t.Sub(process)
	// eller after := time.Since(process).Seconds() / 1000
	fmt.Fprint(w, ticker.TLatest)
}

// GetTickerInfo getst ticker with timestamps with cap if omitted from the user as GET- var
func GetTickerInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	// FIKS PATH ERROR
	// parts := strings.Split(r.URL.Path, "/")

	if r.Method == "GET" {
		reg := regexp.MustCompile("^/api/ticker/([A-Z. a-z0-9:-]*)$")
		parts := reg.FindStringSubmatch(r.URL.Path)

		process := time.Now()
		var default_cap = 5
		cap, ok := r.URL.Query()["cap"]
		if ok {
			if cap[0] != "" {
				c, err := strconv.Atoi(cap[0])
				if err != nil || c > 0 {
					default_cap = c
					fmt.Println(default_cap)
				}
			}
		}

		if parts == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		//Goes to endpoint /api/ticker/<timestamp> when conditions are met
		if len(parts) == 2 && parts[1] != "" {
			err := PopulateTickerInfo(default_cap, true, parts[1])
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}
			t := time.Now()
			ticker.Process = t.Sub(process)
			json.NewEncoder(w).Encode(ticker)
			return
		}

		err := PopulateTickerInfo(default_cap, false, "")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		t := time.Now()
		ticker.Process = t.Sub(process)
		json.NewEncoder(w).Encode(ticker)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// PopulateTickerInfo makes ticker obj dynamically
// the bool which desides if ticker is maid for /api/ticker/<timestamp>, if not it will be false
func PopulateTickerInfo(cap int, which bool, param string) error {
	totTracks, ok := GlobalDB.GetAll()
	ticker = database.Ticker{}
	if !ok {
		// FIKS error handling
	}

	// Handles cap
	for idx, track := range totTracks {
		if idx+1 <= cap {
			if idx == 0 {
				ticker.TStart = track.Timestamp
			}
			if idx+1 == cap {
				ticker.TStop = track.Timestamp
			}
			// Populating []id for tracks accourding to /api/ticker/<timestamp>
			if which {
				match, err := time.Parse(time.RFC3339, param)
				if err != nil {
					return err
				}
				// Checking if the tracks timestamp is higher then the params
				if match.After(track.Timestamp) {
					ticker.Tracks = append(ticker.Tracks, track.TrackID)
					//fmt.Println("Puttet inn track med nr", track.TrackID)
				}
			} else {
				// And now for /api/ticker/
				ticker.Tracks = append(ticker.Tracks, track.TrackID)
			}
		}
		if idx+1 == GlobalDB.Count() {
			ticker.TLatest = track.Timestamp
		}
	}
	return nil
}

/*
** GetUptime updates uptime and formates it in ISO 8601 standard
 */
func GetUptime(t time.Time) (uptime string) {
	now := time.Now()
	newTime := now.Sub(t)
	hours := int(newTime.Hours())
	sek := strconv.Itoa(int(newTime.Seconds()) % 36000 % 60)
	var min, hour, y, m, d string

	// checking and setting when min gets to 1 or more
	if int(newTime.Seconds())%36000 >= 60 {
		min = strconv.Itoa(int(newTime.Minutes()) % 60)
	}

	// checking and setting when hour gets to 1 or more
	if hours >= 1 {
		hour = strconv.Itoa(hours)
	}

	// Setting the days correct
	if hours > 23 {
		d = strconv.Itoa(hours / 24)
		hour = strconv.Itoa(hours % 24)
	}
	days, _ := strconv.Atoi(d)
	// Setting the month correct
	if days > 31 {
		m = strconv.Itoa(days / 31)
		d = strconv.Itoa(days % 31)

	}
	months, _ := strconv.Atoi(m)
	// Setting the year correct
	if months > 12 {
		y = strconv.Itoa(months / 12)
		m = strconv.Itoa(months % 12)
	}

	uptime = "P" + y + "Y" + m + "M" + d + "DT" + hour + "H" + min + "M" + sek + "S"

	return uptime
}
