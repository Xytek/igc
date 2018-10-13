package main

import (
	"net/http"
	"time"
)

var lastUsed int
var startTime = time.Now()

func main() {
	db := TrackDB{}
	db.Init()

	http.HandleFunc("/igcinfo/api", handlerAPI)
	http.HandleFunc("/igcinfo/api/igc", handlerAPIIGC)
	http.HandleFunc("/igcinfo/api/igc/", handlerAPIIGCMORE)
	http.ListenAndServe(":8080", nil)

}
