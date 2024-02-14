package apiserver

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	StoragePath string `toml:"storage_path"`
}

func NewConifg() *Config {
	return &Config{
		BindAddr:    ":8080",
		LogLevel:    "debug",
		StoragePath: "tasks.sqlite3",
	}
}
