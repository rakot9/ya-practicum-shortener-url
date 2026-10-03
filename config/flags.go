package config

import (
	"flag"
)

var FlagRunAddr string
var FlagRunShorternerAddr string
var FlagLog bool

func ParseFlags(args []string) {
	fs := flag.NewFlagSet("shortner", flag.ContinueOnError)

	fs.StringVar(&FlagRunAddr, "a", "localhost:8080", "Адрес запуска http-сервера")

	fs.StringVar(&FlagRunShorternerAddr, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL ")

	fs.BoolVar(&FlagLog, "l", false, "Вывод логов в консоль")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	fs.Parse(args)
}
