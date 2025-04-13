package mangabuff

import (
	"fmt"
	"mbd/core"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type MangabuffParser struct {
	BaseMangaName string
	BaseURL       string
	UserAgent     string
	Client        *http.Client
	core.HttpDownloader
}

func (p *MangabuffParser) ParseChapter(url string) (*core.ChapterInfo, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", p.UserAgent)
	req.Header.Set("Referer", p.BaseURL)

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var chapterInfo core.ChapterInfo
	var pages []core.Page

	chapterTitle := doc.Find("h1.chapter-title").First().Text()
	chapterInfo.Title = strings.TrimSpace(chapterTitle)

	doc.Find("div.reader__pages div.reader__item").Each(func(i int, s *goquery.Selection) {
		pageNumber := i + 1
		if dataPage, exists := s.Attr("data-page"); exists {
			if num, err := strconv.Atoi(dataPage); err == nil {
				pageNumber = num
			}
		}

		img := s.Find("img").First()
		if img.Length() == 0 {
			return
		}

		imgSrc, _ := img.Attr("src")
		if imgSrc == "" {
			imgSrc, _ = img.Attr("data-src")
		}

		cleanedURL := strings.Split(imgSrc, "?")[0]
		ext := path.Ext(cleanedURL)
		if ext == "" {
			ext = "jpg"
		} else {
			ext = strings.TrimPrefix(ext, ".")
		}

		pages = append(pages, core.Page{
			Number:   pageNumber,
			ImageURL: cleanedURL,
			FileExt:  ext,
		})
	})

	if len(pages) == 0 {
		return nil, fmt.Errorf("no pages found")
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Number < pages[j].Number
	})

	chapterInfo.Pages = pages
	return &chapterInfo, nil
}

func (p *MangabuffParser) GetChapterURL(volume, chapter int) string {
	return fmt.Sprintf("%s/%d/%d", p.BaseURL, volume, chapter)
}

func (p *MangabuffParser) GetMangaName() string {
	return p.BaseMangaName
}
