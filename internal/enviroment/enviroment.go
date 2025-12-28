package enviroment

import "github.com/joho/godotenv"

func LoadEnviroment(key_env string) {
	_ = godotenv.Load(".env." + key_env)
}
