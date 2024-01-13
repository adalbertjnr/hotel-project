package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adalbertjnr/hotel-project/api"
	"github.com/adalbertjnr/hotel-project/db"
	"github.com/adalbertjnr/hotel-project/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("dropping db")

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	store := &db.Store{
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Booking: db.NewMongoBookingStore(client, db.DBNAME),
		Room:    db.NewMongoRoomStore(client, hotelStore, db.DBNAME),
		Hotel:   db.NewMongoHotelStore(client, db.DBNAME),
	}

	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println("james ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "fixtures hotel", "fixhotel", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 95.55, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking.ID)
}
