package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	v1 "github.com/jlundan/journeys-api/internal/app/journeys/handlers/v1"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/app/journeys/server"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

const defaultPort = 8080

// This variable is set at build time
//
//goland:noinspection GoUnusedGlobalVariable
var version = "dev"

var dryRun bool
var disableCache bool
var skipValidation bool

var MainCommand = &cobra.Command{
	Use: "journeys",
}

var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Start the API server",
	Long:  "Start the API server",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		baseUrl := os.Getenv("JOURNEYS_BASE_URL")
		vehicleActivityBaseUrl := os.Getenv("JOURNEYS_VA_BASE_URL")
		gtfsPath := os.Getenv("JOURNEYS_GTFS_PATH")

		if baseUrl == "" {
			log.Fatal("JOURNEYS_BASE_URL not set in environment. Cannot proceed.")
		}

		if gtfsPath == "" {
			log.Fatal("JOURNEYS_GTFS_PATH not set in environment. Cannot proceed.")
		}

		serverPort, err := parsePort(os.Getenv("JOURNEYS_PORT"), defaultPort)
		if err != nil {
			log.Fatal(err)
		}

		dataStore, errs := repository.NewJourneysRepository(gtfsPath, skipValidation)

		for _, e := range errs {
			log.Println(e)
		}

		if dryRun {
			os.Exit(0)
		}

		router := mux.NewRouter()

		router.Use(server.CorsMiddleware)

		if !disableCache {
			memcached, err := server.NewMemcachedCacheMiddleware(memcache.New(os.Getenv("MEMCACHED_URL")))
			if err != nil {
				log.Fatal(err)
			}
			router.Use(memcached.Middleware)

			log.Println("Using cache")
		}

		dataService := service.NewJourneysDataService(dataStore)

		router.HandleFunc("/v1/lines", v1.HandleGetAllLines(dataService, baseUrl)).Methods("GET")
		router.HandleFunc(`/v1/lines/{name}`, v1.HandleGetOneLine(dataService, baseUrl)).Methods("GET")
		router.HandleFunc("/v1/journeys", v1.HandleGetAllJourneys(dataService, baseUrl, vehicleActivityBaseUrl)).Methods("GET")
		router.HandleFunc(`/v1/journeys/{name}`, v1.HandleGetOneJourney(dataService, baseUrl, vehicleActivityBaseUrl)).Methods("GET")
		router.HandleFunc("/v1/journey-patterns", v1.HandleGetAllJourneyPatterns(dataService, baseUrl)).Methods("GET")
		router.HandleFunc(`/v1/journey-patterns/{name}`, v1.HandleGetOneJourneyPattern(dataService, baseUrl)).Methods("GET")
		router.HandleFunc("/v1/routes", v1.HandleGetAllRoutes(dataService, baseUrl)).Methods("GET")
		router.HandleFunc(`/v1/routes/{name}`, v1.HandleGetOneRoute(dataService, baseUrl)).Methods("GET")
		router.HandleFunc("/v1/stop-points", v1.HandleGetAllStopPoints(dataService, baseUrl)).Methods("GET")
		router.HandleFunc(`/v1/stop-points/{name}`, v1.HandleGetOneStopPoint(dataService, baseUrl)).Methods("GET")
		router.HandleFunc("/v1/municipalities", v1.HandleGetAllMunicipalities(dataService, baseUrl)).Methods("GET")
		router.HandleFunc(`/v1/municipalities/{name}`, v1.HandleGetOneMunicipality(dataService, baseUrl)).Methods("GET")

		server.Start(router, serverPort, onServerStartupSuccess, onServerStartupFailure, onServerShutdown)
	},
}

func onServerStartupSuccess(port int) {
	log.Println(fmt.Sprintf("listening on port %v", port))
}

func onServerStartupFailure(err error) {
	log.Println(err)
}

func onServerShutdown() {
	log.Println("shutting down")
}

func parsePort(envPort string, defaultPort int) (int, error) {
	var serverPort int

	if envPort == "" {
		serverPort = defaultPort
	} else {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			return 0, err
		}
		serverPort = p
	}

	return serverPort, nil
}

func main() {
	StartCommand.Flags().BoolVar(&disableCache, "disable-cache", false, "Do not use cache")
	StartCommand.Flags().BoolVar(&skipValidation, "skip-validation", false, "Skip all validations")
	StartCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without starting the server")

	MainCommand.AddCommand(StartCommand)
	MainCommand.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:", version)
		},
	})
	_ = MainCommand.Execute()
}
