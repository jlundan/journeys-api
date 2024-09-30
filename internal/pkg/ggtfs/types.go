package ggtfs

import (
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
)

type base struct {
	raw string
}

func (base *base) String() string {
	return base.raw
}

func (base *base) IsEmpty() bool {
	return base.raw == ""
}

func (base *base) Length() int {
	return len(base.raw)
}

// ID represents an internal ID, such as `route_id` or `trip_id`.
type ID struct {
	base
}

// IsValid checks if the ID is not empty.
func (id *ID) IsValid() bool {
	return id.raw != ""
}

// Color represents a color encoded as a six-digit hexadecimal number.
type Color struct {
	base
}

// IsValid checks if the Color is a valid six-digit hexadecimal value.
func (c *Color) IsValid() bool {
	match, _ := regexp.MatchString(`^[0-9A-Fa-f]{6}$`, c.raw)
	return match
}

// Email represents an email address.
type Email struct {
	base
}

// IsValid checks if the Email is in a valid email format.
func (e *Email) IsValid() bool {
	_, err := mail.ParseAddress(e.raw)
	return err == nil
}

// URL represents a fully qualified URL.
type URL struct {
	base
}

// IsValid checks if the URL is well-formed.
func (u *URL) IsValid() bool {
	parsedURL, err := url.ParseRequestURI(u.raw)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

// Time represents a time in HH:MM:SS or H:MM:SS format.
type Time struct {
	base
}

// IsValid checks if the Time is in a valid format.
func (t *Time) IsValid() bool {
	match, _ := regexp.MatchString(`^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$|^([0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$`, t.raw)
	return match
}

// CurrencyCode represents a currency code according to ISO 4217.
type CurrencyCode struct {
	base
}

// IsValid checks if the CurrencyCode is a three-letter code.
func (cc *CurrencyCode) IsValid() bool {
	match, _ := regexp.MatchString(`^[A-Z]{3}$`, cc.raw)
	return match
}

// CurrencyAmount represents a monetary amount.
type CurrencyAmount struct {
	base
}

// IsValid checks if the CurrencyAmount is a valid decimal number.
func (ca *CurrencyAmount) IsValid() bool {
	_, err := strconv.ParseFloat(ca.raw, 64)
	return err == nil
}

// Date represents a date in the format YYYYMMDD.
type Date struct {
	base
}

// IsValid checks if the Date is in the valid YYYYMMDD format.
func (d *Date) IsValid() bool {
	match, _ := regexp.MatchString(`^\d{8}$`, d.raw)
	return match
}

// LanguageCode represents an IETF BCP 47 language code.
type LanguageCode struct {
	base
}

// IsValid checks if the LanguageCode is a valid IETF BCP 47 code.
func (lc *LanguageCode) IsValid() bool {
	// Basic validation for language codes: e.g., "en", "en-US"
	match, _ := regexp.MatchString(`^[a-zA-Z]{2,3}(-[a-zA-Z]{2,3})?$`, lc.raw)
	return match
}

// Latitude represents a WGS84 latitude in decimal degrees.
type Latitude struct {
	base
}

// IsValid checks if the Latitude is a valid decimal value between -90 and 90.
func (lat *Latitude) IsValid() bool {
	value, err := strconv.ParseFloat(lat.raw, 64)
	return err == nil && value >= -90.0 && value <= 90.0
}

// Longitude represents a WGS84 longitude in decimal degrees.
type Longitude struct {
	base
}

// IsValid checks if the Longitude is a valid decimal value between -180 and 180.
func (lon *Longitude) IsValid() bool {
	value, err := strconv.ParseFloat(lon.raw, 64)
	return err == nil && value >= -180.0 && value <= 180.0
}

// PhoneNumber represents a phone number.
type PhoneNumber struct {
	base
}

// IsValid checks if the PhoneNumber has a reasonable length and contains only digits and certain symbols.
func (pn *PhoneNumber) IsValid() bool {
	// Check for minimum length, only contains digits, and common phone number symbols
	match, _ := regexp.MatchString(`^[\d\s\-+()]{5,}$`, pn.raw)
	return match
}

// Text represents a string of UTF-8 characters intended for display.
type Text struct {
	base
}

// IsValid checks if the Text is non-empty.
func (t *Text) IsValid() bool {
	return t.raw != ""
}

// Timezone represents a TZ timezone from the IANA timezone database.
type Timezone struct {
	base
}

// IsValid checks if the Timezone is in a valid format (e.g., "America/New_York").
func (tz *Timezone) IsValid() bool {
	// Basic regex to validate Continent/City or Continent/City_Name format.
	// It checks if we have at least a structure like Continent/City.
	match, _ := regexp.MatchString(`^[A-Za-z]+/[A-Za-z_]+$|^[A-Za-z]+/[A-Za-z]+$`, tz.raw)
	if match {
		return true
	}
	return false
}
