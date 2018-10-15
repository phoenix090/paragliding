package main

import (
	"log"
	"net/http"
	"os"
	"paragliding/config"
	"paragliding/handlers"
	"time"
)

func main() {

	handlers.Start = time.Now()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	var err error
	handlers.GlobalDB, err = config.GetMongoDB()
	if err != nil {
		log.Fatalf("Error connecting to db, %v", err)
	}
	handlers.GlobalDB.Init()

	http.HandleFunc("/paragliding/", handlers.Redirect)
	http.HandleFunc("/paragliding/api/", handlers.Index)
	http.HandleFunc("/paragliding/api/track", handlers.RegAndShowTrack)
	http.HandleFunc("/paragliding/api/track/", handlers.ShowTrackInfo)
	err = http.ListenAndServe(":"+port, nil)
	log.Fatalf("Server error: %s", err)
}
