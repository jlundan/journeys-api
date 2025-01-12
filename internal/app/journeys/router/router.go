package router

import (
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/router/v1"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
)

func New(dataService service.DataService, baseUrl string) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/lines", v1.HandleGetAllLines(dataService, baseUrl)).Methods("GET")
	router.HandleFunc(`/v1/lines/{name}`, v1.HandleGetOneLine(dataService, baseUrl)).Methods("GET")
	router.HandleFunc("/v1/journeys", v1.HandleGetAllJourneys(dataService, baseUrl)).Methods("GET")
	router.HandleFunc(`/v1/journeys/{name}`, v1.HandleGetOneJourney(dataService, baseUrl)).Methods("GET")
	router.HandleFunc("/v1/journey-patterns", v1.HandleGetAllJourneyPatterns(dataService, baseUrl)).Methods("GET")
	router.HandleFunc(`/v1/journey-patterns/{name}`, v1.HandleGetOneJourneyPattern(dataService, baseUrl)).Methods("GET")
	router.HandleFunc("/v1/routes", v1.HandleGetAllRoutes(dataService, baseUrl)).Methods("GET")
	router.HandleFunc(`/v1/routes/{name}`, v1.HandleGetOneRoute(dataService, baseUrl)).Methods("GET")
	router.HandleFunc("/v1/stop-points", v1.HandleGetAllStopPoints(dataService, baseUrl)).Methods("GET")
	router.HandleFunc(`/v1/stop-points/{name}`, v1.HandleGetOneStopPoint(dataService, baseUrl)).Methods("GET")
	router.HandleFunc("/v1/municipalities", v1.HandleGetAllMunicipalities(dataService, baseUrl)).Methods("GET")
	router.HandleFunc(`/v1/municipalities/{name}`, v1.HandleGetOneMunicipality(dataService, baseUrl)).Methods("GET")

	return router
}
