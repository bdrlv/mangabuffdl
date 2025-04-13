package app

import (
	"log"
	"mbd/core"
	"mbd/sources"
	"path"
	"path/filepath"
)

func App(urlFlag *string, volumeFlag, startFlag, endFlag *int) error {
	parser, err := sources.NewParser(*urlFlag)
	if err != nil {
		log.Fatalf("ошибка инициализации: %v", err)
	}

	downloader := &core.HttpDownloader{
		Client: core.CreateDefaultClient(),
	}

	mangaName := parser.GetMangaName()
	if err := downloader.CreateDir(mangaName); err != nil {
		log.Fatalf("ошибка создания каталога тайтла: %v", err)
	}

	for chapterNum := *startFlag; chapterNum <= *endFlag; chapterNum++ {
		log.Printf("Обработка главы %d", chapterNum)

		chapterURL := parser.GetChapterURL(*volumeFlag, chapterNum)

		chapterInfo, err := parser.ParseChapter(chapterURL)
		if err != nil {
			log.Printf("Ошибка парсинга главы %d: %v", chapterNum, err)
			continue
		}

		chapterDir := path.Join(mangaName, core.FormatChapterDir(chapterNum))
		if err := downloader.CreateDir(chapterDir); err != nil {
			log.Printf("Ошибка создания каталога для главы %d: %v", chapterNum, err)
			continue
		}

		for _, page := range chapterInfo.Pages {
			fileName := core.FormatFileName(page.Number, page.FileExt)
			fullPath := filepath.Join(chapterDir, fileName)

			if err := downloader.DownloadImage(page.ImageURL, fullPath); err != nil {
				log.Printf("Ошибка скачивания страницы %d: %v", page.Number, err)
				continue
			}
		}
	}

	return nil
}
