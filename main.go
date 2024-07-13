package main

import (
	"context"
	"fmt"

	"github.com/andrewthowell/budgit/budgit/clients"
	"github.com/andrewthowell/budgit/budgit/db"
	"github.com/andrewthowell/budgit/budgit/svc"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB       *DBConfig     `required:"true" envconfig:"db"`
	Starling *ClientConfig `required:"true" envconfig:"starling"`
}

type ClientConfig struct {
	URL      string `required:"true" envconfig:"url"`
	APIToken string `required:"true" envconfig:"api_token"`
}

type DBConfig struct {
	User     string `required:"true" envconfig:"user"`
	Password string `required:"true" envconfig:"password"`
	Host     string `required:"true" envconfig:"host"`
	Port     string `required:"true" envconfig:"port"`
}

func main() {
	config, err := loadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s?sslmode=disable", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	starlingClient, err := clients.NewStarlingClient(config.Starling.URL, config.Starling.APIToken)
	if err != nil {
		panic(err)
	}

	service := svc.New(conn, db.DB{}, map[string]svc.Provider{starlingClient.ID(): starlingClient})

	accounts, externalAccounts, err := service.LoadAccountsFromProvider(context.Background(), starlingClient.ID())
	if err != nil {
		panic(err)
	}
	fmt.Println(len(accounts), "accounts")
	for _, account := range accounts {
		fmt.Println(fmt.Sprintf("%+v", account))
	}
	fmt.Println(len(externalAccounts), "externalAccounts")
	for _, account := range externalAccounts {
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
