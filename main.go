package main

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit"
	"github.com/andrewthowell/budgit/budgit/clients"
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

	var provider Provider
	provider, err = clients.NewStarlingClient(config.Starling.URL, config.Starling.APIToken)
	if err != nil {
		panic(err)
	}

	accounts, err := provider.GetExternalAccounts(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%+v", accounts[0]))
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

type Provider interface {
	GetExternalAccounts(ctx context.Context) ([]*budgit.ExternalAccount, error)
	GetExternalAccount(ctx context.Context, externalID string) (*budgit.ExternalAccount, error)
}
