package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ernestngugi/medvice-backend/internal/db"
	"github.com/ernestngugi/medvice-backend/internal/web/router"
	"github.com/joho/godotenv"
)

const defaultPort = "8088"

func main() {

	var envFilePath string
	flag.StringVar(&envFilePath, "e", "", "env file path")
	flag.Parse()

	if envFilePath != "" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			panic(fmt.Errorf("load env file err: %v", err))
		}
	}

	fmt.Printf("selected environment = [%v]", os.Getenv("ENVIRONMENT"))

	dB := db.InitDB()
	defer dB.Close()

	appRouter := router.BuildRouter(
		dB,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: appRouter,
	}

	done := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		fmt.Println("shutting down")

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server shut down error: %v", err)
		}

		close(done)
	}()

	fmt.Printf("web-api listening on :%v", port)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("Server shut down")
		} else {
			log.Fatal("Server shut down unexpectedly!")
		}
	}

	timeout := 30 * time.Second

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	code := 0
	select {
	case <-sigint:
		code = 1
		fmt.Println("Process forcibly terminated")
	case <-time.After(timeout):
		code = 1
		fmt.Println("Forcibly shutting down, shutdown timeout")
	case <-done:
		fmt.Println("Shutdown complete.")
	}

	fmt.Println("Server exiting")

	os.Exit(code)
}
