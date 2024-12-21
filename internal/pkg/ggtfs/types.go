package ggtfs

import (
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// base is a struct that contains the raw value and a flag to indicate if the header for the field is present in the
// parsed CSV file.
// The loader method reads only the fields which have a header present in the file. This means that the loader
// skips all fields which do not have a header. To make the processing more straightforward, the isPresent is
// kept at its default state (false), which means that when the loader does not process a missing field, the GTFS
// entity gets the field struct with default values, if it is mandatory field, which marks the field absent (isPresent is false).
// (The optional fields will get a nil pointer, if the header for that field is missing in the CSV file, which means
// that these structs are not created at all for them).

type base struct {
	raw       string
	isPresent bool
}

func (base *base) String() string {
	return base.raw
}

func (base *base) IsEmpty() bool {
	if base == nil {
		return true
	}

	return strings.TrimSpace(base.raw) == ""
}

func (base *base) Length() int {
	if base == nil {
		return 0
	}

	return len(base.raw)
}

func (base *base) IsPresent() bool {
	if base == nil {
		return false
	}

	return base.isPresent
}

// ID represents an internal ID, such as `route_id` or `trip_id`.
type ID struct {
	base
}

// IsValid checks if the ID is not empty.
func (id *ID) IsValid() bool {
	if id == nil {
		return false
	}

	return !id.IsEmpty() && id.IsPresent()
}

func NewID(raw *string) ID {
	if raw == nil {
		return ID{base{raw: ""}}
	}
	return ID{base{raw: *raw, isPresent: true}}
}

func NewOptionalID(raw *string) *ID {
	if raw == nil {
		return nil
	}
	return &ID{base{raw: *raw, isPresent: true}}
}

// Color represents a color encoded as a six-digit hexadecimal number.
type Color struct {
	base
}

// IsValid checks if the Color is a valid six-digit hexadecimal value.
func (c *Color) IsValid() bool {
	if c == nil {
		return false
	}

	match, _ := regexp.MatchString(`^[0-9A-Fa-f]{6}$`, c.raw)
	return match
}

func NewOptionalColor(raw *string) *Color {
	if raw == nil {
		return nil
	}
	return &Color{base{raw: *raw, isPresent: true}}
}

// Email represents an email address.
type Email struct {
	base
}

// IsValid checks if the Email is in a valid email format.
func (e *Email) IsValid() bool {
	if e == nil {
		return false
	}

	_, err := mail.ParseAddress(e.raw)
	return err == nil
}

func NewOptionalEmail(raw *string) *Email {
	if raw == nil {
		return nil
	}
	return &Email{base{raw: *raw, isPresent: true}}
}

// Integer represents a number without floating point.
type Integer struct {
	base
}

// IsValid if the value is a valid Integer.
func (i *Integer) IsValid() bool {
	if i == nil {
		return false
	}

	_, err := strconv.Atoi(i.raw)
	return err == nil
}

// Int returns an integer value of the raw string received from the CSV file.
// Use IsValid to verify that the conversion can be made. This method will return
// zero for strings that cannot be parsed to int.
func (i *Integer) Int() int {
	val, _ := strconv.Atoi(i.raw)
	return val
}

func NewInteger(raw *string) Integer {
	if raw == nil {
		return Integer{base{raw: ""}}
	}
	return Integer{base{raw: *raw, isPresent: true}}
}

func NewOptionalInteger(raw *string) *Integer {
	if raw == nil {
		return nil
	}
	return &Integer{base{raw: *raw, isPresent: true}}
}

// Float represents a number with a floating point.
type Float struct {
	base
}

func NewOptionalFloat(raw *string) *Float {
	if raw == nil {
		return nil
	}
	return &Float{base{raw: *raw, isPresent: true}}
}

// IsValid if the value is a valid Float.
func (f *Float) IsValid() bool {
	if f == nil {
		return false
	}

	_, err := strconv.ParseFloat(f.raw, 64)
	return err == nil
}

// Float64 returns a 64-bit float value of the raw string received from the CSV file.
// Use IsValid to verify that the conversion can be made. This method will return
// zero for strings that cannot be parsed to float.
func (f *Float) Float64() float64 {
	val, _ := strconv.ParseFloat(f.raw, 64)
	return val
}

// URL represents a fully qualified URL.
type URL struct {
	base
}

// IsValid checks if the URL is well-formed.
func (u *URL) IsValid() bool {
	if u == nil {
		return false
	}

	parsedURL, err := url.ParseRequestURI(u.raw)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func NewURL(raw *string) URL {
	if raw == nil {
		return URL{base{raw: ""}}
	}
	return URL{base{raw: *raw, isPresent: true}}
}

func NewOptionalURL(raw *string) *URL {
	if raw == nil {
		return nil
	}
	return &URL{base{raw: *raw, isPresent: true}}
}

// Time represents a time in HH:MM:SS or H:MM:SS format.
type Time struct {
	base
}

// IsValid checks if the Time is in a valid format.
func (t *Time) IsValid() bool {
	if t == nil {
		return false
	}

	match, _ := regexp.MatchString(`^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$|^([0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`, t.raw)
	return match
}

func NewTime(raw *string) Time {
	if raw == nil {
		return Time{base{raw: ""}}
	}
	return Time{base{raw: *raw, isPresent: true}}
}

func NewOptionalTime(raw *string) *Time {
	if raw == nil {
		return nil
	}
	return &Time{base{raw: *raw, isPresent: true}}
}

// CurrencyCode represents a currency code according to ISO 4217.
type CurrencyCode struct {
	base
}

// IsValid checks if the CurrencyCode is a three-letter code.
func (cc *CurrencyCode) IsValid() bool {
	if cc == nil {
		return false
	}

	match, _ := regexp.MatchString(`^[A-Z]{3}$`, cc.raw)
	return match
}

// CurrencyAmount represents a monetary amount.
type CurrencyAmount struct {
	base
}

// IsValid checks if the CurrencyAmount is a valid decimal number.
func (ca *CurrencyAmount) IsValid() bool {
	if ca == nil {
		return false
	}

	_, err := strconv.ParseFloat(ca.raw, 64)
	return err == nil
}

// Date represents a date in the format YYYYMMDD.
type Date struct {
	base
}

// IsValid checks if the Date is in the valid YYYYMMDD format.
func (d *Date) IsValid() bool {
	if d == nil {
		return false
	}

	match, _ := regexp.MatchString(`^(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])$`, d.raw)
	return match
}

func NewDate(raw *string) Date {
	if raw == nil {
		return Date{base{raw: ""}}
	}
	return Date{base{raw: *raw, isPresent: true}}
}

// LanguageCode represents an IETF BCP 47 language code.
type LanguageCode struct {
	base
}

// IsValid checks if the LanguageCode is a valid IETF BCP 47 code.
func (lc *LanguageCode) IsValid() bool {
	if lc == nil {
		return false
	}
	// Basic validation for language codes: e.g., "en", "en-US"
	match, _ := regexp.MatchString(`^[a-zA-Z]{2,3}(-[a-zA-Z]{2,3})?$`, lc.raw)
	return match
}

func NewOptionalLanguageCode(raw *string) *LanguageCode {
	if raw == nil {
		return nil
	}
	return &LanguageCode{base{raw: *raw, isPresent: true}}
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

func NewOptionalLatitude(raw *string) *Latitude {
	if raw == nil {
		return nil
	}
	return &Latitude{base{raw: *raw, isPresent: true}}
}

// IsValid checks if the Latitude is a valid decimal value between -90 and 90.
func (lat *Latitude) IsValid() bool {
	if lat == nil {
		return false
	}
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

func NewOptionalLongitude(raw *string) *Longitude {
	if raw == nil {
		return nil
	}
	return &Longitude{base{raw: *raw, isPresent: true}}
}

// IsValid checks if the Longitude is a valid decimal value between -180 and 180.
func (lon *Longitude) IsValid() bool {
	if lon == nil {
		return false
	}

	value, err := strconv.ParseFloat(lon.raw, 64)
	return err == nil && value >= -180.0 && value <= 180.0
}

// PhoneNumber represents a phone number.
type PhoneNumber struct {
	base
}

// IsValid checks if the PhoneNumber has a reasonable length and contains only digits and certain symbols.
func (pn *PhoneNumber) IsValid() bool {
	if pn == nil {
		return false
	}
	// Check for minimum length, only contains digits, and common phone number symbols
	match, _ := regexp.MatchString(`^[\d\s\-+()]{5,}$`, pn.raw)
	return match
}

func NewOptionalPhoneNumber(raw *string) *PhoneNumber {
	if raw == nil {
		return nil
	}
	return &PhoneNumber{base{raw: *raw, isPresent: true}}
}

// Text represents a string of UTF-8 characters intended for display.
type Text struct {
	base
}

// IsValid checks if the Text is non-empty.
func (t *Text) IsValid() bool {
	if t == nil {
		return false
	}
	return !t.IsEmpty()
}

func NewText(raw *string) Text {
	if raw == nil {
		return Text{base{raw: ""}}
	}
	return Text{base{raw: *raw, isPresent: true}}
}

func NewOptionalText(raw *string) *Text {
	if raw == nil {
		return nil
	}
	return &Text{base{raw: *raw, isPresent: true}}
}

// Timezone represents a TZ timezone from the IANA timezone database.
type Timezone struct {
	base
}

func NewOptionalTimezone(raw *string) *Timezone {
	if raw == nil {
		return nil
	}
	return &Timezone{base{raw: *raw, isPresent: true}}
}

// IsValid checks if the Timezone is in a valid format (e.g., "America/New_York").
func (tz *Timezone) IsValid() bool {
	if tz == nil {
		return false
	}
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
