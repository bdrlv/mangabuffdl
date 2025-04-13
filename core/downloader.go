package core

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type HttpDownloader struct {
	Client *http.Client
}

func (d *HttpDownloader) DownloadImage(imgURL string, filePath string) error {
	parsedURL, err := url.ParseRequestURI(imgURL)
	if err != nil {
		return fmt.Errorf("неверный URL: %w", err)
	}

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return fmt.Errorf("ошибка формирования запроса: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")
	req.Header.Set("Referer", parsedURL.Host)

	startTime := time.Now()
	resp, err := d.Client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d %s",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("ошибка создания каталога: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("файл уже существует: %s", filePath)
		}
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(filePath)
		return fmt.Errorf("ошибка сохранения страницы: %w", err)
	}

	duration := time.Since(startTime)
	speed := float64(size) / duration.Seconds() / 1024 // KB/s

	// fmt.Printf("Загрузка: %s\n", filePath)
	// fmt.Printf("Size: %.2f KB, Time: %v, Speed: %.2f KB/s\n\n",
	// 	float64(size)/1024, duration.Round(time.Millisecond), speed)

	fmt.Printf("Загрузка: %s\n Size: %.2f KB, Time: %v, Speed: %.2f KB/s\n\n", filePath,
		float64(size)/1024, duration.Round(time.Millisecond), speed)

	return nil
}

func (d *HttpDownloader) CreateDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0750)
}
