package service

import "github.com/jlundan/journeys-api/internal/app/journeys/repository"

func NewJourneysDataService(dataStore *repository.JourneysDataStore) *JourneysDataService {
	return &JourneysDataService{
		JourneyPatterns: &JourneyPatternsService{DataStore: dataStore},
		Journeys:        &JourneysService{DataStore: dataStore},
		Lines:           &LinesService{DataStore: dataStore},
		Municipalities:  &MunicipalitiesService{DataStore: dataStore},
		Routes:          &RoutesService{DataStore: dataStore},
		StopPoints:      &StopPointsService{DataStore: dataStore},
	}

}

type JourneysDataService struct {
	JourneyPatterns *JourneyPatternsService
	Journeys        *JourneysService
	Lines           *LinesService
	Municipalities  *MunicipalitiesService
	Routes          *RoutesService
	StopPoints      *StopPointsService
}
