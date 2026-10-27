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
}

func ParseFlags(args []string) Flags {
	fs := flag.NewFlagSet("shortner", flag.ContinueOnError)

	flags := Flags{
		FlagRunAddr:           "",
		FlagRunShorternerAddr: "",
		FlagFileStoragePath:   "",
		FlagInMemoryStorage:   false,
	}

	fs.StringVar(&flags.FlagRunAddr, "a", "localhost:8080", "Адрес запуска http-сервера")

	fs.StringVar(&flags.FlagRunShorternerAddr, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL ")

	fs.StringVar(&flags.FlagFileStoragePath, "f", "storage_shorten_url.txt", "Вывод логов в консоль")

	fs.BoolVar(&flags.FlagInMemoryStorage, "m", false, "Хранение в памяти")

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

	return flags
}
