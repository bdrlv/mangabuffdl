package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type MangaDownloader struct {
	Client     *http.Client
	UserAgent  string
	BaseDelay  time.Duration
	MaxRetries int
}

func NewMangaDownloader() *MangaDownloader {
	return &MangaDownloader{
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: false,
				MaxIdleConns:      10,
				IdleConnTimeout:   30 * time.Second,
			},
		},
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		BaseDelay:  1 * time.Second,
		MaxRetries: 3,
	}
}

func (d *MangaDownloader) DownloadManga(baseURL string, volume, startChapter, endChapter int) error {
	mangaURL := fmt.Sprintf("%s/%d", strings.TrimSuffix(baseURL, "/"), volume)
	mangaName := extractMangaName(baseURL)

	if err := os.MkdirAll(mangaName, 0750); err != nil {
		return fmt.Errorf("ошибка создания папки: %v", err)
	}

	for chapter := startChapter; chapter <= endChapter; chapter++ {
		fmt.Printf("\nТом %d, Глава %d...\n", volume, chapter)

		retries := 0
		for retries < d.MaxRetries {
			err := d.downloadChapter(mangaURL, mangaName, volume, chapter)
			if err == nil {
				break
			}

			retries++
			if retries >= d.MaxRetries {
				fmt.Printf("Превышено количество попыток для тома %d главы %d: %v\n", volume, chapter, err)
				break
			}

			delay := d.BaseDelay * time.Duration(retries)
			fmt.Printf("Ошибка в томе %d главе %d (попытка %d/%d): %v. Повтор через %v...\n",
				volume, chapter, retries, d.MaxRetries, err, delay)
			time.Sleep(delay)
		}

		time.Sleep(d.BaseDelay)
	}

	fmt.Println("\nСкачивание завершено!")
	return nil
}

func (d *MangaDownloader) downloadChapter(mangaURL, mangaName string, volume, chapter int) error {
	chapterURL := fmt.Sprintf("%s/%d", mangaURL, chapter)
	fmt.Printf("URL главы: %s\n", chapterURL)

	req, err := http.NewRequest("GET", chapterURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("User-Agent", d.UserAgent)

	resp, err := d.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("статус код: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка парсинга HTML: %v", err)
	}

	chapterPath := path.Join(mangaName, fmt.Sprintf("Chapter %d", chapter))
	if err := os.MkdirAll(chapterPath, 0750); err != nil {
		return fmt.Errorf("ошибка создания папки: %v", err)
	}

	doc.Find("div.reader__pages div.reader__item").Each(func(i int, item *goquery.Selection) {
		pageNum := item.AttrOr("data-page", strconv.Itoa(i+1))
		img := item.Find("img")
		if img.Length() == 0 {
			fmt.Printf("Том %d, Глава %d, страница %s: нет изображения\n", volume, chapter, pageNum)
			return
		}

		imgSrc := img.AttrOr("src", img.AttrOr("data-src", ""))
		if imgSrc == "" {
			fmt.Printf("Том %d, Глава %d, страница %s: пустая ссылка\n", volume, chapter, pageNum)
			return
		}

		imgSrc = strings.Split(imgSrc, "?")[0]
		ext := strings.ToLower(path.Ext(imgSrc))
		if ext == "" {
			fmt.Printf("Том %d, Глава %d, страница %s: неизвестное расширение\n", volume, chapter, pageNum)
			return
		}
		ext = ext[1:]

		filePath := path.Join(chapterPath, fmt.Sprintf("%s.%s", pageNum, ext))
		if err := d.downloadImage(imgSrc, filePath, mangaURL); err != nil {
			fmt.Printf("Том %d, Глава %d, страница %s: ошибка скачивания: %v\n", volume, chapter, pageNum, err)
		}
	})

	return nil
}

func (d *MangaDownloader) downloadImage(imgURL, filePath, referer string) error {
	if !strings.HasPrefix(imgURL, "http") {
		base, err := url.Parse(referer)
		if err != nil {
			return fmt.Errorf("ошибка парсинга базового URL: %v", err)
		}
		u, err := url.Parse(imgURL)
		if err != nil {
			return fmt.Errorf("ошибка парсинга URL изображения: %v", err)
		}
		imgURL = base.ResolveReference(u).String()
	}

	fmt.Printf("Скачиваем: %s -> %s\n", imgURL, filePath)

	req, err := http.NewRequest("GET", imgURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("User-Agent", d.UserAgent)
	req.Header.Set("Referer", referer)

	resp, err := d.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("статус код: %s", resp.Status)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %v", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("ошибка сохранения файла: %v", err)
	}

	fmt.Printf("Успешно: %s\n", filePath)
	return nil
}

func extractMangaName(mangaURL string) string {
	parts := strings.Split(mangaURL, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] != "" && parts[i] != "manga" {
			return parts[i]
		}
	}
	return "manga"
}

func main() {
	urlFlag := flag.String("u", "", "Базовый URL манги (например, https://mangabuff.ru/manga/ya-budu-korolem-v-etoi-zhizni)")
	volumeFlag := flag.Int("v", 1, "Номер тома")
	startFlag := flag.Int("s", 1, "Номер первой главы")
	endFlag := flag.Int("e", 1, "Номер последней главы")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Необходимо указать URL манги с помощью флага -u")
		flag.Usage()
		os.Exit(1)
	}

	downloader := NewMangaDownloader()
	if err := downloader.DownloadManga(*urlFlag, *volumeFlag, *startFlag, *endFlag); err != nil {
		fmt.Printf("Ошибка при скачивании манги: %v\n", err)
		os.Exit(1)
	}
}
