package ggtfs

import (
	"encoding/csv"
	"strconv"
)

// CalendarItem represents the GTFS calendar file structure.
type CalendarItem struct {
	ServiceId  ID          // service_id
	Monday     WeekdayEnum // monday
	Tuesday    WeekdayEnum // tuesday
	Wednesday  WeekdayEnum // wednesday
	Thursday   WeekdayEnum // thursday
	Friday     WeekdayEnum // friday
	Saturday   WeekdayEnum // saturday
	Sunday     WeekdayEnum // sunday
	StartDate  Date        // start_date
	EndDate    Date        // end_date
	LineNumber int         // CSV row number
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
			calendarItem.ServiceId = NewID(&row[hPos])
		case "monday":
			calendarItem.Monday = NewWeekdayEnum(&row[hPos])
		case "tuesday":
			calendarItem.Tuesday = NewWeekdayEnum(&row[hPos])
		case "wednesday":
			calendarItem.Wednesday = NewWeekdayEnum(&row[hPos])
		case "thursday":
			calendarItem.Thursday = NewWeekdayEnum(&row[hPos])
		case "friday":
			calendarItem.Friday = NewWeekdayEnum(&row[hPos])
		case "saturday":
			calendarItem.Saturday = NewWeekdayEnum(&row[hPos])
		case "sunday":
			calendarItem.Sunday = NewWeekdayEnum(&row[hPos])
		case "start_date":
			calendarItem.StartDate = NewDate(&row[hPos])
		case "end_date":
			calendarItem.EndDate = NewDate(&row[hPos])
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

type WeekdayEnum struct {
	Integer
}

func (w WeekdayEnum) IsValid() bool {
	val, err := strconv.Atoi(w.Integer.base.raw)
	if err != nil {
		return false
	}

	return val == 0 || val == 1
}

func NewWeekdayEnum(raw *string) WeekdayEnum {
	if raw == nil {
		return WeekdayEnum{
			Integer{base: base{raw: ""}}}
	}
	return WeekdayEnum{Integer{base: base{raw: *raw, isPresent: true}}}
}
