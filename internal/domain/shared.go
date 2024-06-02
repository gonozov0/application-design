package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	OrderID uuid.UUID
	UserID  uuid.UUID
	HotelID string
	RoomID  string
)

func ValidateBookingDates(from, to time.Time, err error) error {
	if from.IsZero() {
		return fmt.Errorf("%w: \"from\" cannot be empty", err)
	}
	if to.IsZero() {
		return fmt.Errorf("%w: \"to\" cannot be empty", err)
	}
	if from.After(to) {
		return fmt.Errorf("%w: \"from\" cannot be after \"to\"", err)
	}
	if !isDate(from) {
		return fmt.Errorf("%w: \"from\" must be a date without time", err)
	}
	if !isDate(to) {
		return fmt.Errorf("%w: \"to\" must be a date without time", err)
	}
	return nil
}

func isDate(t time.Time) bool {
	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
}
