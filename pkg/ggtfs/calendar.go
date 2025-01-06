package ggtfs

type CalendarItem struct {
	ServiceId  *string // service_id 	(required)
	Monday     *string // monday		(required)
	Tuesday    *string // tuesday		(required)
	Wednesday  *string // wednesday		(required)
	Thursday   *string // thursday		(required)
	Friday     *string // friday		(required)
	Saturday   *string // saturday		(required)
	Sunday     *string // sunday		(required)
	StartDate  *string // start_date	(required)
	EndDate    *string // end_date		(required)
	LineNumber int
}

func CreateCalendarItem(row []string, headers map[string]int, lineNumber int) *CalendarItem {
	calendarItem := &CalendarItem{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)
		switch hName {
		case "service_id":
			calendarItem.ServiceId = v
		case "monday":
			calendarItem.Monday = v
		case "tuesday":
			calendarItem.Tuesday = v
		case "wednesday":
			calendarItem.Wednesday = v
		case "thursday":
			calendarItem.Thursday = v
		case "friday":
			calendarItem.Friday = v
		case "saturday":
			calendarItem.Saturday = v
		case "sunday":
			calendarItem.Sunday = v
		case "start_date":
			calendarItem.StartDate = v
		case "end_date":
			calendarItem.EndDate = v
		}
	}

	return calendarItem
}

func ValidateCalendarItem(c CalendarItem) []ValidationNotice {
	var validationResults []ValidationNotice

	fields := []struct {
		fieldType FieldType
		name      string
		value     *string
		required  bool
	}{
		{FieldTypeID, "service_id", c.ServiceId, true},
		{FieldTypeCalendarDay, "monday", c.Monday, true},
		{FieldTypeCalendarDay, "tuesday", c.Tuesday, true},
		{FieldTypeCalendarDay, "wednesday", c.Wednesday, true},
		{FieldTypeCalendarDay, "thursday", c.Thursday, true},
		{FieldTypeCalendarDay, "friday", c.Friday, true},
		{FieldTypeCalendarDay, "saturday", c.Saturday, true},
		{FieldTypeCalendarDay, "sunday", c.Sunday, true},
		{FieldTypeDate, "start_date", c.StartDate, true},
		{FieldTypeDate, "end_date", c.EndDate, true},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldType, field.name, field.value, field.required, FileNameCalendar, c.LineNumber)...)
	}

	return validationResults
}

func ValidateCalendarItems(calendarItems []*CalendarItem) []ValidationNotice {
	var results []ValidationNotice

	for _, calendarItem := range calendarItems {
		if calendarItem == nil {
			continue
		}
		results = append(results, ValidateCalendarItem(*calendarItem)...)
	}

	return results
}
