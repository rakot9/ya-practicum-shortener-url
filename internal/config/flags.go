package config

import (
	"flag"
	"os"
)

type Flags struct {
	FlagRunAddr           string
	FlagRunShorternerAddr string
	FlagFileStoragePath   string
	FlagInMemoryStorage   bool
	FlagDBDSN             string
}

func ParseFlags(args []string) Flags {
	fs := flag.NewFlagSet("shortner", flag.ContinueOnError)

	defaultAddr := "localhost:8080"
	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		defaultAddr = envAddr
	}

	defaultBaseURL := "http://localhost:8080"
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		defaultBaseURL = envBaseURL
	}

	defaultFileStorage := "storage_shorten_url.txt"
	if envFileStorage := os.Getenv("FILE_STORAGE_PATH"); envFileStorage != "" {
		defaultFileStorage = envFileStorage
	}

	defaultDBDSN := ""
	if envDBDSN := os.Getenv("DATABASE_DSN"); envDBDSN != "" {
		defaultDBDSN = envDBDSN
	}

	var flags Flags
	fs.StringVar(&flags.FlagRunAddr, "a", defaultAddr, "Адрес запуска http-сервера")
	fs.StringVar(&flags.FlagRunShorternerAddr, "b", defaultBaseURL, "Базовый адрес результирующего сокращённого URL")
	fs.StringVar(&flags.FlagFileStoragePath, "f", defaultFileStorage, "Путь к файлу хранилища")
	fs.BoolVar(&flags.FlagInMemoryStorage, "m", false, "Хранение в памяти")
	fs.StringVar(&flags.FlagDBDSN, "d", defaultDBDSN, "DSN подключения к БД")

	_ = fs.Parse(args)

	return flags
}
