package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andrewthowell/budgit/integrations/starling"
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

	starlingClient, err := newStarlingClient(config.Starling)
	if err != nil {
		panic(err)
	}

	resp, err := starlingClient.GetAccountsWithResponse(context.Background())
	if err != nil {
		panic(fmt.Errorf("getting Accounts: %w", err))
	}
	if resp.JSON4XX != nil {
		panic(fmt.Errorf("getting Accounts: %+v", resp.JSON4XX))
	}
	for _, account := range *resp.JSON200.Accounts {
		fmt.Println(fmt.Sprintf("%+v", account))
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

func newStarlingClient(config *ClientConfig) (*starling.ClientWithResponses, error) {
	bearerToken := fmt.Sprintf("Bearer %s", config.APIToken)
	client, err := starling.NewClientWithResponses(
		config.URL,
		starling.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", bearerToken)
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("initialising starling client: %w", err)
	}
	return client, nil
}
