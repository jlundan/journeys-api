package repository

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/app/journeys/utils"
	"github.com/jlundan/journeys-api/pkg/ggtfs"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

func newJourneysAndJourneyPatternsRepository(stopTimes []*ggtfs.StopTime, trips []*ggtfs.Trip, calendarItems []*ggtfs.CalendarItem,
	calendarDates []*ggtfs.CalendarDate, stopPointDataStore JourneysStopPointsRepository, lineDataStore JourneysLinesRepository,
	routeDataStore JourneysRoutesRepository) (*JourneysJourneyRepository, *JourneysJourneyPatternRepository) {

	var all = make([]*model.Journey, 0)
	var byId = make(map[string]*model.Journey)
	var byActivityId = make(map[string]*model.Journey)
	var allJourneyPatterns = make([]*model.JourneyPattern, 0)
	var journeyPatternsById = make(map[string]*model.JourneyPattern)

	var tripIdToJourneyPattern = make(map[string]*model.JourneyPattern)
	var tripIdToJourneyCalls = make(map[string][]*model.JourneyCall)
	var tripIdToStopTimes = make(map[string][]*ggtfs.StopTime)

	for i, st := range stopTimes {
		if st == nil {
			log.Println(fmt.Sprintf("Nil stopTime detected, number %v in the stopTimes array, newJourneysAndJourneyPatternsRepository function", i))
			continue
		}

		if st.TripId == nil {
			log.Println(fmt.Sprintf("stoptime.TripId is missing, GTFS line: %v", st.LineNumber))
			continue
		}

		tripId := strings.TrimSpace(*st.TripId)
		tripIdToStopTimes[tripId] = append(tripIdToStopTimes[tripId], st)
	}

	usedHashes := make([]string, 0)

	for tripId, stArr := range tripIdToStopTimes {
		sort.Slice(stArr, func(x, y int) bool {
			if stArr == nil || stArr[x] == nil || stArr[x].StopSequence == nil || stArr[y] == nil || stArr[y].StopSequence == nil {
				return false
			}
			sx, err := strconv.Atoi(*stArr[x].StopSequence)
			if err != nil {
				return false
			}
			sy, err := strconv.Atoi(*stArr[y].StopSequence)
			if err != nil {
				return false
			}
			return sx < sy
		})

		// GTFS stop_times.txt contains trips (Journeys) and sequence of stops (JourneyPatterns) merged into one stop time list
		// (the list is a list of stops in sequence added with arrival and departure times for each of the stops).
		// This means the same sequence of stops (JourneyPattern) is repeated for multiple trips (Journeys). Journeys API
		// separates JourneyPatterns and Journeys to their own entities and links them together. This allows it
		// to maintain single list of sequences of stops and refer to those with calculated hash when needed. The
		// hash is a simple md5 has which is created from a list of stop identifiers.

		// JourneyPatterns are tracked with a hash calculated from a list (sequence) of stops. Each JourneyPattern is
		// essentially just a list of model.StopPoints identified by a md5 hash calculated from those StopPoints. We are
		// iterating through a ggtfs.StopTimes list which means that for JourneyPatterns, we should discard any sequence of stops
		// we have previously encountered in the file, and create JourneyPattern for each one which we have not seen previously.
		// We use the stopList hash for that. Each row in the list, on the other hand, maps to a JourneyCall, so we keep track
		// of the JourneyCalls in relation to their respective trip ids. Each trip will be converted to a Journey later on,
		// so we can simply just assign the collected JourneyCalls to the Journey via the trip id since each trip matches to a Journey.

		stopListHash := stopPointIdsToMd5(stArr)

		var jp *model.JourneyPattern
		if !utils.StringArrayContainsItem(usedHashes, stopListHash) {
			jp = &model.JourneyPattern{
				Id: stopListHash,
			}
		}

		for _, stopTime := range stArr {
			sp, spFound := stopPointDataStore.ById[*stopTime.StopId] // stArr is already filtered to contain only non-nil stopTimes

			if !spFound {
				fmt.Println(fmt.Sprintf("Unknown stop point in trip, ignoring it. trip_id:%v, stop_id:%v", tripId, stopTime.TripId))
				continue
			}

			var arrivalTime, departureTime string

			if stopTime.ArrivalTime != nil {
				arrivalTime = strings.TrimSpace(*stopTime.ArrivalTime)
			} else {
				log.Println(fmt.Sprintf("stoptime (on gtfs row %v): ArrivalTime is missing", stopTime.LineNumber))
			}

			if stopTime.DepartureTime != nil {
				departureTime = strings.TrimSpace(*stopTime.DepartureTime)
			} else {
				log.Println(fmt.Sprintf("stoptime (on gtfs row %v): DepartureTime is missing", stopTime.LineNumber))
			}

			tripIdToJourneyCalls[tripId] = append(tripIdToJourneyCalls[tripId], &model.JourneyCall{ // tripId is already trimmed for spaces
				DepartureTime: departureTime,
				ArrivalTime:   arrivalTime,
				StopPoint:     sp,
			})

			if jp != nil {
				jp.StopPoints = append(jp.StopPoints, sp)
			}
		}

		if jp != nil {
			usedHashes = append(usedHashes, stopListHash)
			allJourneyPatterns = append(allJourneyPatterns, jp)
			journeyPatternsById[stopListHash] = jp
		}

		tripIdToJourneyPattern[tripId] = journeyPatternsById[stopListHash]
	}

	calendarMap := buildCalendarMap(calendarItems)
	calendarDateMap := buildCalendarDatesMap(calendarDates)

	for i, trip := range trips {
		if trip == nil {
			log.Println(fmt.Sprintf("Nil trip detected, number %v in the trips array, newJourneysAndJourneyPatternsRepository function", i))
			continue
		}

		if trip.Id == nil {
			fmt.Println(fmt.Sprintf("trip with no id detected, ignoring it. GTFS trip row: %v", trip.LineNumber))
			continue
		}

		if trip.RouteId == nil {
			fmt.Println(fmt.Sprintf("trip with no RouteId detected, ignoring it. GTFS trip row: %v", trip.LineNumber))
			continue
		}

		if trip.ShapeId == nil {
			fmt.Println(fmt.Sprintf("trip with no ShapeId detected, ignoring it. GTFS trip row: %v", trip.LineNumber))
			continue
		}

		if trip.ServiceId == nil {
			fmt.Println(fmt.Sprintf("trip with no ServiceId detected, ignoring it. GTFS trip row: %v", trip.LineNumber))
			continue
		}

		tripId := strings.TrimSpace(*trip.Id)
		routeId := strings.TrimSpace(*trip.RouteId)
		shapeId := strings.TrimSpace(*trip.ShapeId)
		serviceId := strings.TrimSpace(*trip.ServiceId)

		jp, ok := tripIdToJourneyPattern[tripId]
		if !ok {
			fmt.Println(fmt.Sprintf("Journey with no journey pattern detected, ignoring it: %v", tripId))
			continue
		}

		line, lineFound := lineDataStore.ById[routeId]
		if !lineFound {
			fmt.Println(fmt.Sprintf("Journey with no line detected, ignoring it: %v", tripId))
			continue
		}

		route, routeFound := routeDataStore.ById[shapeId]
		if !routeFound {
			fmt.Println(fmt.Sprintf("Journey with no route detected, ignoring it: %v", tripId))
			continue
		}

		cMapItem, ok := calendarMap[serviceId]
		if !ok {
			fmt.Println(fmt.Sprintf("Journey with no service detected, ignoring it: %v", tripId))
			continue
		}

		cdMapItem, ok := calendarDateMap[serviceId]
		if !ok {
			cdMapItem = make([]*model.DayTypeException, 0)
		}

		calls := tripIdToJourneyCalls[tripId]
		if len(calls) == 0 {
			fmt.Println(fmt.Sprintf("Journey with no calls detected, ignoring it: %v", tripId))
			continue
		}

		dtParts := strings.Split(calls[0].DepartureTime, ":")
		dt := strings.Join(dtParts[:2], "")

		firstCall := calls[0]
		lastCall := calls[len(calls)-1]

		activityId := fmt.Sprintf("%v_%v_%v_%v", line.Name, dt, lastCall.StopPoint.ShortName, firstCall.StopPoint.ShortName)

		var headSign, directionId, wheelChairAccessible string

		if trip.HeadSign != nil {
			headSign = strings.TrimSpace(*trip.HeadSign)
		} else {
			log.Println(fmt.Sprintf("trip (on gtfs row %v): HeadSign is missing", trip.LineNumber))
		}

		if trip.DirectionId != nil {
			directionId = strings.TrimSpace(*trip.DirectionId)
		} else {
			log.Println(fmt.Sprintf("trip (on gtfs row %v): DirectionId is missing", trip.LineNumber))
		}

		if trip.WheelchairAccessible != nil {
			wheelChairAccessible = strings.TrimSpace(*trip.WheelchairAccessible)
		} else {
			log.Println(fmt.Sprintf("trip (on gtfs row %v): WheelchairAccessible is missing", trip.LineNumber))
		}

		journey := model.Journey{
			Id:                   tripId,
			HeadSign:             headSign,
			Direction:            directionId,
			WheelchairAccessible: wheelChairAccessible == "1",
			GtfsInfo: &model.JourneyGtfsInfo{
				TripId: tripId,
			},
			DayTypes:          cMapItem.dayTypes,
			DayTypeExceptions: cdMapItem,
			Calls:             calls,
			Line:              line,
			JourneyPattern:    jp,
			ValidFrom:         cMapItem.startDate,
			ValidTo:           cMapItem.endDate,
			Route:             route,
			ArrivalTime:       lastCall.ArrivalTime,
			DepartureTime:     firstCall.DepartureTime,
			ActivityId:        activityId,
		}

		jp.Route = route

		route.Journeys = append(route.Journeys, &journey)
		route.Line = line

		if !routeContainsJourneyPattern(route, jp) {
			journey.Route.JourneyPatterns = append(journey.Route.JourneyPatterns, jp)
		}

		journeyPatternsById[jp.Id].Journeys = append(journeyPatternsById[jp.Id].Journeys, &journey)

		byActivityId[activityId] = &journey
		all = append(all, &journey)
		byId[*trip.Id] = &journey
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].Id < all[y].Id
	})

	sort.Slice(allJourneyPatterns, func(x, y int) bool {
		return allJourneyPatterns[x].Id < allJourneyPatterns[y].Id
	})

	return &JourneysJourneyRepository{
			All:          all,
			ById:         byId,
			ByActivityId: byActivityId,
		},
		&JourneysJourneyPatternRepository{
			All:  allJourneyPatterns,
			ById: journeyPatternsById,
		}
}

func routeContainsJourneyPattern(route *model.Route, journeyPattern *model.JourneyPattern) bool {
	for _, jp := range route.JourneyPatterns {
		if jp.Id == journeyPattern.Id {
			return true
		}
	}

	return false
}

type JourneysJourneyRepository struct {
	All          []*model.Journey
	ById         map[string]*model.Journey
	ByActivityId map[string]*model.Journey
}

type JourneysJourneyPatternRepository struct {
	All  []*model.JourneyPattern
	ById map[string]*model.JourneyPattern
}

type calendarFileRow struct {
	serviceId string
	startDate string
	endDate   string
	dayTypes  []string
}

func buildCalendarMap(calendarItems []*ggtfs.CalendarItem) map[string]calendarFileRow {
	result := make(map[string]calendarFileRow)
	for i, ci := range calendarItems {
		if ci == nil {
			log.Println(fmt.Sprintf("Nil calendar item detected, number %v in the calendarItems array, buildCalendarMap function", i))
			continue
		}

		if ci.ServiceId == nil || ci.Monday == nil || ci.Tuesday == nil || ci.Wednesday == nil || ci.Thursday == nil ||
			ci.Friday == nil || ci.Saturday == nil || ci.Sunday == nil || ci.StartDate == nil || ci.EndDate == nil {
			log.Println(fmt.Sprintf("malformed calendar item, GTFS row: %v", ci.LineNumber))
			continue
		}

		serviceId := strings.TrimSpace(*ci.ServiceId)
		monday := strings.TrimSpace(*ci.Monday)
		tuesday := strings.TrimSpace(*ci.Tuesday)
		wednesday := strings.TrimSpace(*ci.Wednesday)
		thursday := strings.TrimSpace(*ci.Thursday)
		friday := strings.TrimSpace(*ci.Friday)
		saturday := strings.TrimSpace(*ci.Saturday)
		sunday := strings.TrimSpace(*ci.Sunday)
		startDate := strings.TrimSpace(*ci.StartDate)
		endDate := strings.TrimSpace(*ci.EndDate)

		days := make([]string, 0)
		if monday == "1" {
			days = append(days, "monday")
		}
		if tuesday == "1" {
			days = append(days, "tuesday")
		}
		if wednesday == "1" {
			days = append(days, "wednesday")
		}
		if thursday == "1" {
			days = append(days, "thursday")
		}
		if friday == "1" {
			days = append(days, "friday")
		}
		if saturday == "1" {
			days = append(days, "saturday")
		}
		if sunday == "1" {
			days = append(days, "sunday")
		}

		var formattedStartDate string
		parsedStartDate, err := time.Parse("20060102", startDate)
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing start date for calendar item, GTFS row: %v", ci.LineNumber))
			formattedStartDate = startDate
		} else {
			formattedStartDate = parsedStartDate.Format("2006-01-02")
		}

		var formattedEndDate string
		parsedEndDate, err := time.Parse("20060102", endDate)
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing end date for calendar item, GTFS row: %v", ci.LineNumber))
			formattedEndDate = endDate
		} else {
			formattedEndDate = parsedEndDate.Format("2006-01-02")
		}

		result[serviceId] = calendarFileRow{
			serviceId: serviceId,
			startDate: formattedStartDate,
			endDate:   formattedEndDate,
			dayTypes:  days,
		}
	}

	return result
}

func buildCalendarDatesMap(calendarDates []*ggtfs.CalendarDate) map[string][]*model.DayTypeException {
	result := make(map[string][]*model.DayTypeException)

	for i, cd := range calendarDates {
		if cd == nil {
			log.Println(fmt.Sprintf("Nil calendar date detected, number %v in the calendarDates array, buildCalendarDatesMap function", i))
			continue
		}

		if cd.ServiceId == nil || cd.ExceptionType == nil {
			log.Println(fmt.Sprintf("malformed calendar date, GTFS row: %v", cd.LineNumber))
			continue
		}

		date := strings.TrimSpace(*cd.Date)
		serviceId := strings.TrimSpace(*cd.ServiceId)
		exceptionType := strings.TrimSpace(*cd.ExceptionType)

		var formattedDate string

		parsedTime, err := time.Parse("20060102", date)
		if err != nil {
			log.Println(fmt.Sprintf("Error parsing start date for calendar date, GTFS row: %v", cd.LineNumber))
			formattedDate = date
		} else {
			formattedDate = parsedTime.Format("2006-01-02")
		}

		result[serviceId] = append(result[serviceId], &model.DayTypeException{
			From: formattedDate,
			To:   formattedDate,
			Runs: exceptionType == "1",
		})
	}

	return result
}

func stopPointIdsToMd5(stopTimes []*ggtfs.StopTime) string {
	bucket := md5.New()
	for i, st := range stopTimes {
		if st == nil {
			log.Println(fmt.Sprintf("Nil stopTime detected, number %v in the stopTimes array, stopPointIdsToMd5 function", i))
			continue
		}

		if st.StopId == nil {
			log.Println(fmt.Sprintf("stoptime.StopId is missing, GTFS line: %v", st.LineNumber))
			continue
		}

		bucket.Write([]byte(strings.TrimSpace(*st.StopId)))
	}

	return hex.EncodeToString(bucket.Sum(nil))
}
