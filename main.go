package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stan-project/db"
	"stan-project/handler"
	"stan-project/logic"
	"syscall"
	"time"
)

func main() {

	ctx := context.Background()

	log.Printf("Initializing DB...")

	postgresDB, err := db.InitDB(ctx)
	if err != nil {
		panic(fmt.Sprintf("error initializing postgres DB: %s", err))
	}

	log.Printf("Running migrations...")

	err = postgresDB.RunMigrations(ctx)
	if err != nil {
		panic(fmt.Sprintf("error running migrations"))
	}

	riskDB := db.NewRisksDB(postgresDB)
	riskLogic := logic.NewRiskLogic(riskDB)
	riskHandler := handler.NewRiskHandler(riskLogic)

	log.Printf("Starting HTTP server...")

	h := handler.NewHandler(riskHandler)
	router := handler.NewRouter(h)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	terminationChannel := make(chan os.Signal, 1)
	signal.Notify(terminationChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Error starting HTTP server: %s", err))
		}
	}()

	sig := <-terminationChannel

	log.Printf("Termination signal '%s' received, initiating graceful shutdown...", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(25)*time.Second)
	defer cancel()

	shutdownGracefully(ctx, httpServer, postgresDB.Close)
}

func shutdownGracefully(ctx context.Context, httpServer *http.Server, postgresClose func(ctx context.Context) error) {
	//shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown HTTP server: %s", err)
	} else {
		log.Printf("successfully and gracefully shutdown HTTP server.")
	}

	err := postgresClose(ctx)
	if err != nil {
		log.Printf("failed to gracefully close postgres connection")
	} else {
		log.Printf("successfully and gracefully closed postgres connection")
	}

	log.Printf("Exiting Risks service.....")
}
