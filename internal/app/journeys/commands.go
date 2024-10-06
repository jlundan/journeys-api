package journeys

import (
	"bytes"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/jlundan/journeys-api/internal/app/journeys/context/tre"
	"github.com/jlundan/journeys-api/internal/app/journeys/routes"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

var mc *memcache.Client
var dryRun bool

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// If it's a preflight request, stop here
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func init() {
	if os.Getenv("MEMCACHED_URL") != "" {
		mc = memcache.New(os.Getenv("MEMCACHED_URL"))
	}

	mainCommand.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without starting the server")
}

func cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.String()

		item, err := mc.Get(key)
		if err == nil {
			// Cache hit, send response
			_, err = w.Write(item.Value)
			if err != nil {
				rw := NewResponseWriter(w)
				next.ServeHTTP(rw, r)
			}
			return
		}

		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		// Cache the new response
		_ = mc.Set(&memcache.Item{Key: key, Value: rw.body.Bytes()})
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, new(bytes.Buffer)}
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

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
		if mc != nil {
			log.Println("Using memcached")
			r.Use(cacheMiddleware)
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
