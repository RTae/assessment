package settings

import "os"

type Config struct {
	Port        string
	DatabaseUrl string
	Url         string
}

func Setting() Config {
	return Config{
		Port:        os.Getenv("PORT"),
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		Url:         "localhost",
	}
}
