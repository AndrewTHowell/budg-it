package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/andrewthowell/budgit/integrations/starling"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	apiToken, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		panic("Missing `API_TOKEN` env variable")
	}
	bearerToken := fmt.Sprintf("Bearer %s", apiToken)

	c, err := starling.NewClientWithResponses("https://api.starlingbank.com", starling.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
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
