package core

import "time"

type ChapterInfo struct {
	Number int
	Pages  []Page
	Title  string
}

type Page struct {
	Number   int
	ImageURL string
	FileExt  string
}

type Downloader interface {
	DownloadImage(url string, path string) error
	CreateDir(path string) error
}

type Parser interface {
	ParseChapter(url string) (*ChapterInfo, error)
	GetChapterURL(volume, chapter int) string
	GetMangaName() string
}

type ParallelEngine struct {
	Chapters int
	Pages    int
	Delay    time.Duration
}
