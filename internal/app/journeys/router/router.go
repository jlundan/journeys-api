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

	return router
}
