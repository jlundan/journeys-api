package service

import (
	"github.com/jlundan/journeys-api/internal/app/repository"
)

type DataService struct {
	DataStore *repository.JourneysDataStore
}
