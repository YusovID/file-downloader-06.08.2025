package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                string   `yaml:"env" env-required:"true"`
	MaxFilesPerTask    int      `yaml:"max_files_per_task" env=required:"true"`
	MaxConcurrentTasks int      `yaml:"max_concurrent_tasks" env=required:"true"`
	AllowedExtensions  []string `yaml:"allowed_extensions" env=required:"true"`
	TempDir            string   `yaml:"temp_dir" env=required:"true"`
	HTTPServer         `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %v", err)
	}

	return &cfg
}
