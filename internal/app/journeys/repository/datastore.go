package repository

func NewJourneysDataStore(gtfsPath string, skipValidation bool) *JourneysDataStore {
	bundle := newGTFSBundle(gtfsPath, skipValidation)

	linesDataStore := newLineDataStore(bundle.Routes)
	routesDataStore := newRouteDataStore(bundle.Shapes)
	municipalityDataStore := newMunicipalityDataStore(*bundle.Municipalities)
	stopPointDataStore := newStopPointDataStore(bundle.Stops, municipalityDataStore)
	journeyDataStore, journeyPatternDataStore := newJourneyAndJourneyPatternDatastore(bundle.StopTimes, bundle.Trips, bundle.CalendarItems, bundle.CalendarDates, *stopPointDataStore, *linesDataStore, *routesDataStore)

	return &JourneysDataStore{
		Lines:           linesDataStore,
		StopPoints:      stopPointDataStore,
		Municipalities:  municipalityDataStore,
		Routes:          routesDataStore,
		Journeys:        journeyDataStore,
		JourneyPatterns: journeyPatternDataStore,
	}
}

type JourneysDataStore struct {
	Lines           *JourneysLineDataStore
	StopPoints      *JourneysStopPointDataStore
	Municipalities  *JourneysMunicipalityDataStore
	Routes          *JourneysRouteDataStore
	Journeys        *JourneysJourneyDataStore
	JourneyPatterns *JourneysJourneyPatternDataStore
}
