package service

import "github.com/jlundan/journeys-api/internal/app/journeys/repository"

func NewJourneysDataService(journeysRepository *repository.JourneysRepository) *JourneysDataService {
	return &JourneysDataService{
		JourneyPatterns: &JourneyPatternsService{Repository: journeysRepository},
		Journeys:        &JourneysService{Repository: journeysRepository},
		Lines:           &LinesService{Repository: journeysRepository},
		Municipalities:  &MunicipalitiesService{Repository: journeysRepository},
		Routes:          &RoutesService{Repository: journeysRepository},
		StopPoints:      &StopPointsService{Repository: journeysRepository},
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
