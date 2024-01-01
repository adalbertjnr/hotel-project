package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adalbertjnr/hotel-project/api"
	"github.com/adalbertjnr/hotel-project/db"
	"github.com/adalbertjnr/hotel-project/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStorer
	hotelStore   db.HotelStorer
	userStore    db.UserStorer
	bookingStore db.BookingStorer
	ctx          = context.Background()
)

func seedUser(fname, lname, email, password string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("seeding user -> %s %s\n", user.Email, api.CreateTokenFromUser(user, "5as4c56as4d654as569C8AS908"))
	return insertedUser
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	rooms := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.InsertRoom(context.TODO(), rooms)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("seeding room ->", insertedRoom)
	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, numPersons int, fromDate, untilDate time.Time, canceled bool) {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: numPersons,
		FromDate:   fromDate,
		Tilldate:   untilDate,
		Canceled:   canceled,
	}

	if _, err := bookingStore.InsertBooking(context.TODO(), booking); err != nil {
		log.Fatal(err)
	}
	fmt.Println("seed booking ->", booking)
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("seeding hotel ->", insertedHotel)
	return insertedHotel
}

func main() {
	james := seedUser("james", "foo", "james@foo.com", "password", false)
	seedUser("admin", "admin", "admin@foo.com", "adminpassword", true)
	seedHotel("Bellucia", "France", 3)
	seedHotel("The cozy hotel", "Netherlands", 4)
	gH := seedHotel("5 Star hotel", "Germany", 5)
	seedRoom("small", true, 89.99, gH.ID)
	seedRoom("medium", true, 189.99, gH.ID)
	room := seedRoom("large", true, 289.99, gH.ID)
	seedBooking(james.ID, room.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2), false)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME)
	bookingStore = db.NewMongoBookingStore(client)
}
