package main

import (
	"flag"
	"log"
	"task-scheduler/internal/app/apiserver"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config")
}

func main() {
	flag.Parse()

	config := apiserver.NewConifg()

	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal("No such config file")
	}

	s := apiserver.New(config)

	if err := s.Start(); err != nil {
		log.Fatal()
	}
}
