package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Content string `json:"content"`
}

func main() {
	// Holds the track counts
	countOld := 0
	countNew := 0
	url := "https://discordapp.com/api/webhooks/506087941764939790/4CuwPYPKQy1UylbdbOuf2fgcoL5Dmtk7-W7D3yNrQQ9P_bxxKUL0dT5AjrNI81i_b8yd"

	// Start the infinite loop
	for true {
		// To calculate processing for the message
		funcStart := time.Now()

		// Get the tracks
		resp, err := http.Get("https://igcinf.herokuapp.com/paragliding/api/track")
		if err != nil {
			fmt.Printf("Error making the GET request, %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
			return
		}

		// Move the last rounds new tracks to this rounds old
		countOld = countNew

		// Get the track array
		var a []int
		err = json.NewDecoder(resp.Body).Decode(&a)

		// Update new track count
		countNew = len(a)

		// Get an array with all the new tracks since the last round for the message
		var newTracks []int
		for i := countNew; i > countOld; i-- {
			newTracks = append(newTracks, a[i-1])
		}

		// If there's any new tracks since the last check
		if countNew > countOld {
			// Get the latest ticker for the message
			resp, err := http.Get("https://igcinf.herokuapp.com/paragliding/api/ticker/latest")
			if err != nil {
				fmt.Printf("Error making the GET request, %s", err)
			}
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
				return
			}

			var latestTicker int64
			err = json.NewDecoder(resp.Body).Decode(&latestTicker)

			// Creates a string in the format 1, 2, 3 from the id array
			tracksArrayPresentable := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(newTracks)), ", "), "[]")

			// Checks the processing time since the start of the function (FindUsers) in ms
			processing := time.Since(funcStart) / 1000000

			var message Message
			// Creates the string that'll be shown to the user
			message.Content = fmt.Sprintf("Latest timestamp: %v, %v new tracks are: %s. (processing: %vs %vms)",
				latestTicker, len(newTracks), tracksArrayPresentable, int(processing/1000), int(processing%1000))

			msg, _ := json.Marshal(message)

			// POST the message to the webhook
			resp, err = http.Post(url, "application/json", bytes.NewBuffer(msg))
			if err != nil {
				panic(err)
			}

		}

		// Sleep for 10 minutes before doing the loop again
		time.Sleep(10 * time.Minute)
	}

}
