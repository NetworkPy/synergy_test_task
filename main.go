package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NetworkPy/synergy_test_task/internal/handler"
	"github.com/NetworkPy/synergy_test_task/internal/model"
	"github.com/NetworkPy/synergy_test_task/internal/repository"
	"github.com/NetworkPy/synergy_test_task/internal/service"
)

func main() {
	methodUrlList := []model.MethodUrl{
		{
			Method: http.MethodGet,
			Url:    "https://novasite.su/test1.php",
		},
		{
			Method: http.MethodGet,
			Url:    "https://novasite.su/test2.php",
		},
	}

	cacheDataRepository := repository.NewCacheDataRepository()
	requestDataService, err := service.NewRequestDataService(&service.RDSConfig{
		Endpoints:           methodUrlList,
		CacheDataRepository: cacheDataRepository,
	})
	if err != nil {
		log.Fatal(err)
	}

	requestDataService.Start()

	mux := http.NewServeMux()

	handler.NewDataHandler(&handler.DHConfig{
		DataService: requestDataService,
		Mux:         mux,
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
