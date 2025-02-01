package repository

import (
	"errors"
)

func NewJourneysRepository(gtfsPath string, skipValidation bool) (*JourneysRepository, []error) {
	bundle := newGTFSBundle(gtfsPath, skipValidation)

	linesRepository := newLinesRepository(bundle.Routes)
	routesRepository := newRoutesRepository(bundle.Shapes)
	municipalitiesRepository := newMunicipalitiesRepository(*bundle.Municipalities)
	stopPointsRepository := newStopPointsRepository(bundle.Stops, municipalitiesRepository)
	journeyRepository, journeyPatternRepository := newJourneysAndJourneyPatternsRepository(bundle.StopTimes, bundle.Trips, bundle.CalendarItems, bundle.CalendarDates, *stopPointsRepository, *linesRepository, *routesRepository)

	errs := getBundleErrorsNotices(bundle)

	return &JourneysRepository{
		Lines:           linesRepository,
		StopPoints:      stopPointsRepository,
		Municipalities:  municipalitiesRepository,
		Routes:          routesRepository,
		Journeys:        journeyRepository,
		JourneyPatterns: journeyPatternRepository,
	}, errs
}

type JourneysRepository struct {
	Lines           *JourneysLinesRepository
	StopPoints      *JourneysStopPointsRepository
	Municipalities  *JourneysMunicipalitiesRepository
	Routes          *JourneysRoutesRepository
	Journeys        *JourneysJourneyRepository
	JourneyPatterns *JourneysJourneyPatternRepository
}

func getBundleErrorsNotices(bundle *GTFSBundle) []error {
	var errs []error

	errs = append(errs, bundle.Errors...)
	for _, notice := range bundle.ValidationNotices {
		errs = append(errs, errors.New(notice.AsText()))
	}

	return errs
}
