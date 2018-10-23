package database

import (
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Operations
type Operations interface {
	Count() int
	Get(int) (Track, bool)
	GetAll() ([]Track, bool)
}

// MongoDB the db object
type MongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

// URL used to parse the url post request
type URL struct {
	Url string `json:"url"`
}

// Information used in the index handler
type Information struct {
	Uptime  string
	Info    string
	Version string
}

// Track igc track record
type Track struct {
	//Id          bson.ObjectId `bson:"_id,omitempty"`
	TrackID       int       `json:"trackid"`
	HDate         time.Time `json:"H_date"`
	Pilot         string    `json:"pilot"`
	Glider        string    `json:"glider"`
	GliderId      string    `json:"glider_id`
	TrackLength   float64   `json:"track_length`
	Track_src_url string    `json:"track_src_url"`
	Timestamp     time.Time
}

// Ticker Flytte denne?
type Ticker struct {
	TLatest time.Time
	TStart  time.Time
	TStop   time.Time
	Tracks  []int
	Process time.Duration
}

/*
db init
*/
func (db *MongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)

	if err != nil {
		log.Fatalf("Can't connect to mongodb")
	}

	defer session.Close()

	index := mgo.Index{
		Key:        []string{"TrackID"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.DatabaseName).C(db.CollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

/*
adds track to db
*/
func (db *MongoDB) Add(t Track) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(t)
	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return err
	}

	return nil
}

/*
Returns the nr of tracks in db
*/
func (db *MongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		log.Fatalf("Can't connect to mongodb")
	}

	defer session.Close()

	// handle to "db"
	count, err := session.DB(db.DatabaseName).C(db.CollectionName).Count()
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

/*
Get returns a track with a given ID or empty track struct.
*/
func (db *MongoDB) Get(keyID int) (Track, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	Track := Track{}
	ok := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"trackid": keyID}).One(&Track)
	if err != nil {
		ok = false
	}
	return Track, ok
}

/*
Gets all the tracks from the db
*/
func (db *MongoDB) GetAll() ([]Track, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []Track
	ok := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		ok = false
	}

	return all, ok
}

// DeleteAll deletes all track records from db
func (db *MongoDB) DeleteAll() error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).DropCollection()
	if err != nil {
		return err
	}
	return nil
}

/** 	Ticker db- operations		**/

/*
// Adds ticker into db
func (db * MongoDB) Add(t Ticker) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(t)
	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
		return err
	}

	return nil
}


// Get gets and returns a ticker with a given ID or empty track struct.

func (db *MongoDB) Get(keyID int) (Track, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	Track := Track{}
	ok := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"trackid": keyID}).One(&Track)
	if err != nil {
		ok = false
	}
	return Track, ok
}
*/
