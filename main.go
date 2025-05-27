package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	controllers "snake_ladder/controller"
	"snake_ladder/repository"
	"snake_ladder/service"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	userRepo := repository.NewUserRepo()
	gameRepo := repository.NewGameRepo()

	userService := service.NewUserService(userRepo)
	gameService := service.NewGameService(gameRepo, userService)
	matchmakingService := service.NewMatchMakingService(gameService)

	router := mux.NewRouter()

	router.HandleFunc("/playgames",controllers.WebsocketHandler(userService,matchmakingService,gameService))

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("Listening on port :8080")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	sig := <-sigChan
	log.Println("Recieved signal to terminate:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	log.Println("Server exited properly")
}
