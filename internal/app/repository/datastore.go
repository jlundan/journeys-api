package repository

import (
	"github.com/jlundan/journeys-api/internal/app/journeys/model"
)

func NewJourneysDataStore(gtfsPath string, skipValidation bool) *JourneysDataStore {
	ctx := NewGTFSBundle(gtfsPath, skipValidation)

	allLines, linesById := buildLineIndexes(ctx.Routes)

	return &JourneysDataStore{
		Lines:     allLines,
		LinesById: linesById,
	}
}

type JourneysDataStore struct {
	Lines     []*model.Line
	LinesById map[string]*model.Line
}
