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

/*
func ParseFlags(args []string) Flags {
	fs := flag.NewFlagSet("shortner", flag.ContinueOnError)

	flags := Flags{
		FlagRunAddr:           "",
		FlagRunShorternerAddr: "",
		FlagFileStoragePath:   "",
		FlagInMemoryStorage:   false,
		FlagDBDSN:             "",
	}

	fs.StringVar(&flags.FlagRunAddr, "a", "localhost:8080", "Адрес запуска http-сервера")

	fs.StringVar(&flags.FlagRunShorternerAddr, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL ")

	fs.StringVar(&flags.FlagFileStoragePath, "f", "storage_shorten_url.txt", "Вывод логов в консоль")

	fs.BoolVar(&flags.FlagInMemoryStorage, "m", false, "Хранение в памяти")

	fs.StringVar(&flags.FlagDBDSN, "d", "", "DSN подключения к БД")

	fs.Parse(args)

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		flags.FlagRunAddr = serverAddress
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		flags.FlagRunShorternerAddr = baseURL
	}

	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		flags.FlagFileStoragePath = fileStoragePath
	}

	if dbDSN := os.Getenv("DATABASE_DSN"); dbDSN != "" {
		flags.FlagDBDSN = dbDSN
	}

	return flags
}
*/

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
