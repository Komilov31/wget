package parser

import (
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/Komilov31/wget/internal/files"
	"github.com/PuerkitoBio/goquery"
)

func (p *Parser) parse(urlString string, resp io.Reader) error {
	urlPath, err := url.ParseRequestURI(urlString)
	if err != nil {
		return err
	}

	u := urlPath.Host + urlPath.Path
	_, ok := p.processedUrls.Load(u)
	if ok {
		return ErrAlreadyProcessedUrl
	}
	p.processedUrls.Store(u, struct{}{})

	p.saveAllResourses(urlPath, resp)

	return nil
}

func (p *Parser) saveAllResourses(parsedUrl *url.URL, resp io.Reader) {
	doc, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		return
	}

	// Найти все CSS
	doc.Find("link[rel='stylesheet']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			_, ok := p.processedUrls.Load(href)
			if ok {
				return
			}
			p.processedUrls.Store(href, struct{}{})

			url := href
			if !strings.Contains(href, "https://") {
				url = parsedUrl.Scheme + "://" + parsedUrl.Host + href
			}

			data, err := p.getData(url)
			if err != nil {
				log.Println("could not save CSS:", err)
				return
			}
			files.SaveFile(url, data)
		}
	})

	// Найти все JavaScript
	doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			_, ok := p.processedUrls.Load(src)
			if ok {
				return
			}
			p.processedUrls.Store(src, struct{}{})

			url := src
			if !strings.Contains(src, "https://") {
				url = parsedUrl.Scheme + "://" + parsedUrl.Host + src
			}

			data, err := p.getData(url)
			if err != nil {
				log.Println("could not save JS:", err)
				return
			}
			files.SaveFile(url, data)
		}
	})

	// Найти все изображения
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			_, ok := p.processedUrls.Load(src)
			if ok {
				return
			}
			p.processedUrls.Store(src, struct{}{})

			url := src
			if !strings.Contains(src, "https://") {
				url = parsedUrl.Scheme + "://" + parsedUrl.Host + src
			}

			data, err := p.getData(url)
			if err != nil {
				log.Println("could not save Image:", err)
				return
			}
			files.SaveFile(url, data)
		}
	})

	// Найти все icon
	doc.Find("link[rel='icon']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			_, ok := p.processedUrls.Load(href)
			if ok {
				return
			}
			p.processedUrls.Store(href, struct{}{})

			url := href
			if !strings.Contains(href, "https://") {
				url = parsedUrl.Scheme + "://" + parsedUrl.Host + href
			}

			data, err := p.getData(url)
			if err != nil {
				log.Println("could not save icon:", err)
				return
			}
			files.SaveFile(url, data)
		}
	})
}

// получить байты из url
func (p *Parser) getData(url string) ([]byte, error) {
	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// составить правильный url из входного
func getCorrectUrl(urlString string) string {
	urlString = strings.TrimSpace(urlString)

	if !strings.HasPrefix(urlString, "https://") {
		urlString = "https://" + urlString
	}

	if !strings.HasSuffix(urlString, "/") {
		urlString = urlString + "/"
	}

	return urlString
}

// парсинг всех ссылок из страницы
func (p *Parser) ParseAllUrls(urlString string, urls chan string, curDepth int) error {
	if p.maxDepth != 0 && curDepth >= p.maxDepth {
		return ErrMaxDepthConceded
	}

	urlString = getCorrectUrl(urlString)

	resp, err := p.client.Get(urlString)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	urlPath, err := url.ParseRequestURI(urlString)
	if err != nil {
		return err
	}

	h := urlPath.Host + urlPath.Path
	_, ok := p.parsedUrls.Load(h)
	if ok {
		return ErrAlreadyProcessedUrl
	}
	p.parsedUrls.Store(h, struct{}{})

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			path, err := url.ParseRequestURI(href)
			if err != nil {
				return
			}

			host := urlPath.Host
			if !strings.HasPrefix(host, "www.") {
				host = "www." + host
			}

			if path.Host == host {
				url := path.Host + path.Path

				_, ok := p.processedUrls.Load(url)
				if ok {
					return
				}

				urls <- href
				p.ParseAllUrls(href, urls, curDepth+1)
			}
		}
	})

	return nil
}
