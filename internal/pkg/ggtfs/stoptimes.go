package ggtfs

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
)

type StopTime struct {
	TripId            string
	ArrivalTime       string
	DepartureTime     string
	StopId            string
	StopSequence      int
	StopHeadSign      *string
	PickupType        *int
	DropOffType       *int
	ContinuousPickup  *int
	ContinuousDropOff *int
	ShapeDistTraveled *float64
	Timepoint         *int
	lineNumber        int
}

func ExtractStopTimesByTrips(input *csv.Reader, output *csv.Writer, trips map[string]struct{}) (map[string]struct{}, []error) {
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(input)
	if err != nil {
		errs = append(errs, createFileError(StopTimesFileName, fmt.Sprintf("read error: %v", err.Error())))
		return nil, errs
	}

	if headers == nil { // EOF
		return nil, nil
	}

	var tripIdPos uint8
	if pos, hasTripId := headers["trip_id"]; hasTripId {
		tripIdPos = pos
	} else {
		errs = append(errs, createFileError(StopTimesFileName, "cannot stop times trips without trip_id column"))
		return nil, errs
	}

	var stopIdPos *uint8
	if pos, hasId := headers["stop_id"]; hasId {
		stopIdPos = &pos
	} else {
		errs = append(errs, createFileError(StopTimesFileName, "no stop_id column, stops will not be extracted"))
		stopIdPos = nil
	}

	err = writeHeaderRow(headers, output)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	stopIdMap := make(map[string]struct{}, 0)
	for {
		row, rErr := ReadDataRow(input)
		if rErr != nil {
			errs = append(errs, createFileError(StopTimesFileName, fmt.Sprintf("%v", rErr.Error())))
			continue
		}

		if row == nil { // EOF
			break
		}

		if _, shouldBeExtracted := trips[row[tripIdPos]]; shouldBeExtracted {
			wErr := writeDataRow(row, output)
			if wErr != nil {
				errs = append(errs, wErr)
				return nil, errs
			}

			if stopIdPos == nil {
				continue
			}
			stopId := row[*stopIdPos]
			if _, alreadyExists := stopIdMap[stopId]; !alreadyExists {
				stopIdMap[stopId] = struct{}{}
			}
		}
	}

	return stopIdMap, nil
}

func LoadStopTimes(csvReader *csv.Reader) ([]*StopTime, []error) {
	stopPoints := make([]*StopTime, 0)
	errs := make([]error, 0)

	headers, err := ReadHeaderRow(csvReader)
	if err != nil {
		errs = append(errs, createFileError(StopTimesFileName, fmt.Sprintf("read error: %v", err.Error())))
		return stopPoints, errs
	}
	if headers == nil {
		return stopPoints, errs
	}

	index := 0
	for {
		row, err := ReadDataRow(csvReader)
		if err != nil {
			errs = append(errs, createFileError(StopTimesFileName, fmt.Sprintf("%v", err.Error())))
			index++
			continue
		}
		if row == nil {
			break
		}

		rowErrs := make([]error, 0)
		stopPoint := StopTime{
			lineNumber: index,
		}

		var (
			tripId        *string
			stopId        *string
			arrivalTime   *string
			departureTime *string
			sequence      *int
		)
		for name, column := range headers {
			switch name {
			case "trip_id":
				tripId = handleIDField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "arrival_time":
				arrivalTime = handleTimeField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "departure_time":
				departureTime = handleTimeField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "stop_id":
				stopId = handleIDField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "stop_sequence":
				sequence = handleIntField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "stop_headsign":
				stopPoint.StopHeadSign = &row[column]
			case "pickup_type":
				stopPoint.PickupType = handlePickupField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "drop_off_type":
				stopPoint.DropOffType = handleDropOffField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "continuous_pickup":
				stopPoint.ContinuousPickup = handleContinuousPickupField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "continuous_drop_off":
				stopPoint.ContinuousDropOff = handleContinuousDropOffField(row[column], StopTimesFileName, name, index, &rowErrs)
			case "shape_dist_traveled":
				stopPoint.ShapeDistTraveled = handleFloat64Field(row[column], StopTimesFileName, name, index, &rowErrs)
			case "timepoint":
				stopPoint.Timepoint = handleTimePointField(row[column], StopTimesFileName, name, index, &rowErrs)
			}
		}

		if tripId == nil {
			rowErrs = append(rowErrs, createFileRowError(StopTimesFileName, stopPoint.lineNumber, "trip_id must be specified"))
		} else {
			stopPoint.TripId = *tripId
		}

		if stopId == nil {
			rowErrs = append(rowErrs, createFileRowError(StopTimesFileName, stopPoint.lineNumber, "stop_id must be specified"))
		} else {
			stopPoint.StopId = *stopId
		}

		if sequence == nil {
			rowErrs = append(rowErrs, createFileRowError(StopTimesFileName, stopPoint.lineNumber, "stop_sequence must be specified"))
		} else {
			stopPoint.StopSequence = *sequence
		}

		if arrivalTime == nil {
			rowErrs = append(rowErrs, createFileRowError(StopTimesFileName, stopPoint.lineNumber, "arrival_time must be specified"))
		} else {
			stopPoint.ArrivalTime = *arrivalTime
		}

		if departureTime == nil {
			rowErrs = append(rowErrs, createFileRowError(StopTimesFileName, stopPoint.lineNumber, "departure_time must be specified"))
		} else {
			stopPoint.DepartureTime = *departureTime
		}

		if len(rowErrs) > 0 {
			errs = append(errs, rowErrs...)
		} else {
			stopPoints = append(stopPoints, &stopPoint)
		}

		index++
	}

	return stopPoints, errs
}

func ValidateStoptimes(stopTimes []*StopTime, stops []*Stop) []error {
	var validationErrors []error

	if stopTimes == nil {
		return validationErrors
	}

	// Stop times are grouped by trip id, and we want to make sure that each trip id has at least two stops
	// We do this by first finding out all trips in the stopTimes and setting their stop count to zero in a map. Then we loop
	// through all stop times while incrementing the stop counts in our map. Eventually we iterate through the map and check
	// if any trip has stop count less than a zero.
	tripIdToStopCount := make(map[string]uint64)
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
				if stopTimeItem.StopId == stop.Id {
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

	return validationErrors
}

func handlePickupField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v <= 4 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}

func handleDropOffField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v <= 4 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}

func handleTimePointField(str string, fileName string, fieldName string, index int, errs *[]error) *int {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, err))
		return nil
	}

	if v := int(n); v <= 4 {
		return &v
	} else {
		*errs = append(*errs, createFieldError(fileName, fieldName, index, errors.New(invalidValue)))
		return nil
	}
}
