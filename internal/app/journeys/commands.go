package journeys

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/context/tre"
	"github.com/jlundan/journeys-api/internal/app/journeys/routes"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var dryRun bool
var disableCache bool

func init() {
	MainCommand.Flags().BoolVar(&disableCache, "disable-cache", false, "Do not use cache")
	MainCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without starting the server")
}

var MainCommand = &cobra.Command{
	Use:   "journeys",
	Short: "journeys",
	Long:  "Run the app",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("JOURNEYS_BASE_URL") == "" {
			log.Println("JOURNEYS_BASE_URL not set in environment. Cannot proceed.")
			os.Exit(1)
		}

		ctx, errs, warnings, recommendations := tre.NewContext()

		for _, err := range errs {
			log.Println(err.Error())
		}

		for _, warning := range warnings {
			log.Println(warning.Error())
		}

		for _, recommendation := range recommendations {
			log.Println(recommendation)
		}

		if dryRun && len(errs) > 0 {
			os.Exit(0)
		}

		r := mux.NewRouter()
		r.Use(corsMiddleware)

		if !disableCache {
			cmw, err := NewMemcachedCacheMiddleware(memcache.New(os.Getenv("MEMCACHED_URL")))
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("Using cache")
			r.Use(cmw.Middleware)
		}

		routes.InjectMunicipalityRoutes(r, ctx)
		routes.InjectLineRoutes(r, ctx)
		routes.InjectJourneyPatternRoutes(r, ctx)
		routes.InjectStopPointRoutes(r, ctx)
		routes.InjectJourneyRoutes(r, ctx)
		routes.InjectRouteRoutes(r, ctx)

		var port string
		if val := os.Getenv("JOURNEYS_PORT"); val == "" {
			port = "8080"
		} else {
			port = val
		}

		startServer(r, port)
	},
}
