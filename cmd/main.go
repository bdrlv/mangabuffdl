package main

import (
	"flag"
	"fmt"
	"log"

	"mbd/cmd/app"
)

func main() {
	urlFlag := flag.String("u", "", "Ссылка на мангу")
	volumeFlag := flag.Int("v", 1, "Номер тома")
	startFlag := flag.Int("s", 1, "Номер начальной главы")
	endFlag := flag.Int("e", 1, "Номер конечной главы")
	flag.Parse()

	if *urlFlag == "" {
		log.Fatal("URL обязателен")
	}
	if *startFlag > *endFlag {
		log.Fatal("Начальный эпизод не может быть больше конечного")
	}

	err := app.App(urlFlag, volumeFlag, startFlag, endFlag)
	if err != nil {
		log.Fatalf("ошибка загрузки: %s", err)
	}

	fmt.Println("Загрузка завершена!")
}
