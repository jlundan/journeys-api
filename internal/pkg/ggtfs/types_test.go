package ggtfs

import "testing"

func TestID_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid ID", "route_123", true},
		{"Empty ID", "", false},
	}

	for _, tt := range tests {
		id := ID{base{raw: tt.value}}
		if got := id.IsValid(); got != tt.valid {
			t.Errorf("ID.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestColor_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Color", "FFFFFF", true},
		{"Invalid Color with #", "#FFFFFF", false},
		{"Invalid Color length", "FFFFF", false},
		{"Invalid Color characters", "ZZZZZZ", false},
	}

	for _, tt := range tests {
		color := Color{base{raw: tt.value}}
		if got := color.IsValid(); got != tt.valid {
			t.Errorf("Color.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestEmail_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Email", "test@example.com", true},
		{"Missing domain", "test@", false},
		{"Invalid format", "testexample.com", false},
		{"Empty Email", "", false},
	}

	for _, tt := range tests {
		email := Email{base{raw: tt.value}}
		if got := email.IsValid(); got != tt.valid {
			t.Errorf("Email.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestURL_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid HTTP URL", "http://example.com", true},
		{"Valid HTTPS URL", "https://example.com", true},
		{"Invalid URL without scheme", "example.com", false},
		{"Invalid URL with unknown scheme", "ftp://example.com", false},
		{"Empty URL", "", false},
	}

	for _, tt := range tests {
		url := URL{base{raw: tt.value}}
		if got := url.IsValid(); got != tt.valid {
			t.Errorf("URL.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestTime_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Time HH:MM:SS", "14:30:00", true},
		{"Valid Time H:MM:SS", "2:30:00", true},
		{"Invalid Time format", "14:30", false},
		{"Invalid Time with out of bounds value", "25:30:00", false},
		{"Empty Time", "", false},
	}

	for _, tt := range tests {
		time := Time{base{raw: tt.value}}
		if got := time.IsValid(); got != tt.valid {
			t.Errorf("Time.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestCurrencyCode_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Currency Code", "USD", true},
		{"Lowercase Currency Code", "usd", false},
		{"Invalid Length", "US", false},
		{"Invalid Characters", "U$D", false},
	}

	for _, tt := range tests {
		cc := CurrencyCode{base{raw: tt.value}}
		if got := cc.IsValid(); got != tt.valid {
			t.Errorf("CurrencyCode.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestCurrencyAmount_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Decimal Amount", "12.50", true},
		{"Valid Integer Amount", "100", true},
		{"Invalid Characters", "12.50$", false},
		{"Invalid Format", "12,50", false},
	}

	for _, tt := range tests {
		ca := CurrencyAmount{base{raw: tt.value}}
		if got := ca.IsValid(); got != tt.valid {
			t.Errorf("CurrencyAmount.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestDate_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Date", "20220101", true},
		{"Invalid Format", "2022/01/01", false},
		{"Invalid Length", "202201", false},
		{"Empty Date", "", false},
	}

	for _, tt := range tests {
		date := Date{base{raw: tt.value}}
		if got := date.IsValid(); got != tt.valid {
			t.Errorf("Date.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestLanguageCode_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Language Code", "en", true},
		{"Valid Extended Language Code", "en-US", true},
		{"Invalid Code with Spaces", "en US", false},
		{"Invalid Characters", "en_US", false},
		{"Empty Language Code", "", false},
	}

	for _, tt := range tests {
		lc := LanguageCode{base{raw: tt.value}}
		if got := lc.IsValid(); got != tt.valid {
			t.Errorf("LanguageCode.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestLatitude_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Latitude", "41.890169", true},
		{"Negative Latitude", "-41.890169", true},
		{"Out of Range Latitude", "95.000000", false},
		{"Invalid Format", "41,890169", false},
		{"Empty Latitude", "", false},
	}

	for _, tt := range tests {
		lat := Latitude{base{raw: tt.value}}
		if got := lat.IsValid(); got != tt.valid {
			t.Errorf("Latitude.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestLongitude_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Longitude", "12.492269", true},
		{"Negative Longitude", "-12.492269", true},
		{"Out of Range Longitude", "195.000000", false},
		{"Invalid Format", "12,492269", false},
		{"Empty Longitude", "", false},
	}

	for _, tt := range tests {
		lon := Longitude{base{raw: tt.value}}
		if got := lon.IsValid(); got != tt.valid {
			t.Errorf("Longitude.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestPhoneNumber_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Phone Number", "+1-800-123-4567", true},
		{"Valid Phone Number", "+358114041122", true},
		{"Valid Short Phone Number", "12345", true},
		{"Invalid Characters", "123-abc-7890", false},
		{"Empty Phone Number", "", false},
	}

	for _, tt := range tests {
		pn := PhoneNumber{base{raw: tt.value}}
		if got := pn.IsValid(); got != tt.valid {
			t.Errorf("PhoneNumber.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestText_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Text", "Some GTFS Text", true},
		{"Empty Text", "", false},
	}

	for _, tt := range tests {
		text := Text{base{raw: tt.value}}
		if got := text.IsValid(); got != tt.valid {
			t.Errorf("Text.IsValid() = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}

func TestTimezone_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{"Valid Timezone", "America/New_York", true},
		{"Valid Timezone with underscore", "America/Argentina_Buenos_Aires", true},
		{"Invalid Timezone with space", "America New_York", false},
		{"Invalid Characters", "123/Invalid", false},
		{"Empty Timezone", "", false},
	}

	for _, tt := range tests {
		tz := Timezone{base{raw: tt.value}}
		if got := tz.IsValid(); got != tt.valid {
			t.Errorf("got = %v, want %v, case: %s", got, tt.valid, tt.name)
		}
	}
}
