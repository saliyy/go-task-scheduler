package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"task-scheduler/internal/app/apiserver"
	dto "task-scheduler/internal/app/dto/task"
	"task-scheduler/internal/app/storage/sqlite"

	"github.com/BurntSushi/toml"
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
		log.Fatal(err.Error())
		log.Fatal("error to initializate db")
		os.Exit(1)
	}

	taskDTO := &dto.CreateTaskDTO{Name: "Is My task", IsCompleted: false}
	entity, err := storage.SaveTask(taskDTO)

	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	fmt.Print(entity)

	if err := s.Start(); err != nil {
		log.Fatal()
	}
}
