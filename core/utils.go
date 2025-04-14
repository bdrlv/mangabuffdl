package core

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CreateDefaultClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}
}

func FormatChapterDir(chapter int) string {
	return fmt.Sprintf("Chapter_%04d", chapter)
}

func FormatFileName(number int, ext string) string {
	return fmt.Sprintf("%04d.%s", number, ext)
}

func ExtractMangaNameFromURL(rawURL string) string {
	// пока только для mangabuff
	u, err := url.Parse(rawURL)
	if err != nil {
		return "manga"
	}

	parts := strings.Split(u.Path, "/")
	// for i := range parts {
	// 	fmt.Println(parts[i])
	// }
	// panic(parts)
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i])
		if part != "" && part != "manga" {
			return part
			// return strings.ReplaceAll(part, "-", " ")
		}
	}
	return "manga"
}
