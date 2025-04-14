package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"mbd/cmd/app"
	"mbd/core"
)

func main() {
	starttime := time.Now()
	defer func(t time.Time) {
		fmt.Printf("Затрачено времени: %v\n", time.Since(t))
	}(starttime)

	urlFlag := flag.String("u", "", "Ссылка на мангу")
	volumeFlag := flag.Int("v", 1, "Номер тома")
	startFlag := flag.Int("s", 1, "Номер начальной главы")
	endFlag := flag.Int("e", 1, "Номер конечной главы")
	parallelChaptersFlag := flag.Int("pc", 1, "Количество параллельно загружаемых глав")
	parallelPagesFlag := flag.Int("pp", 3, "Количество параллельно загружаемых страниц")
	parallelDelayFlag := flag.Int("d", 1000, "Задержка при параллельной загрузке, мс")
	flag.Parse()

	if *urlFlag == "" {
		log.Fatal("URL обязателен")
	}
	if *startFlag > *endFlag {
		log.Fatal("Начальный эпизод не может быть больше конечного")
	}

	engine := &core.ParallelEngine{
		Chapters: *parallelChaptersFlag,
		Pages:    *parallelPagesFlag,
		Delay:    time.Duration(*parallelDelayFlag) * time.Millisecond,
	}

	downloader, err := app.NewDownloaderApp(*urlFlag, engine)
	if err != nil {
		log.Fatal(err)
	}

	if err := downloader.Run(*volumeFlag, *startFlag, *endFlag); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Загрузка завершена!")
}
