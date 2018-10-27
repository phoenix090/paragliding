package main

import (
	"log"
	"net/http"
	"os"
	"paragliding/admin"
	"paragliding/config"
	"paragliding/database"
	"paragliding/handlers"
	"time"
)

func main() {
	handlers.Start = time.Now()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	//fmt.Println(handlers.SendDiscord("..."))
	var err error
	dbURL, ok := os.LookupEnv("DBURL")
	dbName, ok2 := os.LookupEnv("AuthDatabase")
	dbCollection, ok3 := os.LookupEnv("DBCollection")
	if !ok || !ok2 || !ok3 {
		handlers.GlobalDB, err = config.GetMongoDB()
	} else {
		handlers.GlobalDB = database.MongoDB{DatabaseURL: dbURL, DatabaseName: dbName, CollectionName: dbCollection}
	}

	if err != nil {
		log.Fatalf("Error connecting to db, %v", err)
	}
	handlers.GlobalDB.Init()

	http.HandleFunc("/paragliding/", handlers.Redirect)
	http.HandleFunc("/paragliding/api/", handlers.Index)
	http.HandleFunc("/paragliding/api/track", handlers.RegAndShowTrack)
	http.HandleFunc("/paragliding/api/track/", handlers.ShowTrackInfo)

	// Ticker handlers
	http.HandleFunc("/api/ticker/", handlers.GetTickerInfo)
	http.HandleFunc("/api/ticker/latest", handlers.GetLatestTicker)

	//	Admin handlers
	http.HandleFunc("/admin/api/tracks_count", admin.GetTracksCount)
	http.HandleFunc("/admin/api/tracks", admin.DeleteAllTracks)

	//Webhooks  /api/webhook/new_track/
	http.HandleFunc("/api/webhook/new_track/", handlers.RegisterWebhook)

	err = http.ListenAndServe(":"+port, nil)
	log.Fatalf("Server error: %s", err)
	// Checks every then minutes and sends webhook notification to subs

	for {
		handlers.NotifySubs()
		delay := time.Minute * 10
		time.Sleep(delay)
	}
}
