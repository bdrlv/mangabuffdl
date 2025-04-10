package main

import (
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

type Title struct {
	Name         string
	BaseURL      string
	Volume       int
	ChapterFirst int
	ChapterLast  int
	SlidesLimit  int
	Formats      []string
	Client       *http.Client
}

func main() {
	manga := &Title{
		Name:         "ya-budu-korolem-v-etoi-zhizni",
		BaseURL:      "https://mangabuff.ru/manga/ya-budu-korolem-v-etoi-zhizni/1",
		Volume:       1,
		ChapterFirst: 1,
		ChapterLast:  114,
		SlidesLimit:  500,
		Formats:      []string{"jpeg", "jpg", "png", "webp"},
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: false,
			},
		},
	}
	manga.Run()
}

func (t *Title) Run() error {
	for chapter := t.ChapterFirst; chapter <= t.ChapterLast; chapter++ {
		fmt.Printf("\nГлава %d...\n", chapter)
		err := t.downloadChapter(chapter)
		if err != nil {
			fmt.Printf("Ошибка в главе %d: %v\n", chapter, err)
			continue
		}
		time.Sleep(1 * time.Second) // Задержка между главами
	}
	fmt.Println("\nСкачивание завершено!")
	return nil
}

func (t *Title) downloadChapter(chapter int) error {
	chapterURL := fmt.Sprintf("%s/%d", t.BaseURL, chapter)
	fmt.Printf("URL главы: %s\n", chapterURL)

	// Получаем HTML страницы главы
	req, err := http.NewRequest("GET", chapterURL, nil)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := t.Client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("статус код: %s", resp.Status)
	}

	// Парсим HTML с помощью goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка парсинга HTML: %v", err)
	}

	// Создаем папку для главы
	chapterPath := path.Join(t.Name, fmt.Sprintf("Chapter %d", chapter))
	if err := os.MkdirAll(chapterPath, 0750); err != nil {
		return fmt.Errorf("ошибка создания папки: %v", err)
	}

	// Ищем ВСЕ изображения в reader__pages
	doc.Find("div.reader__pages div.reader__item").Each(func(i int, item *goquery.Selection) {
		pageNum := item.AttrOr("data-page", strconv.Itoa(i+1))
		img := item.Find("img")
		if img.Length() == 0 {
			fmt.Printf("Глава %d, страница %s: нет изображения\n", chapter, pageNum)
			return
		}

		// Получаем src (или data-src, если src пустой)
		imgSrc := img.AttrOr("src", img.AttrOr("data-src", ""))
		if imgSrc == "" {
			fmt.Printf("Глава %d, страница %s: пустая ссылка\n", chapter, pageNum)
			return
		}

		// Очищаем URL от параметров (все после ?)
		imgSrc = strings.Split(imgSrc, "?")[0]

		// Проверяем расширение файла
		ext := strings.ToLower(path.Ext(imgSrc))
		if ext == "" {
			fmt.Printf("Глава %d, страница %s: неизвестное расширение\n", chapter, pageNum)
			return
		}
		ext = ext[1:] // Убираем точку в начале

		// Проверяем поддерживаемый формат
		validFormat := false
		for _, format := range t.Formats {
			if ext == format {
				validFormat = true
				break
			}
		}
		if !validFormat {
			fmt.Printf("Глава %d, страница %s: неподдерживаемый формат (%s)\n", chapter, pageNum, ext)
			return
		}

		// Скачиваем изображение
		filePath := path.Join(chapterPath, fmt.Sprintf("%s.%s", pageNum, ext))
		if err := t.downloadImage(imgSrc, filePath); err != nil {
			fmt.Printf("Глава %d, страница %s: ошибка скачивания: %v\n", chapter, pageNum, err)
		}
	})

	return nil
}

func (t *Title) downloadImage(imgURL, filePath string) error {
	// Проверяем, является ли URL абсолютным
	if !strings.HasPrefix(imgURL, "http") {
		base, err := url.Parse(t.BaseURL)
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
	req.Header.Set("Referer", t.BaseURL)

	resp, err := t.Client.Do(req)
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
