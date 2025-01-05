package journeys

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(r *mux.Router, port int, onStartupSuccess func(port int), onStartupError func(err error), onShutdown func()) {
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%v", port),
		// Avoid Slow loris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handlers.CompressHandler(r),
	}

	go func() {
		if onStartupSuccess != nil {
			onStartupSuccess(port)
		}
		if err := srv.ListenAndServe(); err != nil {
			if onStartupError != nil {
				onStartupError(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_ = srv.Shutdown(ctx)

	if onShutdown != nil {
		onShutdown()
	}

	os.Exit(0)
}
