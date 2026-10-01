package config

import (
	"flag"
)

var FlagRunAddr string
var FlagRunShorternerAddr string
var FlagLog bool

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "Адрес запуска http-сервера")

	flag.StringVar(&FlagRunShorternerAddr, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL ")

	flag.BoolVar(&FlagLog, "l", false, "Вывод логов в консоль")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
