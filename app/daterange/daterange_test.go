package daterange

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	dummyYear    = 2020
	dummyMonth   = 1
	dummyDay     = 2
	dummyHour    = 3
	dummyMinute  = 4
	dummySecond  = 5
	dummyNSecond = 6

	invalidDayStringFormat = "02-12-2022"
	fromYear               = 2022
	fromMonth              = 11
	fromDay                = 20
	toYear                 = 2022
	toMonth                = 11
	toDay                  = 30
)

var (
	dummyTime = time.Date(
		dummyYear,
		dummyMonth,
		dummyDay,
		dummyHour,
		dummyMinute,
		dummySecond,
		dummyNSecond,
		time.UTC)
	dummyFrom = fmt.Sprintf("%04d-%02d-%02d", fromYear, fromMonth, fromDay)
	dummyTo   = fmt.Sprintf("%04d-%02d-%02d", toYear, toMonth, toDay)
)

func TestDayStringFromTimeShallConvertTimeToDayString(t *testing.T) {
	result := TimeToDayString(dummyTime)
	assert.Equal(t, result, fmt.Sprintf("%04d-%02d-%02d", dummyYear, dummyMonth, dummyDay))
}

func TestConvertRangeParametersToTimeShallConvertTimeRangeFromDayStringsToTimesAndReturnNoErrorWhenStringFormatsAreValidAndToDayFollowsFromDay(t *testing.T) {
	from, to, err := ConvertRangeParametersToTime(dummyFrom, dummyTo)
	assert.Equal(t, from, time.Date(fromYear, fromMonth, fromDay, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, to, time.Date(toYear, toMonth, toDay, 0, 0, 0, 0, time.UTC))
	assert.Nil(t, err)
}

func TestConvertRangeParametersToTimeShallReturnFromParameterInvalidFormatWhenFromDayStringHasInvalidFormat(t *testing.T) {
	_, _, err := ConvertRangeParametersToTime(invalidDayStringFormat, dummyTo)
	assert.Equal(t, err, FromParameterInvalidFormat)
}

func TestConvertRangeParametersToTimeShallReturnToParameterInvalidFormatWhenToDayStringHasInvalidFormat(t *testing.T) {
	_, _, err := ConvertRangeParametersToTime(dummyFrom, invalidDayStringFormat)
	assert.Equal(t, err, ToParameterInvalidFormat)
}

func TestConvertRangeParametersToTimeShallReturnInvalidDateRangeErrorWhenFromDayFollowsAfterToDay(t *testing.T) {
	_, _, err := ConvertRangeParametersToTime(dummyTo, dummyFrom)
	assert.Equal(t, err, InvalidDateRangeError)
}
