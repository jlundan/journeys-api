package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jlundan/journeys-api/internal/app/journeys/repository"
	"github.com/jlundan/journeys-api/internal/app/journeys/router"
	"github.com/jlundan/journeys-api/internal/app/journeys/server"
	"github.com/jlundan/journeys-api/internal/app/journeys/service"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

const defaultPort = 8080

var dryRun bool
var disableCache bool
var skipValidation bool

func main() {
	MainCommand.Flags().BoolVar(&disableCache, "disable-cache", false, "Do not use cache")
	MainCommand.Flags().BoolVar(&skipValidation, "skip-validation", false, "Skip all validations")
	MainCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without starting the server")
	_ = MainCommand.Execute()
}

var MainCommand = &cobra.Command{
	Use:   "journeys",
	Short: "journeys",
	Long:  "Start Journey API server",
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

		//ctx := tre.NewContext(gtfsPath, skipValidation)
		//
		//for _, parseErrors := range ctx.GetParseErrors() {
		//	log.Println(parseErrors)
		//}
		//for _, v := range ctx.GetViolations() {
		//	log.Println(v)
		//}
		//for _, r := range ctx.GetRecommendations() {
		//	log.Println(r)
		//}
		//for _, i := range ctx.GetInfos() {
		//	log.Println(i)
		//}
		//
		//if dryRun {
		//	os.Exit(0)
		//}

		dataStore := repository.NewJourneysDataStore(gtfsPath, skipValidation)
		dataService := service.NewJourneysDataService(dataStore)

		r := router.New(dataService, baseUrl, vehicleActivityBaseUrl)

		if !disableCache {
			memcached, err := server.NewMemcachedCacheMiddleware(memcache.New(os.Getenv("MEMCACHED_URL")))
			if err != nil {
				log.Fatal(err)
			}
			r.Use(memcached.Middleware)

			log.Println("Using cache")
		}

		r.Use(server.CorsMiddleware)

		server.Start(r, serverPort, onServerStartupSuccess, onServerStartupFailure, onServerShutdown)
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
