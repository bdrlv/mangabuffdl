package sources

import (
	"errors"
	"mbd/core"
	"mbd/sources/mangabuff"
	"net/url"
)

func NewParser(sourceURL string) (core.Parser, error) {

	mangaName := core.ExtractMangaNameFromURL(sourceURL)

	u, _ := url.Parse(sourceURL)

	switch u.Host {
	case "mangabuff.ru":
		return &mangabuff.MangabuffParser{
			BaseMangaName: mangaName,
			BaseURL:       sourceURL,
			UserAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36...",
			Client:        core.CreateDefaultClient(),
		}, nil
	}

	return nil, errors.New("источник не поддерживается")
}
