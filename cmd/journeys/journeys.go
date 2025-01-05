package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys"
	"github.com/jlundan/journeys-api/internal/app/journeys/context/tre"
	"github.com/jlundan/journeys-api/internal/app/journeys/routes"
	"github.com/spf13/cobra"
	"log"
	"os"
)

const defaultPort = 8080

var dryRun bool
var disableCache bool

func main() {
	MainCommand.Flags().BoolVar(&disableCache, "disable-cache", false, "Do not use cache")
	MainCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without starting the server")
	_ = MainCommand.Execute()
}

var MainCommand = &cobra.Command{
	Use:   "journeys",
	Short: "journeys",
	Long:  "Start Journey API server",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("JOURNEYS_BASE_URL") == "" {
			log.Fatal("JOURNEYS_BASE_URL not set in environment. Cannot proceed.")
		}

		if os.Getenv("JOURNEYS_GTFS_PATH") == "" {
			log.Fatal("JOURNEYS_GTFS_PATH not set in environment. Cannot proceed.")
		}

		serverPort, err := journeys.ParsePort(os.Getenv("JOURNEYS_PORT"), defaultPort)
		if err != nil {
			log.Fatal(err)
		}

		ctx := tre.NewContext(os.Getenv("JOURNEYS_GTFS_PATH"))

		parseErrors := ctx.GetParseErrors()
		for _, err := range parseErrors {
			log.Println(err)
		}

		violations := ctx.GetViolations()
		for _, v := range violations {
			log.Println(v)
		}

		recommendations := ctx.GetRecommendations()
		for _, r := range recommendations {
			log.Println(r)
		}

		infos := ctx.GetInfos()
		for infos := range infos {
			log.Println(infos)
		}

		if dryRun && len(parseErrors) > 0 && len(violations) > 0 && len(recommendations) > 0 {
			os.Exit(1)
		}

		router := mux.NewRouter()

		routes.InjectMunicipalityRoutes(router, ctx)
		routes.InjectLineRoutes(router, ctx)
		routes.InjectJourneyPatternRoutes(router, ctx)
		routes.InjectStopPointRoutes(router, ctx)
		routes.InjectJourneyRoutes(router, ctx)
		routes.InjectRouteRoutes(router, ctx)

		if !disableCache {
			cmw, err := journeys.NewMemcachedCacheMiddleware(memcache.New(os.Getenv("MEMCACHED_URL")))
			if err != nil {
				log.Fatal(err)
			}
			router.Use(cmw.Middleware)

			log.Println("Using cache")
		}

		router.Use(journeys.CorsMiddleware)

		journeys.StartServer(router, serverPort, onServerStartupSuccess, onServerStartupFailure, onServerShutdown)
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
