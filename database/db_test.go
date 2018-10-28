package database

import (
	"testing"
)

func TestSomething(t *testing.T) {

}

/*

func SetupDB(t *testing.T) *MongoDB {
	db := MongoDB{DatabaseURL: config.DBURL, DatabaseName: config.AuthDatabase, CollectionName: config.TESTCOL}
	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()
	if err != nil {
		t.Errorf("Error setting up connection %v", err)
	}

	return &db
}

// Tearing down the db after the test is finished

/*
func Teardown(t *testing.T, db *MongoDB) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		t.Error(err)
	}
	err = session.DB(db.DatabaseName).DropDatabase()
	if err != nil {
		t.Error(err)
	}
}

// Testing the connection to the db.
func TestConnection(t *testing.T) {
	// Testing with db- creds from config
	db := MongoDB{DatabaseURL: config.DBURL, DatabaseName: config.AuthDatabase, CollectionName: config.TESTCOL}
	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()
	defer Teardown(db)
	db.Init()

	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()
	if err != nil {
		t.Errorf("Error dialing db, err: %v", err)
	}

	//mocking it to fail
	session, err = mgo.Dial(db.DatabaseURL + "rubbish")
	if err == nil {
		t.Errorf("Expected error got: %v", err)
	}
}

// Adding track to db and testing if it worked
func TestAddingTrackToDB(t *testing.T) {
	db := MongoDB{DatabaseURL: config.DBURL, DatabaseName: config.AuthDatabase, CollectionName: config.TESTCOL}
	session, err := mgo.Dial(db.DatabaseURL)
	defer session.Close()
	defer Teardown(db)
	db.Init()

	//checking that the db is empthy as it should be
	if db.Count != 0 {
		t.Errorf("db was not empty, got %v", db.Count)
	}
	newTrack := Track{TrackID: 100, HDate: time.Now(), Pilot: "Carlos", Glider: "PVP", GliderId: "23", TrackLength: 234.333, Timestamp: time.Now()}
	err := db.Add(newTrack)
	if err != nil {
		t.Errorf("Error adding track into db, got %v", err)
	}

}
*/
