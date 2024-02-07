package db

const (
	COLL_ROOM     = "rooms"
	COLL_USER     = "users"
	COLL_HOTEL    = "hotels"
	COLL_BOOKINGS = "bookings"
	DBURI         = "mongodb://root:example@localhost:27017"
	DBNAME        = "hotel-project"
)

type Store struct {
	User    UserStorer
	Hotel   HotelStorer
	Room    RoomStorer
	Booking BookingStorer
}
