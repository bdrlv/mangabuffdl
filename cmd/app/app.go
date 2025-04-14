package app

import (
	"log"
	"mbd/core"
	"mbd/sources"
	"path/filepath"
)

type App struct {
	parser     core.Parser
	downloader *core.HttpDownloader
	engine     *core.ParallelEngine
	mangaName  string
}

func NewDownloaderApp(url string, engine *core.ParallelEngine) (*App, error) {
	parser, err := sources.NewParser(url)
	if err != nil {
		return nil, err
	}

	return &App{
		parser:     parser,
		downloader: &core.HttpDownloader{Client: core.CreateDefaultClient()},
		engine:     engine,
		mangaName:  parser.GetMangaName(),
	}, nil
}

func (a *App) Run(volume, startChapter, endChapter int) error {
	// Создаем корневую папку
	if err := a.downloader.CreateDir(a.mangaName); err != nil {
		return err
	}

	// Обрабатываем главы параллельно
	a.engine.ProcessChapters(endChapter-startChapter+1, func(chapter int) {
		chapterNum := startChapter + chapter - 1
		log.Printf("Обработка главы %d", chapterNum)

		chapterURL := a.parser.GetChapterURL(volume, chapterNum)
		chapterInfo, err := a.parser.ParseChapter(chapterURL)
		if err != nil {
			log.Printf("Ошибка парсинга главы %d: %v", chapterNum, err)
			return
		}

		chapterDir := filepath.Join(a.mangaName, core.FormatChapterDir(chapterNum))
		if err := a.downloader.CreateDir(chapterDir); err != nil {
			log.Printf("Ошибка создания каталога главы %d: %v", chapterNum, err)
			return
		}

		// Обрабатываем страницы параллельно с задержкой
		a.engine.ProcessPages(chapterInfo.Pages, func(page core.Page) {
			fileName := core.FormatFileName(page.Number, page.FileExt)
			fullPath := filepath.Join(chapterDir, fileName)

			if err := a.downloader.DownloadImage(page.ImageURL, fullPath); err != nil {
				log.Printf("Ошибка скачивания страницы %d: %v", page.Number, err)
			}
		})
	})

	return nil
}
