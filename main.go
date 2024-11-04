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

func main() {
	if os.Getenv("PORT") == "" {
		fmt.Println("PORT is not set")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/addElements", tools.AddElements)
	mux.HandleFunc("/getAllElements", tools.GetAllElements)
	mux.HandleFunc("/listElements", tools.ListElements)
	mux.HandleFunc("/getElement", tools.GetElement)
	mux.HandleFunc("/listDatasets", tools.ListDatasets)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
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
