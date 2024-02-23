package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"task-scheduler/internal/app/apiserver"
	"task-scheduler/internal/app/apiserver/handlers/task/list"
	"task-scheduler/internal/app/apiserver/handlers/task/save"
	"task-scheduler/internal/app/storage/sqlite"
	taskrepo "task-scheduler/internal/app/storage/sqlite/repos"

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

	storage, err := sqlite.New(config.StoragePath, config.DumpPath)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("error to init db")
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	taskRepo := taskrepo.New(storage)
	router.Post("/task", save.New(slog.Default(), taskRepo))
	router.Get("/tasks", list.New(slog.Default(), taskRepo))

	if err := http.ListenAndServe(config.BindAddr, router); err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal("Server is unvailable")
}
