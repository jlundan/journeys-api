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
	"time"
)

const defaultPort = 8080
const defaultShortCacheDuration = 30 * time.Minute
const defaultLongCacheDuration = 2 * time.Hour
const defaultShortCacheLowerBound = 0
const defaultShortCacheUpperBound = 5

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

		serverPort, err := parseIntFromString(os.Getenv("JOURNEYS_PORT"), defaultPort)
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
			scLowerBound, err := parseIntFromString(os.Getenv("JOURNEYS_SHORT_CACHE_LOWER_BOUND"), defaultShortCacheLowerBound)
			if err != nil {
				log.Println(fmt.Sprintf("Error parsing short-cache lower bound: %s. Using default value: %v", err.Error(), defaultShortCacheLowerBound))
			}

			scUpperBound, err := parseIntFromString(os.Getenv("JOURNEYS_SHORT_CACHE_UPPER_BOUND"), defaultShortCacheUpperBound)
			if err != nil {
				log.Println(fmt.Sprintf("Error parsing short-cache upper bound: %s. Using default value: %v", err.Error(), defaultShortCacheUpperBound))
			}

			memcached, err := server.NewMemcachedCacheMiddleware(memcache.New(os.Getenv("MEMCACHED_URL")), getShortCacheDuration(), getLongCacheDuration(), scLowerBound, scUpperBound)
			if err != nil {
				log.Fatal(err)
			}

			flushErr := memcached.Flush()
			if flushErr != nil {
				log.Println(flushErr)
				os.Exit(1)
			}

			router.Use(memcached.Middleware)

			log.Println(fmt.Sprintf("Using cache. Short cache duration: %v, Long cache duration: %v. Short cache duration hours %v -> %v", getShortCacheDuration(), getLongCacheDuration(), scLowerBound, scUpperBound))
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

func parseIntFromString(source string, defaultValue int) (int, error) {
	var result int

	if source == "" {
		result = defaultValue
	} else {
		v, err := strconv.Atoi(source)
		if err != nil {
			return defaultValue, err
		}
		result = v
	}

	return result, nil
}

func getShortCacheDuration() time.Duration {
	cacheDurationStr := os.Getenv("JOURNEYS_SHORT_CACHE_DURATION")
	if cacheDurationStr == "" {
		return defaultShortCacheDuration
	}

	duration, err := time.ParseDuration(cacheDurationStr)
	if err != nil {
		fmt.Printf("Invalid JOURNEYS_SHORT_CACHE_DURATION format: %s, using default %v\n", cacheDurationStr, defaultShortCacheDuration)
		return defaultShortCacheDuration
	}

	return duration
}

func getLongCacheDuration() time.Duration {
	cacheDurationStr := os.Getenv("JOURNEYS_LONG_CACHE_DURATION")
	if cacheDurationStr == "" {
		return defaultLongCacheDuration
	}

	duration, err := time.ParseDuration(cacheDurationStr)
	if err != nil {
		fmt.Printf("Invalid JOURNEYS_LONG_CACHE_DURATION format: %s, using default %v\n", cacheDurationStr, defaultLongCacheDuration)
		return defaultLongCacheDuration
	}

	return duration
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
