package main

import (
	"log"

	"github.com/linolabx/charts-proxy/config"
	"github.com/linolabx/charts-proxy/syncer"
	"github.com/robfig/cron/v3"
)

func main() {
	config := config.LoadConfig()

	if err := config.Validate(); err != nil {
		log.Fatalf("Error validating config: %s", err)
	}

	c := cron.New()
	for _, repo := range config.Repos {
		if err := syncer.Sync(&repo, config); err != nil {
			log.Fatalf("Error syncing repo %s: %s", repo.Name, err)
		}

		_, err := c.AddFunc(repo.Cron, func() {
			err := syncer.Sync(&repo, config)
			if err != nil {
				log.Fatalf("Error syncing repo %s: %s", repo.Name, err)
			}
		})

		if err != nil {
			log.Fatalf("Error adding cron job for repo %s: %s", repo.Name, err)
		}

		log.Printf("Added cron job for repo %s: %s", repo.Name, repo.Cron)
	}

	c.Start()

	log.Printf("Started cron jobs")
	select {}
}
