package initializers

import "github.com/joho/godotenv"

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		logger.Info("Error loading .env file - possibly not found. Relying on environment variables.")
	}
}
