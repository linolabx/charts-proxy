package config

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

type ChartFilter struct {
	Name         string `yaml:"name"`
	VersionsFrom string `yaml:"versions_from"`
	StableOnly   bool   `yaml:"stable_only"`
}

type Repo struct {
	Name         string        `yaml:"name"`
	Kind         string        `yaml:"kind"`
	Url          string        `yaml:"url"`
	Cron         string        `yaml:"cron"`
	ChartFilters []ChartFilter `yaml:"charts"`
}

type Config struct {
	TargetDir string `yaml:"target_dir"`
	Repos     []Repo `yaml:"repos"`

	DefaultProxy string `yaml:"default_proxy"`
}

func (c *Config) Validate() error {
	if c.TargetDir == "" {
		return fmt.Errorf("target_dir is required")
	}

	for _, repo := range c.Repos {
		if repo.Cron == "" {
			return fmt.Errorf("repo %s has no cron", repo.Name)
		}

		_, err := cron.Parse(repo.Cron)
		if err != nil {
			return fmt.Errorf("repo %s has invalid cron", repo.Name)
		}
	}

	return nil
}

func LoadConfig() *Config {
	godotenv.Load()

	config_file := os.Getenv("CONFIG_FILE")

	if config_file == "" {
		config_file = "/config.yaml"
	}

	yaml_file, err := os.ReadFile(config_file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(yaml_file, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
