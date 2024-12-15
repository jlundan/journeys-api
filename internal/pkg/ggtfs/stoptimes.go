package ggtfs

import (
	"encoding/csv"
	"fmt"
)

// StopTime struct with fields as strings and optional fields as string pointers.
type StopTime struct {
	TripId            string  // trip_id
	ArrivalTime       string  // arrival_time
	DepartureTime     string  // departure_time
	StopId            string  // stop_id
	StopSequence      string  // stop_sequence (string for consistency)
	StopHeadSign      *string // stop_headsign (optional)
	PickupType        *string // pickup_type (optional)
	DropOffType       *string // drop_off_type (optional)
	ContinuousPickup  *string // continuous_pickup (optional)
	ContinuousDropOff *string // continuous_drop_off (optional)
	ShapeDistTraveled *string // shape_dist_traveled (optional)
	Timepoint         *string // timepoint (optional)
	LineNumber        int     // Line number in the CSV file for error reporting
}

// ValidStopTimeHeaders defines the headers expected in the stop_times CSV file.
var validStopTimeHeaders = []string{"trip_id", "arrival_time", "departure_time", "stop_id", "stop_sequence",
	"stop_headsign", "pickup_type", "drop_off_type", "continuous_pickup", "continuous_drop_off",
	"shape_dist_traveled", "timepoint"}

// LoadStopTimes loads stop times from a CSV reader and returns them along with any errors.
func LoadStopTimes(csvReader *csv.Reader) ([]*StopTime, []error) {
	entities, errs := loadEntities(csvReader, validStopTimeHeaders, CreateStopTime, StopTimesFileName)

	stopTimes := make([]*StopTime, 0, len(entities))
	for _, entity := range entities {
		if stopTime, ok := entity.(*StopTime); ok {
			stopTimes = append(stopTimes, stopTime)
		}
	}

	return stopTimes, errs
}

// CreateStopTime creates and validates a StopTime instance from the CSV row data.
func CreateStopTime(row []string, headers map[string]int, lineNumber int) interface{} {
	var parseErrors []error

	stopTime := StopTime{
		LineNumber: lineNumber,
	}

	for hName, hPos := range headers {
		switch hName {
		case "trip_id":
			stopTime.TripId = getField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "arrival_time":
			stopTime.ArrivalTime = getField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "departure_time":
			stopTime.DepartureTime = getField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "stop_id":
			stopTime.StopId = getField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "stop_sequence":
			stopTime.StopSequence = getField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "stop_headsign":
			stopTime.StopHeadSign = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "pickup_type":
			stopTime.PickupType = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "drop_off_type":
			stopTime.DropOffType = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "continuous_pickup":
			stopTime.ContinuousPickup = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "continuous_drop_off":
			stopTime.ContinuousDropOff = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "shape_dist_traveled":
			stopTime.ShapeDistTraveled = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		case "timepoint":
			stopTime.Timepoint = getOptionalField(row, hName, hPos, &parseErrors, lineNumber, StopTimesFileName)
		}
	}

	if len(parseErrors) > 0 {
		return &stopTime
	}
	return &stopTime
}

// ValidateStopTimes performs additional validation for a list of StopTime instances.
func ValidateStopTimes(stopTimes []*StopTime, stops []*Stop) ([]error, []string) {
	var validationErrors []error
	var recommendations []string

	if stopTimes == nil {
		return validationErrors, recommendations
	}

	for _, stopTime := range stopTimes {
		// Additional required field checks for individual StopTime.
		if stopTime.TripId == "" {
			validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, stopTime.LineNumber, "trip_id must be specified"))
		}
		if stopTime.ArrivalTime == "" {
			validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, stopTime.LineNumber, "arrival_time must be specified"))
		}
		if stopTime.DepartureTime == "" {
			validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, stopTime.LineNumber, "departure_time must be specified"))
		}
		if stopTime.StopId == "" {
			validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, stopTime.LineNumber, "stop_id must be specified"))
		}
		if stopTime.StopSequence == "" {
			validationErrors = append(validationErrors, createFileRowError(StopTimesFileName, stopTime.LineNumber, "stop_sequence must be specified"))
		}
	}

	// Group stop times by trip_id to ensure each trip has at least two stops.
	tripIdToStopCount := make(map[string]int)
	for _, stopTimeItem := range stopTimes {
		if stopTimeItem == nil {
			continue
		}
		tripIdToStopCount[stopTimeItem.TripId] = 0
	}

	for _, stopTimeItem := range stopTimes {
		if stopTimeItem == nil {
			continue
		}
		stopFound := false
		if stops != nil {
			for _, stop := range stops {
				if stop == nil {
					continue
				}
				if stopTimeItem.StopId == stop.Id.String() {
					stopFound = true
					break
				}
			}
		}
		if !stopFound {
			validationErrors = append(validationErrors, createFileError(StopTimesFileName, fmt.Sprintf("trip (%v) references to an unknown stop_id (%s)", stopTimeItem.TripId, stopTimeItem.StopId)))
		} else {
			tripIdToStopCount[stopTimeItem.TripId]++
		}
	}

	for tripId, stopCount := range tripIdToStopCount {
		if stopCount < 2 {
			validationErrors = append(validationErrors, createFileError(StopTimesFileName, fmt.Sprintf("trip (%v) has less than two defined stop times", tripId)))
		}
	}

	return validationErrors, recommendations
}
