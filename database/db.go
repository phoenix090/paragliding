package database

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"time"
)

/* Global vars */
//var TrackRec map[int]Track

type MongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

type URL struct {
	Url string `json:"url"`
}

type Information struct {
	Uptime  string
	Info    string
	Version string
}

type Track struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	TrackID     int
	HDate       time.Time
	Pilot       string
	Glider      string
	GliderId    string
	TrackLength float64
	Track_src_url string
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
	allWasGood := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"trackid": keyID}).One(&Track)
	if err != nil {
		allWasGood = false
	}
	return Track, allWasGood
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