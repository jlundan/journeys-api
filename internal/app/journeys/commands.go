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

func init() {
	mc = memcache.New("localhost:11211")
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

		ctx, errs, _ := tre.NewContext()
		if len(errs) > 0 {
			for _, err := range errs {
				log.Println(err)
			}
			os.Exit(1)
		}

		r := mux.NewRouter()
		r.Use(cacheMiddleware)
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
