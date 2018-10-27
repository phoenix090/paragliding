package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"paragliding/handlers"
)

// For openstack
func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	err := handlers.Connect()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	http.HandleFunc("/paragliding/", handlers.Redirect)
	http.HandleFunc("/paragliding/api/", handlers.Index)
	http.HandleFunc("/paragliding/api/track", handlers.RegAndShowTrack)
	http.HandleFunc("/paragliding/api/track/", handlers.ShowTrackInfo)

	// Ticker handlers
	http.HandleFunc("/api/ticker/", handlers.GetTickerInfo)
	http.HandleFunc("/api/ticker/latest", handlers.GetLatestTicker)

	//	Admin handlers
	http.HandleFunc("/admin/api/tracks_count", handlers.GetTracksCount)
	http.HandleFunc("/admin/api/tracks", handlers.DeleteAllTracks)

	//Webhooks  /api/webhook/new_track/
	http.HandleFunc("/api/webhook/new_track/", handlers.RegisterWebhook)

	err = http.ListenAndServe(":"+port, nil)
	log.Fatalf("Server error: %s", err)
	// Checks every then minutes and sends webhook notification to subs

	fmt.Println("Hello from openstack sin main")
	/*
		for {
			handlers.NotifySubs()
			delay := time.Minute * 10
			time.Sleep(delay)
		}
	*/
}
