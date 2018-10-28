package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	DBInit()
	IDsInit()
	r := mux.NewRouter()
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	// Part 1, meta logic
	r.HandleFunc("/paragliding/", apiRedirect)
	r.HandleFunc("/paragliding/api", apiGet).Methods("GET")

	// Part 2, track logic
	r.HandleFunc("/paragliding/api/track", trackGet).Methods("GET")
	r.HandleFunc("/paragliding/api/track", trackPost).Methods("POST")
	r.HandleFunc("/paragliding/api/track/{igcId}", trackIDGet).Methods("GET")
	r.HandleFunc("/paragliding/api/track/{igcId}/{igcField}", trackIDFieldGet).Methods("GET")

	// Part 3, ticker logic
	r.HandleFunc("/paragliding/api/ticker", tickerGet).Methods("GET")
	r.HandleFunc("/paragliding/api/ticker/latest", tickerLatestGet).Methods("GET")
	r.HandleFunc("/paragliding/api/ticker/{timestamp}", tickerTimestampGet).Methods("GET")

	// Part 4, webhook logic
	r.HandleFunc("/paragliding/api/webhook/new_track", webhookNewTrackGet).Methods("POST")
	r.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}", webhookNewTrackIDGet).Methods("GET")
	r.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}", webhookNewTrackIDDelete).Methods("DELETE")

	// Part 5, admin logic
	r.HandleFunc("/paragliding/admin/api/tracks_count", adminCount).Methods("GET")
	r.HandleFunc("/paragliding/admin/api/tracks", adminDelete).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(addr, r))
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
