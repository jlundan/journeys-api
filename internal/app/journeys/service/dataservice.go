package service

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
)

type DataService struct {
	DataStore *repository.JourneysDataStore
}
