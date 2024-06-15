package config

type Config struct {
	TelegramAPIToken string `env:"TELEGRAM_API_TOKEN"`
	PostgresURL      string `env:"POSTGRES_URL"`
	MongoURL         string `env:"MONGO_URL"`
}
