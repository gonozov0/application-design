package internal

import (
	"fmt"
	"time"
)

func ValidateBookingDates(from, to time.Time, err error) error {
	now := time.Now()
	if now.After(from) {
		return fmt.Errorf("%w: \"from\" cannot be in the past", err)
	}
	if now.After(to) {
		return fmt.Errorf("%w: \"to\" cannot be in the past", err)
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
