package daterange

import (
	"errors"
	"fmt"
	"time"
)

var (
	FromParameterInvalidFormat = errors.New("'from' parameter format is invalid. Should be YYYY-MM-DD.")
	ToParameterInvalidFormat   = errors.New("'to' parameter format is invalid. Should be YYYY-MM-DD.")
	InvalidDateRangeError      = errors.New("'from' date is not before 'to'")
)

func ConvertRangeParametersToTime(from, to string) (time.Time, time.Time, error) {
	fromTime, fromTimeErr := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", from))
	if fromTimeErr != nil {
		return time.Time{}, time.Time{}, FromParameterInvalidFormat
	}

	toTime, toTimeErr := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", to))
	if toTimeErr != nil {
		return time.Time{}, time.Time{}, ToParameterInvalidFormat
	}

	if fromTime.Unix() > toTime.Unix() {
		return time.Time{}, time.Time{}, InvalidDateRangeError
	}

	return fromTime, toTime, nil
}

func TimeToDayString(date time.Time) string {
	dateStr := date.String()
	return dateStr[0:10]
}
