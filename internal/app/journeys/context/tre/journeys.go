package tre

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
	"github.com/jlundan/journeys-api/internal/pkg/ggtfs"
	"github.com/jlundan/journeys-api/internal/pkg/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Journeys struct {
	all          []*model.Journey
	byId         map[string]*model.Journey
	byActivityId map[string]*model.Journey
}

func (journeys Journeys) GetOne(name string) (*model.Journey, error) {
	if _, ok := journeys.byId[name]; !ok {
		return &model.Journey{}, errors.New("no such element")
	}
	return journeys.byId[name], nil
}
func (journeys Journeys) GetOneByActivityId(name string) (*model.Journey, error) {
	if _, ok := journeys.byActivityId[name]; !ok {
		return &model.Journey{}, errors.New("no such element")
	}
	return journeys.byActivityId[name], nil
}
func (journeys Journeys) GetAll() []*model.Journey {
	return journeys.all
}

type JourneyPatterns struct {
	all  []*model.JourneyPattern
	byId map[string]*model.JourneyPattern
}

func (journeyPatterns JourneyPatterns) GetOne(name string) (*model.JourneyPattern, error) {
	if _, ok := journeyPatterns.byId[name]; !ok {
		return &model.JourneyPattern{}, errors.New("no such element")
	}
	return journeyPatterns.byId[name], nil
}

func (journeyPatterns JourneyPatterns) GetAll() []*model.JourneyPattern {
	return journeyPatterns.all
}

func buildJourneys(g GTFSContext, lines Lines, routes Routes, stopPoints StopPoints) (Journeys, JourneyPatterns) {
	var all = make([]*model.Journey, 0)
	var byId = make(map[string]*model.Journey)
	var byActivityId = make(map[string]*model.Journey)
	var allJourneyPatterns = make([]*model.JourneyPattern, 0)
	var journeyPatternsById = make(map[string]*model.JourneyPattern)

	var tripIdToJourneyPattern = make(map[string]*model.JourneyPattern)
	var tripIdToJourneyCalls = make(map[string][]*model.JourneyCall)
	var tripIdToStopTimes = make(map[string][]*ggtfs.StopTime)

	for _, st := range g.StopTimes {
		tripIdToStopTimes[st.TripId] = append(tripIdToStopTimes[st.TripId], st)
	}

	usedHashes := make([]string, 0)

	for tripId, stArr := range tripIdToStopTimes {
		sort.Slice(stArr, func(x, y int) bool {
			return stArr[x].StopSequence < stArr[y].StopSequence
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
			sp, err := stopPoints.GetOne(stopTime.StopId)
			if err != nil {
				fmt.Println(fmt.Sprintf("Unknown stop point in trip, ignoring it. trip_id:%v, stop_id:%v", tripId, stopTime.TripId))
				continue
			}

			tripIdToJourneyCalls[tripId] = append(tripIdToJourneyCalls[tripId], &model.JourneyCall{
				DepartureTime: stopTime.DepartureTime,
				ArrivalTime:   stopTime.ArrivalTime,
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

	calendarMap := buildCalendarMap(g)
	calendarDateMap := buildCalendarDatesMap(g)

	for _, trip := range g.Trips {
		jp, ok := tripIdToJourneyPattern[trip.Id]
		if !ok {
			fmt.Println(fmt.Sprintf("Journey with no journey pattern detected, ignoring it: %v", trip.Id))
			continue
		}

		line, err := lines.GetOne(trip.RouteId)
		if err != nil {
			fmt.Println(fmt.Sprintf("Journey with no line detected, ignoring it: %v", trip.Id))
			continue
		}

		if trip.ShapeId == nil {
			fmt.Println(fmt.Sprintf("Journey with no route detected, ignoring it: %v", trip.Id))
			continue
		}
		route, err := routes.GetOne(*trip.ShapeId)
		if err != nil {
			fmt.Println(fmt.Sprintf("Journey with no route detected, ignoring it: %v", trip.Id))
			continue
		}

		cMapItem, ok := calendarMap[trip.ServiceId]
		if !ok {
			fmt.Println(fmt.Sprintf("Journey with no service detected, ignoring it: %v", trip.Id))
			continue
		}

		cdMapItem, ok := calendarDateMap[trip.ServiceId]
		if !ok {
			cdMapItem = make([]*model.DayTypeException, 0)
		}

		calls := tripIdToJourneyCalls[trip.Id]

		dtParts := strings.Split(calls[0].DepartureTime, ":")
		dt := strings.Join(dtParts[:2], "")

		firstCall := calls[0]
		lastCall := calls[len(calls)-1]

		activityId := fmt.Sprintf("%v_%v_%v_%v", line.Name, dt, lastCall.StopPoint.ShortName, firstCall.StopPoint.ShortName)

		journey := model.Journey{
			Id:                   trip.Id,
			HeadSign:             *trip.HeadSign,
			Direction:            strconv.Itoa(*trip.DirectionId),
			WheelchairAccessible: *trip.WheelchairAccessible == 1,
			GtfsInfo: &model.JourneyGtfsInfo{
				TripId: trip.Id,
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
		// jp.Name = fmt.Sprintf("%s - %s", jp.StopPoints[0].Name, jp.StopPoints[len(jp.StopPoints)-1].Name)

		route.Journeys = append(route.Journeys, &journey)
		// route.Name = jp.Name
		route.Line = line

		if !routeContainsJourneyPattern(route, jp) {
			journey.Route.JourneyPatterns = append(journey.Route.JourneyPatterns, jp)
		}

		journeyPatternsById[jp.Id].Journeys = append(journeyPatternsById[jp.Id].Journeys, &journey)

		byActivityId[activityId] = &journey
		all = append(all, &journey)
		byId[trip.Id] = &journey
	}

	sort.Slice(all, func(x, y int) bool {
		return all[x].Id < all[y].Id
	})

	sort.Slice(allJourneyPatterns, func(x, y int) bool {
		return allJourneyPatterns[x].Id < allJourneyPatterns[y].Id
	})

	return Journeys{
			all:          all,
			byId:         byId,
			byActivityId: byActivityId,
		},
		JourneyPatterns{
			all:  allJourneyPatterns,
			byId: journeyPatternsById,
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

type calendarFileRow struct {
	serviceId string
	startDate string
	endDate   string
	dayTypes  []string
}

func buildCalendarMap(g GTFSContext) map[string]calendarFileRow {
	result := make(map[string]calendarFileRow)
	for _, calendarItem := range g.CalendarItems {
		days := make([]string, 0)
		if calendarItem.Monday == "1" {
			days = append(days, "monday")
		}
		if calendarItem.Tuesday == "1" {
			days = append(days, "tuesday")
		}
		if calendarItem.Wednesday == "1" {
			days = append(days, "wednesday")
		}
		if calendarItem.Thursday == "1" {
			days = append(days, "thursday")
		}
		if calendarItem.Friday == "1" {
			days = append(days, "friday")
		}
		if calendarItem.Saturday == "1" {
			days = append(days, "saturday")
		}
		if calendarItem.Sunday == "1" {
			days = append(days, "sunday")
		}

		serviceId := calendarItem.ServiceId
		result[serviceId] = calendarFileRow{
			serviceId: serviceId,
			startDate: calendarItem.StartDate,
			endDate:   calendarItem.EndDate,
			dayTypes:  days,
		}
	}

	return result
}

func buildCalendarDatesMap(g GTFSContext) map[string][]*model.DayTypeException {
	result := make(map[string][]*model.DayTypeException)

	for _, calendarDate := range g.CalendarDates {
		var date string

		parsedTime, err := time.Parse("20060102", calendarDate.Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			date = calendarDate.Date
		} else {
			date = parsedTime.Format("2006-01-02")
		}

		result[calendarDate.ServiceId] = append(result[calendarDate.ServiceId], &model.DayTypeException{
			From: date,
			To:   date,
			Runs: calendarDate.ExceptionType == "1",
		})
	}

	return result
}

func stopPointIdsToMd5(arr []*ggtfs.StopTime) string {
	bucket := md5.New()
	for _, v := range arr {
		bucket.Write([]byte(v.StopId))
	}

	return hex.EncodeToString(bucket.Sum(nil))
}
