package parser

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Komilov31/wget/internal/files"
)

var (
	ErrMaxDepthConceded    = errors.New("max depth of downloading exceeded")
	ErrAlreadyProcessedUrl = errors.New("this url is already processed")
)

type Parser struct {
	client        *http.Client
	maxDepth      int
	processedUrls sync.Map
	parsedUrls    sync.Map
}

func New(maxDepth int) *Parser {
	return &Parser{
		maxDepth:      maxDepth,
		processedUrls: sync.Map{},
		parsedUrls:    sync.Map{},
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (p *Parser) ProcessUrl(urlString string) error {
	url := getCorrectUrl(urlString)

	resp, err := p.client.Get(url)
	if err != nil {
		log.Fatal("unable to resolve host address:", url)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

	}()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(data)
	err = p.parse(url, reader)
	if err != nil {
		return err
	}

	err = files.SaveFile(url, data)
	if err != nil && err != ErrMaxDepthConceded {
		log.Fatal("could not save file:", err)
	}

	return nil
}
