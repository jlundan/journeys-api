package ggtfs

type GtfsEntity interface {
	*Shape | *Stop | *Agency | *CalendarItem | *CalendarDate | *Route | *StopTime | *Trip | any
}

type entityCreator[T GtfsEntity] func(row []string, headers map[string]int, lineNumber int) T

type FieldType string

const (
	FieldTypeColor              FieldType = "Color"
	FieldTypeCurrencyCode       FieldType = "CurrencyCode"
	FieldTypeCurrencyAmount     FieldType = "CurrencyAmount"
	FieldTypeDate               FieldType = "Date"
	FieldTypeEmail              FieldType = "Email"
	FieldTypeID                 FieldType = "ID"
	FieldTypeLanguageCode       FieldType = "LanguageCode"
	FieldTypeLatitude           FieldType = "Latitude"
	FieldTypeLongitude          FieldType = "Longitude"
	FieldTypeFloat              FieldType = "Float"
	FieldTypeInteger            FieldType = "Integer"
	FieldTypePhoneNumber        FieldType = "PhoneNumber"
	FieldTypeTime               FieldType = "Time"
	FieldTypeText               FieldType = "Text"
	FieldTypeTimezone           FieldType = "Timezone"
	FieldTypeURL                FieldType = "URL"
	FieldTypeCalendarDay        FieldType = "CalendarDay"
	FieldTypeCalendarException  FieldType = "CalendarException"
	FieldTypeRouteType          FieldType = "RouteType"
	FieldTypeContinuousPickup   FieldType = "ContinuousPickup"
	FieldTypeContinuousDropOff  FieldType = "ContinuousDropOff"
	FieldTypeLocationType       FieldType = "LocationType"
	FieldTypeWheelchairBoarding FieldType = "WheelchairBoarding"
	FieldTypePickupType         FieldType = "PickupType"
	FieldTypeDropOffType        FieldType = "DropOffType"
	FieldTypeTimepoint          FieldType = "Timepoint"

	FieldTypeDirectionId          FieldType = "DirectionId"
	FieldTypeWheelchairAccessible FieldType = "WheelchairAccessible"
	FieldTypeBikesAllowed         FieldType = "BikesAllowed"

	FileNameAgency       = "agency.txt"
	FileNameCalendar     = "calendar.txt"
	FileNameCalendarDate = "calendar_dates.txt"
	FileNameRoutes       = "routes.txt"
	FileNameShapes       = "shapes.txt"
	FileNameStops        = "stops.txt"
	FileNameStopTimes    = "stop_times.txt"
	FileNameTrips        = "trips.txt"
)
