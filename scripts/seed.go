package main

import (
	"context"
	"fmt"
	"log"

	"github.com/souzagmu/hotel-project/db"
	"github.com/souzagmu/hotel-project/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name, location string) error {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 122.9,
		},
	}

	// Insere o hotel no banco e retorna o ID da inserção no banco
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	// Utiliza o ID gerado do hotel acima e o adiciona para cada quarto por causa da estrutura do quarto que possui o ID do hotel
	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("insertedRoom: %v", insertedRoom)
	}
	return nil
}

func main() {
	seedHotel("Bellucia", "France")
	seedHotel("The cozy hotel", "Netherlands")
	seedHotel("5 Star hotel", "Germany")
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
}
