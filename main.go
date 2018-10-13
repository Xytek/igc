package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var lastUsed int
var startTime = time.Now()

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	db := TrackDB{}
	db.Init()

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/igcinfo/api", handlerAPI)
	http.HandleFunc("/igcinfo/api/igc", handlerAPIIGC)
	http.HandleFunc("/igcinfo/api/igc/", handlerAPIIGCMORE)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
