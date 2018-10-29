package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"paragliding/handlers"
	"sync"
	"time"
)

const discordUrl = "https://discordapp.com/api/webhooks/503581701981732865/OCYiZbwzXY0LlwCZ6J8DhxuD6PFCja7PNC08Du9B0SfcNTR-LzLvaqJit5FJfbbL0Aod"

// Uses the discord webhook above to send info about new tracks registered.
func sendDiscord() {
	count := 0
	var err error
	for {
		delay := time.Minute * 10
		count, err = handlers.SendLogToDiscord(discordUrl, &count)
		if err != nil {
			fmt.Println("Exiting, something went wrong")
		}
		time.Sleep(delay)
	}
}

// For Heroku
func main() {
	handlers.Connect()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		sendDiscord()
		wg.Done()
	}()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
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

	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("Server error: %s", err)

	wg.Wait()
}
