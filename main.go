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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	DB       *DBConfig     `required:"true" envconfig:"db"`
	Logger   *LoggerConfig `required:"true" envconfig:"logger"`
	Starling *ClientConfig `required:"true" envconfig:"starling"`
}

func (c Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("DB", c.DB)
	enc.AddObject("Logger", c.Logger)
	enc.AddObject("Starling", c.Starling)
	return nil
}

type DBConfig struct {
	User     string `required:"true" envconfig:"user"`
	Password string `required:"true" envconfig:"password"`
	Host     string `required:"true" envconfig:"host"`
	Port     string `required:"true" envconfig:"port"`
}

func (c DBConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("User", c.User)
	enc.AddString("Password", "**REDACTED**")
	enc.AddString("Host", c.Host)
	enc.AddString("Port", c.Port)
	return nil
}

type LoggerConfig struct {
	IsDev bool `required:"true" envconfig:"is_dev"`
}

func (c LoggerConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddBool("IsDev", c.IsDev)
	return nil
}

type ClientConfig struct {
	URL      string `required:"true" envconfig:"url"`
	APIToken string `required:"true" envconfig:"api_token"`
}

func (c ClientConfig) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("URL", c.URL)
	enc.AddString("APIToken", "**REDACTED**")
	return nil
}

func main() {
	config, err := loadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	log, err := newLogger(config)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Infow("Starting Budgit", zap.Any("config", config))

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/budgit?sslmode=disable", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Panic("Connecting to Postgres", zap.Error(err))
	}
	defer conn.Close(context.Background())

	starlingClient, err := clients.NewStarlingClient(log, config.Starling.URL, config.Starling.APIToken)
	if err != nil {
		log.Panic("Connecting to Starling", zap.Error(err))
	}

	db := db.New(log)
	service := svc.New(log, conn, db, []svc.Integration{starlingClient})

	accounts, err := service.LoadAccountsFromIntegration(context.Background(), starlingClient.ID())
	if err != nil {
		log.Panic("Loading accounts from Starling", zap.Error(err))
	}
	fmt.Println(len(accounts), "accounts")
	for _, account := range accounts {
		fmt.Println(fmt.Sprintf("%+v", account))
	}

	log.Info("Exiting Budgit")
}

func newLogger(config *Config) (*zap.SugaredLogger, error) {
	cfg, encoderCfg := zap.NewProductionConfig(), zap.NewProductionEncoderConfig()
	if config.Logger.IsDev {
		cfg, encoderCfg = zap.NewDevelopmentConfig(), zap.NewDevelopmentEncoderConfig()
	}

	cfg.EncoderConfig = encoderCfg

	l, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("starting logger: %w", err)
	}
	return l.Sugar(), nil
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
