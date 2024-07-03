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
	Starling Client `required:"true" envconfig:"starling"`
}

type Client struct {
	URL      string `required:"true" envconfig:"url"`
	APIToken string `required:"true" envconfig:"api_token"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	config := &Config{}
	envconfig.MustProcess("", config)

	bearerToken := fmt.Sprintf("Bearer %s", config.Starling.APIToken)

	c, err := starling.NewClientWithResponses(config.Starling.URL, starling.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", bearerToken)
		return nil
	}))
	if err != nil {
		panic(err)
	}

	resp, err := c.GetAccountsWithResponse(context.Background())
	if err != nil {
		panic(err)
	}
	for _, account := range *resp.JSON200.Accounts {
		fmt.Println(fmt.Sprintf("%+v", account))
	}
}
