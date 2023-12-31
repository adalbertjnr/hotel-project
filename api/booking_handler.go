package api

import (
	"fmt"

	"github.com/adalbertjnr/hotel-project/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(genericResp{
			Type:    "error",
			Message: "not authorized",
		})
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), booking.ID.Hex(), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{Type: "status", Message: "updated"})
}

// admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return c.JSON(bookings)
}

// user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookingByID(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(genericResp{
			Type:    "error",
			Message: "not authorized",
		})
	}
	return c.JSON(booking)
}
