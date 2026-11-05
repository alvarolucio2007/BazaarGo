package app

import "github.com/joho/godotenv"

func LoadEnv(envfile string) error {
	return godotenv.Load(envfile)
}
