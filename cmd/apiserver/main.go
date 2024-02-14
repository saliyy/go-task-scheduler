package main

import (
	"flag"
	"log"
	"os"
	"task-scheduler/internal/app/apiserver"
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

	_ = storage

	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("error to initializate db")
		os.Exit(1)
	}

	if err := s.Start(); err != nil {
		log.Fatal()
	}
}
