package main

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DBInfo stores the details of the DB connection
type DBInfo struct {
	DBURL                 string
	DBName                string
	TrackCollectionName   string
	WebhookCollectionName string
}

// db stores the credentials of our database
var db DBInfo

// DBInit fills DBInfo with the information about our database
func DBInit() {
	db.DBURL = "mongodb://username:password123@ds143293.mlab.com:43293/paragliding"
	db.DBName = "paragliding"
	db.TrackCollectionName = "track"
	db.WebhookCollectionName = "webhook"
}

// AddTrack adds new tracks to the storage
func (db *DBInfo) addTrack(t Track) Track {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	// Adds a new ID to the ID management system
	IDs = append(IDs, lastUsed)
	t.SimpleID = lastUsed
	lastUsed++

	// Adds a timestamp to the track. The if sentence is for testing purposes
	if t.Timestamp != 10 {
		t.Timestamp = bson.NewObjectId().Time().Unix()
	}

	// Inserts the track into the database
	err = session.DB(db.DBName).C(db.TrackCollectionName).Insert(t)
	if err != nil {
		fmt.Printf("Error in AddTrack(): %v", err.Error())
	}

	// Check to see if any users need an update
	CheckUsers()

	return t
}

// CountTrack returns the current count of the tracks in storage
func (db *DBInfo) countTrack() int {
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	count, err := session.DB(db.DBName).C(db.TrackCollectionName).Count()
	if err != nil {
		fmt.Printf("Error in Count(): %v", err.Error())
		return -1
	}
	return count
}

// GetTrack returns a track with a given ID or empty struct
func (db *DBInfo) getTrack(keyID int) (Track, bool) {
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	track := Track{}
	allGood := true
	err = session.DB(db.DBName).C(db.TrackCollectionName).Find(bson.M{"simpleid": keyID}).One(&track)
	if err != nil {
		allGood = false
	}
	return track, allGood
}

// getAllTracks returns a slice of all tracks from the database
func (db *DBInfo) getAllTracks() []Track {
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var allTracks []Track
	err = session.DB(db.DBName).C(db.TrackCollectionName).Find(bson.M{}).All(&allTracks)
	if err != nil {
		fmt.Printf("Error in getAllTracks(): %v", err.Error())
	}
	return allTracks
}

// getTracksAfter returns a slice of all tracks after the attached timestamp
func (db *DBInfo) getTracksAfter(ts int64) []Track {
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	var someTracks []Track
	err = session.DB(db.DBName).C(db.TrackCollectionName).Find(bson.M{"timestamp": bson.M{"$gt": ts}}).All(&someTracks)
	if err != nil {
		fmt.Printf("Error in getTracksBetween(): %v", err.Error())
	}
	return someTracks
}

// latestTicker finds the track with the latest timestamp
func (db *DBInfo) latestTicker() int64 {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	var track Track
	err = session.DB(db.DBName).C(db.TrackCollectionName).Find(nil).Sort("-$natural").One(&track)
	if err != nil {
		fmt.Printf("Error in latestTicker(): %v", err.Error())
	}

	return track.Timestamp
}

// deleteAllTracks deletes everything from the collection in the database
func (db *DBInfo) deleteAllTracks() {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	// Reset internal ID managing systems
	lastUsed = 0
	IDs = nil
	IDs = make([]int, 0)

	// Delete everything from the collection
	session.DB(db.DBName).C(db.TrackCollectionName).RemoveAll(nil)
}

/*
-----------------------------------------------------------
--------------------------WEBHOOK--------------------------
-----------------------------------------------------------
*/

// addWebhook adds new tracks to the storage
func (db *DBInfo) addWebhook(wh Webhook) Webhook {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	// Assign an ID
	wh.ID = bson.NewObjectId()
	wh.LastCheck = wh.ID.Time().Unix()

	// Inserts the webhook into the database
	err = session.DB(db.DBName).C(db.WebhookCollectionName).Insert(wh)
	if err != nil {
		fmt.Printf("Error in AddWebhook(): %v", err.Error())
	}
	return wh
}

// getWebhook returns a webhook with a given ID or empty struct
func (db *DBInfo) getWebhook(keyID string) Webhook {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	wh := Webhook{}

	err = session.DB(db.DBName).C(db.WebhookCollectionName).Find(bson.M{"_id": bson.ObjectIdHex(keyID)}).One(&wh)
	if err != nil {
		fmt.Printf("Error in GetWebhook(): %v", err.Error())
	}
	return wh
}

// deleteWebhook deletes a specific element by ID
func (db *DBInfo) deleteWebhook(id string) {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	err = session.DB(db.DBName).C(db.WebhookCollectionName).Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	if err != nil {
		fmt.Printf("Error in DeleteWebhook(): %v", err.Error())
	}
}

// getAllWebhooks gets all the webhook elements from the database
func (db *DBInfo) getAllWebhooks() []Webhook {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	var wh []Webhook
	err = session.DB(db.DBName).C(db.WebhookCollectionName).Find(bson.M{}).All(&wh)
	if err != nil {
		fmt.Printf("Error in GetAllWebhooks(): %v", err.Error())
	}

	return wh
}

// addWebhook adds new tracks to the storage
func (db *DBInfo) updateWebhookTimestamp(wh Webhook, ts int64) {
	// Creates a connection
	session, err := mgo.Dial(db.DBURL)
	if err != nil {
		panic(err)
	}
	// Ensures the connection closes afterwards
	defer session.Close()

	// Inserts the webhook into the database
	err = session.DB(db.DBName).C(db.WebhookCollectionName).Update(bson.M{"_id": wh.ID}, bson.M{"$set": bson.M{"lastCheck": ts}})
	if err != nil {
		fmt.Printf("Error in AddWebhook(): %v", err.Error())
	}

	return
}
