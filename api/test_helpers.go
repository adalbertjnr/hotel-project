package api

import (
	"context"
	"log"
	"testing"

	"github.com/adalbertjnr/hotel-project/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb://root:example@localhost:27017"
	dbname    = "hotel-project-test"
)

type testdb struct {
	client *mongo.Client
	db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	const TEST_DATABASE = "hotel-project-test"

	hotelStore := db.NewMongoHotelStore(client, TEST_DATABASE)
	return &testdb{
		client: client,
		Store: db.Store{
			User:    db.NewMongoUserStore(client, TEST_DATABASE),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore, TEST_DATABASE),
			Booking: db.NewMongoBookingStore(client, TEST_DATABASE),
		},
	}
}
