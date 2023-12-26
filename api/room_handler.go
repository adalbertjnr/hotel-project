package api

import (
	"fmt"
	"time"

	"github.com/adalbertjnr/hotel-project/db"
	"github.com/adalbertjnr/hotel-project/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("error booking the room. time must be bigger than current time")
	}
	return nil
}

type RoomHandler struct {
	store db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: *store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomOID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(genericResp{
			Type:    "error",
			Message: "internal server error",
		})
	}

	where := bson.M{
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(c.Context(), where)
	if err != nil {
		return err
	}

	if len(bookings) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(genericResp{
			Type:    "error",
			Message: fmt.Sprintf("room %s already booked", roomOID),
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomOID,
		FromDate:   params.FromDate,
		Tilldate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(genericResp{
			Type:    "error",
			Message: "internal server error",
		})
	}
	return c.JSON(insertedBooking)
}
