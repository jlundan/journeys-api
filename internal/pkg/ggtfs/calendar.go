package ggtfs

import (
	"encoding/csv"
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

var validCalendarHeaders = []string{
	"service_id", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday", "start_date", "end_date",
}

// LoadCalendarItems reads and parses GTFS calendar data from the provided CSV reader.
func LoadCalendarItems(csvReader *csv.Reader) ([]*CalendarItem, []error) {
	entities, errs := loadEntities(csvReader, validCalendarHeaders, CreateCalendarItem, CalendarFileName)

	calendarItems := make([]*CalendarItem, 0, len(entities))
	for _, entity := range entities {
		if calendarItem, ok := entity.(*CalendarItem); ok {
			calendarItems = append(calendarItems, calendarItem)
		}
	}

	return calendarItems, errs
}

// CreateCalendarItem creates a CalendarItem from a CSV row, using the provided headers.
func CreateCalendarItem(row []string, headers map[string]int, lineNumber int) interface{} {
	calendarItem := &CalendarItem{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "service_id":
			calendarItem.ServiceId = NewID(getRowValue(row, hPos))
		case "monday":
			calendarItem.Monday = NewAvailableForWeekdayInfo(getRowValue(row, hPos))
		case "tuesday":
			calendarItem.Tuesday = NewAvailableForWeekdayInfo(getRowValue(row, hPos))
		case "wednesday":
			calendarItem.Wednesday = NewAvailableForWeekdayInfo(getRowValue(row, hPos))
		case "thursday":
			calendarItem.Thursday = NewAvailableForWeekdayInfo(getRowValue(row, hPos))
		case "friday":
			calendarItem.Friday = NewAvailableForWeekdayInfo(getRowValue(row, hPos))
		case "saturday":
			calendarItem.Saturday = NewAvailableForWeekdayInfo(getRowValue(row, hPos))
		case "sunday":
			calendarItem.Sunday = NewAvailableForWeekdayInfo(getRowValue(row, hPos))
		case "start_date":
			calendarItem.StartDate = NewDate(getRowValue(row, hPos))
		case "end_date":
			calendarItem.EndDate = NewDate(getRowValue(row, hPos))
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
