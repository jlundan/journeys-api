package ggtfs

import (
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type base struct {
	raw       string
	isPresent bool
}

func (base base) Raw() string {
	return base.raw
}

func (base base) IsEmpty() bool {
	return strings.TrimSpace(base.raw) == ""
}

func (base base) Length() int {
	return len(base.raw)
}

func (base base) IsPresent() bool {
	return base.isPresent
}

type ID struct {
	base
}

func (id ID) IsValid() bool {
	return !id.IsEmpty() && id.IsPresent()
}

func NewID(raw *string) ID {
	if raw == nil {
		return ID{base{raw: ""}}
	}
	return ID{base{raw: *raw, isPresent: true}}
}

type Color struct {
	base
}

func (c Color) IsValid() bool {
	match, _ := regexp.MatchString(`^[0-9A-Fa-f]{6}$`, c.raw)
	return match
}

func NewColor(raw *string) Color {
	if raw == nil {
		return Color{base{raw: ""}}
	}
	return Color{base{raw: *raw, isPresent: true}}
}

type Email struct {
	base
}

func (e Email) IsValid() bool {
	_, err := mail.ParseAddress(e.raw)
	return err == nil
}

func NewEmail(raw *string) Email {
	if raw == nil {
		return Email{base{raw: ""}}
	}
	return Email{base{raw: *raw, isPresent: true}}
}

type Integer struct {
	base
}

func (i Integer) IsValid() bool {
	_, err := strconv.Atoi(i.raw)
	return err == nil
}

func (i Integer) Int() int {
	val, _ := strconv.Atoi(i.raw)
	return val
}

func NewInteger(raw *string) Integer {
	if raw == nil {
		return Integer{base{raw: ""}}
	}
	return Integer{base{raw: *raw, isPresent: true}}
}

type Float struct {
	base
}

func NewFloat(raw *string) Float {
	if raw == nil {
		return Float{base{raw: ""}}
	}
	return Float{base{raw: *raw, isPresent: true}}
}

func (f Float) IsValid() bool {
	_, err := strconv.ParseFloat(f.raw, 64)
	return err == nil
}

func (f Float) Float64() float64 {
	val, _ := strconv.ParseFloat(f.raw, 64)
	return val
}

type URL struct {
	base
}

func (u URL) IsValid() bool {
	parsedURL, err := url.ParseRequestURI(u.raw)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func NewURL(raw *string) URL {
	if raw == nil {
		return URL{base{raw: ""}}
	}
	return URL{base{raw: *raw, isPresent: true}}
}

// Time represents a time in HH:MM:SS or H:MM:SS format.
type Time struct {
	base
}

// IsValid checks if the Time is in the valid HH:MM:SS or H:MM:SS format. The hour is between 0 and 47, since the trips on the service day might run
// through the night. For example, 25:00:00 represents 1:00:00 AM the next day.
func (t Time) IsValid() bool {
	match, _ := regexp.MatchString(`^(0[0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-7]):([0-5][0-9]):([0-5][0-9])$|^([0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`, t.raw)
	return match
}

func NewTime(raw *string) Time {
	if raw == nil {
		return Time{base{raw: ""}}
	}
	return Time{base{raw: *raw, isPresent: true}}
}

// CurrencyCode represents a currency code according to ISO 4217.
type CurrencyCode struct {
	base
}

// IsValid checks if the CurrencyCode is a three-letter code.
func (cc CurrencyCode) IsValid() bool {
	match, _ := regexp.MatchString(`^[A-Z]{3}$`, cc.raw)
	return match
}

type CurrencyAmount struct {
	base
}

func (ca CurrencyAmount) IsValid() bool {
	_, err := strconv.ParseFloat(ca.raw, 64)
	return err == nil
}

type Date struct {
	base
}

// IsValid checks if the Date is in the valid YYYYMMDD format.
func (d Date) IsValid() bool {
	match, _ := regexp.MatchString(`^(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])$`, d.raw)
	return match
}

func NewDate(raw *string) Date {
	if raw == nil {
		return Date{base{raw: ""}}
	}
	return Date{base{raw: *raw, isPresent: true}}
}

type LanguageCode struct {
	base
}

// IsValid checks if the LanguageCode is a valid IETF BCP 47 code.
func (lc LanguageCode) IsValid() bool {
	// Basic validation for language codes: e.g., "en", "en-US"
	match, _ := regexp.MatchString(`^[a-zA-Z]{2,3}(-[a-zA-Z]{2,3})?$`, lc.raw)
	return match
}

func NewLanguageCode(raw *string) LanguageCode {
	if raw == nil {
		return LanguageCode{base{raw: ""}}
	}
	return LanguageCode{base{raw: *raw, isPresent: true}}
}

// Latitude represents a WGS84 latitude in decimal degrees.
type Latitude struct {
	base
}

func NewLatitude(raw *string) Latitude {
	if raw == nil {
		return Latitude{base{raw: ""}}
	}
	return Latitude{base{raw: *raw, isPresent: true}}
}

// IsValid checks if the Latitude is a valid decimal value between -90 and 90.
func (lat Latitude) IsValid() bool {
	value, err := strconv.ParseFloat(lat.raw, 64)
	return err == nil && value >= -90.0 && value <= 90.0
}

// Longitude represents a WGS84 longitude in decimal degrees.
type Longitude struct {
	base
}

func NewLongitude(raw *string) Longitude {
	if raw == nil {
		return Longitude{base{raw: ""}}
	}
	return Longitude{base{raw: *raw, isPresent: true}}
}

// IsValid checks if the Longitude is a valid decimal value between -180 and 180.
func (lon Longitude) IsValid() bool {
	value, err := strconv.ParseFloat(lon.raw, 64)
	return err == nil && value >= -180.0 && value <= 180.0
}

type PhoneNumber struct {
	base
}

func (pn PhoneNumber) IsValid() bool {
	// Check for minimum length, only contains digits, and common phone number symbols
	match, _ := regexp.MatchString(`^[\d\s\-+()]{5,}$`, pn.raw)
	return match
}

func NewPhoneNumber(raw *string) PhoneNumber {
	if raw == nil {
		return PhoneNumber{base{raw: ""}}
	}
	return PhoneNumber{base{raw: *raw, isPresent: true}}
}

type Text struct {
	base
}

func (t Text) IsValid() bool {
	return !t.IsEmpty()
}

func NewText(raw *string) Text {
	if raw == nil {
		return Text{base{raw: ""}}
	}
	return Text{base{raw: *raw, isPresent: true}}
}

// Timezone represents a TZ timezone from the IANA timezone database.
type Timezone struct {
	base
}

// IsValid checks if the Timezone is in a valid format (e.g., "America/New_York").
func (tz Timezone) IsValid() bool {
	// Basic regex to validate Continent/City or Continent/City_Name format.
	// It checks if we have at least a structure like Continent/City.
	match, _ := regexp.MatchString(`^[A-Za-z]+/[A-Za-z_]+$|^[A-Za-z]+/[A-Za-z]+$`, tz.raw)
	if match {
		return true
	}
	return false
}

func NewTimezone(raw *string) Timezone {
	if raw == nil {
		return Timezone{base{raw: ""}}
	}
	return Timezone{base{raw: *raw, isPresent: true}}
}

type PositiveInteger struct {
	Integer
}

func (pi PositiveInteger) IsValid() bool {
	return pi.Integer.IsValid() && pi.Integer.Int() > 0
}

func NewPositiveInteger(raw *string) PositiveInteger {
	if raw == nil {
		return PositiveInteger{
			Integer{base: base{raw: ""}}}
	}
	return PositiveInteger{Integer{base: base{raw: *raw, isPresent: true}}}
}

type PositiveFloat struct {
	Float
}

func (pf PositiveFloat) IsValid() bool {
	return pf.Float.IsValid() && pf.Float.Float64() > 0
}

func NewPositiveFloat(raw *string) PositiveFloat {
	if raw == nil {
		return PositiveFloat{
			Float{base: base{raw: ""}}}
	}
	return PositiveFloat{Float{base: base{raw: *raw, isPresent: true}}}
}

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

	FileNameAgency       = "agency.txt"
	FileNameCalendar     = "calendar.txt"
	FileNameCalendarDate = "calendar_dates.txt"
	FileNameRoutes       = "routes.txt"
	FileNameShapes       = "shapes.txt"
	FileNameStops        = "stops.txt"
	FileNameStopTimes    = "stop_times.txt"
)
