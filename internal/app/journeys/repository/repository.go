package repository

func NewJourneysRepository(gtfsPath string, skipValidation bool) *JourneysRepository {
	bundle := newGTFSBundle(gtfsPath, skipValidation)

	linesRepository := newLinesRepository(bundle.Routes)
	routesRepository := newRoutesRepository(bundle.Shapes)
	municipalitiesRepository := newMunicipalitiesRepository(*bundle.Municipalities)
	stopPointsRepository := newStopPointsRepository(bundle.Stops, municipalitiesRepository)
	journeyRepository, journeyPatternRepository := newJourneysAndJourneyPatternsRepository(bundle.StopTimes, bundle.Trips, bundle.CalendarItems, bundle.CalendarDates, *stopPointsRepository, *linesRepository, *routesRepository)

	return &JourneysRepository{
		Lines:           linesRepository,
		StopPoints:      stopPointsRepository,
		Municipalities:  municipalitiesRepository,
		Routes:          routesRepository,
		Journeys:        journeyRepository,
		JourneyPatterns: journeyPatternRepository,
	}
}

type JourneysRepository struct {
	Lines           *JourneysLinesRepository
	StopPoints      *JourneysStopPointsRepository
	Municipalities  *JourneysMunicipalitiesRepository
	Routes          *JourneysRoutesRepository
	Journeys        *JourneysJourneyRepository
	JourneyPatterns *JourneysJourneyPatternRepository
}
