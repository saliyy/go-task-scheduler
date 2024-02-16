package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"task-scheduler/internal/app/apiserver"
	"task-scheduler/internal/app/apiserver/handlers/task/save"
	"task-scheduler/internal/app/storage/sqlite"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/apiserver.toml", "path to config")
}

func main() {
	flag.Parse()

	config := apiserver.NewConifg()

	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal("No such config file")
	}

	s := apiserver.New(config)
	storage, err := sqlite.New(config.StoragePath)

	if err != nil {
		log.Fatal("error to initializate db")
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/task", save.New(slog.Default(), storage))

	if err := s.Start(); err != nil {
		log.Fatal()
	}
}
