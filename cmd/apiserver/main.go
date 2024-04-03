package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"task-scheduler/internal/app/apiserver"
	"task-scheduler/internal/app/apiserver/handlers/auth/oauth"
	"task-scheduler/internal/app/apiserver/handlers/list"
	task "task-scheduler/internal/app/apiserver/handlers/task/list"
	"task-scheduler/internal/app/apiserver/handlers/task/save"
	"task-scheduler/internal/app/apiserver/middlewares/auth"
	"task-scheduler/internal/app/events"
	"task-scheduler/internal/app/listeners"
	"task-scheduler/internal/app/storage/sqlite"
	listrepo "task-scheduler/internal/app/storage/sqlite/repos/list"
	taskrepo "task-scheduler/internal/app/storage/sqlite/repos/task"
	userrepo "task-scheduler/internal/app/storage/sqlite/repos/user"

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
		log.Fatal("error to init db %w", err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	taskRepo := taskrepo.New(storage)
	listRepo := listrepo.New(storage)

	// todo move to factory
	userCreatedListener := listeners.NewDefaultListCreator(listRepo)

	listeners := &listeners.Listeners{
		CreateDefaultListListener: userCreatedListener,
	}

	events.Init(listeners)

	router.Route("/tasks", func(r chi.Router) {
		r.Use(auth.CurrentUserCtx)
		r.Post("/", save.New(slog.Default(), taskRepo))
		r.Get("/", task.New(slog.Default(), taskRepo))
	})

	router.Route("/lists", func(r chi.Router) {
		r.Use(auth.CurrentUserCtx)
		r.Post("/", list.New(slog.Default(), listRepo))
		// add user to list (sending email you are added to list)
	})

	userRepo := userrepo.New(storage)
	router.Post("/oauth/signup", oauth.New(slog.Default(), userRepo))
	router.Post("/oauth/authorize", oauth.Authorize(slog.Default(), userRepo))

	if err := http.ListenAndServe(config.BindAddr, router); err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal("Server is unvailable")
}
