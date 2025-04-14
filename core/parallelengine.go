package core

import (
	"sync"
	"time"
)

func (e *ParallelEngine) ProcessChapters(totalChapters int, processChapter func(chapter int)) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, e.Chapters)

	for ch := 1; ch <= totalChapters; ch++ {
		sem <- struct{}{}
		wg.Add(1)

		go func(chapter int) {
			defer func() {
				<-sem
				wg.Done()
			}()
			processChapter(chapter)
		}(ch)
	}
	wg.Wait()
}

func (e *ParallelEngine) ProcessPages(pages []Page, processPage func(page Page)) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, e.Pages)

	for _, p := range pages {
		sem <- struct{}{}
		wg.Add(1)

		go func(page Page) {
			defer func() {
				<-sem
				wg.Done()
			}()
			time.Sleep(e.Delay)
			processPage(page)
		}(p)
	}
	wg.Wait()
}
