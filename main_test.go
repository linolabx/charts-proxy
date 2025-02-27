package main_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/linolabx/charts-proxy/config"
	"github.com/linolabx/charts-proxy/syncer"
)

func TestSync(t *testing.T) {
	config := config.LoadConfig()

	for _, repo := range config.Repos {
		err := syncer.Sync(&repo, config)
		if err != nil {
			log.Fatalf("Error syncing repo %s: %s", repo.Name, err)
		}
	}

	fmt.Println(config)

}
