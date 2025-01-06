package ggtfs

type CalendarDate struct {
	ServiceId     *string // service_id 	(required)
	Date          *string // date 			(required)
	ExceptionType *string // exception_type (required)
	LineNumber    int
}

func CreateCalendarDate(row []string, headers map[string]int, lineNumber int) *CalendarDate {
	calendarDate := &CalendarDate{
		LineNumber: lineNumber,
	}

	for hName := range headers {
		v := getRowValueForHeaderName(row, headers, hName)
		switch hName {
		case "service_id":
			calendarDate.ServiceId = v
		case "date":
			calendarDate.Date = v
		case "exception_type":
			calendarDate.ExceptionType = v
		}
	}

	return calendarDate
}

func ValidateCalendarDate(cd CalendarDate) []ValidationNotice {
	var validationResults []ValidationNotice

	fields := []struct {
		fieldType FieldType
		name      string
		value     *string
		required  bool
	}{
		{FieldTypeID, "service_id", cd.ServiceId, true},
		{FieldTypeDate, "date", cd.Date, true},
		{FieldTypeCalendarException, "exception_type", cd.ExceptionType, true},
	}

	for _, field := range fields {
		validationResults = append(validationResults, validateField(field.fieldType, field.name, field.value, field.required, FileNameCalendarDate, cd.LineNumber)...)
	}

	return validationResults
}

func ValidateCalendarDates(calendarDates []*CalendarDate, calendarItems []*CalendarItem) []ValidationNotice {
	var results []ValidationNotice

	for _, calendarDate := range calendarDates {
		if calendarDate == nil {
			continue
		}
		results = append(results, ValidateCalendarDate(*calendarDate)...)
	}

	if calendarItems != nil {
		validateCalendarDateReferences(calendarDates, calendarItems, &results)
	}

	return results
}

func validateCalendarDateReferences(calendarDates []*CalendarDate, calendarItems []*CalendarItem, results *[]ValidationNotice) {
	serviceIDMap := make(map[string]struct{})
	for _, item := range calendarItems {
		if item != nil && !StringIsNilOrEmpty(item.ServiceId) {
			serviceIDMap[*item.ServiceId] = struct{}{}
		}
	}

	for _, calendarDate := range calendarDates {
		if calendarDate == nil || StringIsNilOrEmpty(calendarDate.ServiceId) {
			continue
		}
		if _, found := serviceIDMap[*calendarDate.ServiceId]; !found {
			*results = append(*results, ForeignKeyViolationNotice{
				ReferencingFileName:  FileNameCalendarDate,
				ReferencingFieldName: "service_id",
				ReferencedFieldName:  FileNameCalendar,
				ReferencedFileName:   "service_id",
				OffendingValue:       *calendarDate.ServiceId,
				ReferencedAtRow:      calendarDate.LineNumber,
			})
		}
	}
}
