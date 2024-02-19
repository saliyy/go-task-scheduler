package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"task-scheduler/internal/app/apiserver"
	"task-scheduler/internal/app/apiserver/handlers/task/save"
	"task-scheduler/internal/app/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/apiserver.toml", "path to config")
}

func main() {

	flag.Parse()

	config := apiserver.MustLoad(configPath)

	storage, err := sqlite.New(config.StoragePath)

	if err != nil {
		log.Fatal("error to init db")
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/task", save.New(slog.Default(), storage))

	if err := http.ListenAndServe(config.BindAddr, router); err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal("Server is unvailable")
}
