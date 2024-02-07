package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/adalbertjnr/hotel-project/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		user    = fixtures.AddUser(&db.Store, "james", "foo", false)
		hotel   = fixtures.AddHotel(&db.Store, "foo hotel", "a", 4, nil)
		room    = fixtures.AddRoom(&db.Store, "small", true, 4.4, hotel.ID)
		from    = time.Now()
		until   = from.AddDate(0, 0, 5)
		booking = fixtures.AddBooking(&db.Store, user.ID, room.ID, from, until)
	)
	fmt.Println(booking)

}
