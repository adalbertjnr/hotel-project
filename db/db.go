package db

const (
	COLL_ROOM  = "rooms"
	COLL_USER  = "users"
	COLL_HOTEL = "hotels"
	DBNAME     = "hotel-project"
	DBURI      = "mongodb://root:example@localhost:27017"
)

type Store struct {
	User  UserStorer
	Hotel HotelStorer
	Room  RoomStorer
}
