package worker

import (
	"log"
	"sync"

	"github.com/Komilov31/wget/internal/parser"
)

type Worker struct {
	maxWorkers int
	urls       chan string
	p          *parser.Parser
	startUrl   string
}

func New(startUrl string, maxWorkers int, p *parser.Parser) *Worker {
	return &Worker{
		startUrl:   startUrl,
		maxWorkers: maxWorkers,
		urls:       make(chan string),
		p:          p,
	}
}

func (w *Worker) Start() {
	wg := new(sync.WaitGroup)
	for i := 0; i < w.maxWorkers; i++ {
		wg.Add(1)
		go w.worker(wg)
	}
	w.urls <- w.startUrl

	wg.Add(1)
	go w.urlProcessor(wg)

	wg.Wait()
}

func (w *Worker) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range w.urls {
		err := w.p.ProcessUrl(url)
		if err != nil && err != parser.ErrAlreadyProcessedUrl {
			log.Println(err)
		}
	}
}

func (w *Worker) urlProcessor(wg *sync.WaitGroup) {
	defer wg.Done()
	err := w.p.ParseAllUrls(w.startUrl, w.urls, 0)
	if err != nil {
		if err != parser.ErrMaxDepthConceded && err != parser.ErrAlreadyProcessedUrl {
			log.Fatal(err)
		}
	}
	close(w.urls)
}
