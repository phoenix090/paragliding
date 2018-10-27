package admin

import (
	"fmt"
	"net/http"
	"paragliding/handlers"
)

/** Admin operations, exposed only to admim */

// GetTracksCount gets the current track count in the db
func GetTracksCount(w http.ResponseWriter, r *http.Request) {
	// simple auth for admin users
	code, ok := r.URL.Query()["admincode"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if code[0] != "12345" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
	if r.Method == "GET" {
		fmt.Fprint(w, handlers.GlobalDB.Count())
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// DeleteAllTracks deletes all tracks from db
func DeleteAllTracks(w http.ResponseWriter, r *http.Request) {
	// simple auth for admin users
	code, ok := r.URL.Query()["admincode"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if code[0] != "12345" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
	if r.Method == "DELETE" {
		if handlers.GlobalDB.Count() == 0 {
			fmt.Fprint(w, "Track records are empty...")
		} else {
			err := handlers.GlobalDB.DeleteAll()
			if err != nil {
				fmt.Fprint(w, "Error deleting all the records")
				return
			}
			fmt.Fprint(w, "OK!")
		}
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
