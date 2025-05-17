package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gptscript-ai/datasets/pkg/tools"
)

var tokenFromEnv = os.Getenv("GPTSCRIPT_DAEMON_TOKEN")

func main() {
	if os.Getenv("PORT") == "" {
		fmt.Println("PORT is not set")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /addElements", authenticatedHandler(tools.AddElements))
	mux.HandleFunc("POST /getAllElements", authenticatedHandler(tools.GetAllElements))
	mux.HandleFunc("POST /listElements", authenticatedHandler(tools.ListElements))
	mux.HandleFunc("POST /getElement", authenticatedHandler(tools.GetElement))
	mux.HandleFunc("POST /listDatasets", authenticatedHandler(tools.ListDatasets))
	mux.HandleFunc("POST /outputFilter", authenticatedHandler(tools.OutputFilter))
	mux.HandleFunc("/", health)

	srv := &http.Server{
		Addr:    "127.0.0.1:" + os.Getenv("PORT"),
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error: %v\n", err)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("server forced to shut down: %v\n", err)
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}

func authenticatedHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authenticate(r.Header) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func authenticate(headers http.Header) bool {
	return headers.Get("X-GPTScript-Daemon-Token") == tokenFromEnv
}
