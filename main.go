package main

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit/clients"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/svc"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Starling *ClientConfig `required:"true" envconfig:"starling"`
}

type ClientConfig struct {
	URL      string `required:"true" envconfig:"url"`
	APIToken string `required:"true" envconfig:"api_token"`
}

func main() {
	config, err := loadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	budgetID := "main"

	starlingClient, err := clients.NewStarlingClient(config.Starling.URL, config.Starling.APIToken)
	if err != nil {
		panic(err)
	}

	db := db.New()
	service := svc.New(db, map[string]svc.Provider{starlingClient.ID(): starlingClient})

	if err = service.LoadAccountsFromProvider(context.Background(), budgetID, starlingClient.ID()); err != nil {
		panic(err)
	}
}

func loadConfigFromEnv() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("loading config from env: %w", err)
	}

	config := &Config{}
	envconfig.MustProcess("", config)
	return config, nil
}
