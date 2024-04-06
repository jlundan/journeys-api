package journeys

import (
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/context/tre"
	"github.com/jlundan/journeys-api/internal/app/journeys/routes"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func Run() {
	_ = mainCommand.Execute()
}

var mainCommand = &cobra.Command{
	Use:   "journeys",
	Short: "journeys",
	Long:  "Run the app",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv("JOURNEYS_BASE_URL") == "" {
			log.Println("JOURNEYS_BASE_URL not set in environment. Cannot proceed.")
			os.Exit(1)
		}

		ctx, errs, _ := tre.NewContext()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Println(err)
			}
			os.Exit(1)
		}

		r := mux.NewRouter()
		routes.InjectMunicipalityRoutes(r, ctx)
		routes.InjectLineRoutes(r, ctx)
		routes.InjectJourneyPatternRoutes(r, ctx)
		routes.InjectStopPointRoutes(r, ctx)
		routes.InjectJourneyRoutes(r, ctx)
		routes.InjectRouteRoutes(r, ctx)

		var port string
		if os.Getenv("JOURNEYS_PORT") == "" {
			port = "5678"
		}

		startServer(r, port)
	},
}
