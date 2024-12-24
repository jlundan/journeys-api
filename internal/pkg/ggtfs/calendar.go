package ggtfs

import (
	"strconv"
)

// CalendarItem represents the GTFS calendar file structure.
type CalendarItem struct {
	ServiceId  ID                      // service_id
	Monday     AvailableForWeekdayInfo // monday
	Tuesday    AvailableForWeekdayInfo // tuesday
	Wednesday  AvailableForWeekdayInfo // wednesday
	Thursday   AvailableForWeekdayInfo // thursday
	Friday     AvailableForWeekdayInfo // friday
	Saturday   AvailableForWeekdayInfo // saturday
	Sunday     AvailableForWeekdayInfo // sunday
	StartDate  Date                    // start_date
	EndDate    Date                    // end_date
	LineNumber int                     // CSV row number
}

func (c CalendarItem) Validate() []error {
	var validationErrors []error

	fields := []struct {
		fieldName string
		field     ValidAndPresentField
	}{
		{"service_id", &c.ServiceId},
		{"monday", &c.Monday},
		{"tuesday", &c.Tuesday},
		{"wednesday", &c.Wednesday},
		{"thursday", &c.Thursday},
		{"friday", &c.Friday},
		{"saturday", &c.Saturday},
		{"sunday", &c.Sunday},
		{"start_date", &c.StartDate},
		{"end_date", &c.EndDate},
	}

	for _, f := range fields {
		validationErrors = append(validationErrors, validateFieldIsPresentAndValid(f.field, f.fieldName, c.LineNumber, CalendarFileName)...)
	}

	return validationErrors
}

// CreateCalendarItem creates a CalendarItem from a CSV row, using the provided headers.
func CreateCalendarItem(row []string, headers map[string]int, lineNumber int) *CalendarItem {
	calendarItem := &CalendarItem{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)
		switch hName {
		case "service_id":
			calendarItem.ServiceId = NewID(v)
		case "monday":
			calendarItem.Monday = NewAvailableForWeekdayInfo(v)
		case "tuesday":
			calendarItem.Tuesday = NewAvailableForWeekdayInfo(v)
		case "wednesday":
			calendarItem.Wednesday = NewAvailableForWeekdayInfo(v)
		case "thursday":
			calendarItem.Thursday = NewAvailableForWeekdayInfo(v)
		case "friday":
			calendarItem.Friday = NewAvailableForWeekdayInfo(v)
		case "saturday":
			calendarItem.Saturday = NewAvailableForWeekdayInfo(v)
		case "sunday":
			calendarItem.Sunday = NewAvailableForWeekdayInfo(v)
		case "start_date":
			calendarItem.StartDate = NewDate(v)
		case "end_date":
			calendarItem.EndDate = NewDate(v)
		}
	}

	return calendarItem
}

func ValidateCalendarItems(calendarItems []*CalendarItem) ([]error, []string) {
	var validationErrors []error

	for _, calendarItem := range calendarItems {
		validationErrors = append(validationErrors, calendarItem.Validate()...)
	}

	return validationErrors, nil
}

const (
	CalendarAvailableForWeekday    string = "1"
	CalendarNotAvailableForWeekday string = "0"
)

type AvailableForWeekdayInfo struct {
	Integer
}

func (w AvailableForWeekdayInfo) IsValid() bool {
	val, err := strconv.Atoi(w.Integer.base.raw)

	if err != nil {
		return false
	}

	return val == 0 || val == 1
}

func NewAvailableForWeekdayInfo(raw *string) AvailableForWeekdayInfo {
	if raw == nil {
		return AvailableForWeekdayInfo{
			Integer{base: base{raw: ""}}}
	}
	return AvailableForWeekdayInfo{Integer{base: base{raw: *raw, isPresent: true}}}
}
