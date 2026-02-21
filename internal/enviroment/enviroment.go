package enviroment

import "github.com/joho/godotenv"

func LoadEnviroment(key_env string) {
	if key_env != "" {
		_ = godotenv.Load(".env." + key_env)
	} else {
		_ = godotenv.Load(".env")
	}
}
