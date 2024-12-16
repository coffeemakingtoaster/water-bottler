package utils

import (
	"time"

	"github.com/rs/zerolog/log"
)

// Function to check if a date is in the future with generics str or date
func DateInFuture[T string | time.Time](date T) (bool, error) {
	var validUntil time.Time
	dateNow := time.Now().UTC()
	// Check of the type of the input is a string
	switch typedDate := any(date).(type) {
	case string:
		var err error
		validUntil, err = time.Parse(time.RFC3339, typedDate)
		if err != nil {
			log.Err(err).Msg("Error parsing date")
			return false, err
		}
	case time.Time:
		validUntil = typedDate
	}
	// Check if the date is in the future
	return validUntil.After(dateNow), nil
}
