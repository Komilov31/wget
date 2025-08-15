package files

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func SaveFile(urlString string, data []byte) error {
	dirPath, filePath := parseURL(urlString)

	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		log.Fatal("could not create directory:", err)
	}

	file, err := os.Create(filePath)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal("could not close file:", err)
		}
	}()

	if err != nil {
		log.Fatal("could not create file:", err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatal("could not write to file: ", err)
	}

	fmt.Println("saved: ", filePath)
	return nil
}

func parseURL(input string) (string, string) {
	parsed, err := url.ParseRequestURI(input)
	if err != nil {
		return "", ""
	}

	host := GetHostFromUrl(parsed)
	localPath := parsed.Path

	dirPath := localPath
	file := localPath
	base := filepath.Base(localPath)
	if strings.Contains(base, ".") {
		dirPath = filepath.Dir(localPath)
	} else {
		if !strings.HasSuffix(localPath, "/") {
			file = file + "/"
		}
		file = file + "index.html"
	}

	dirFull := host + dirPath
	fileFull := host + file

	return dirFull, fileFull
}

func GetHostFromUrl(path *url.URL) string {
	host := path.Host
	host = strings.TrimPrefix(host, "https://www.")
	host = strings.TrimPrefix(host, "http://www.")
	host = strings.TrimPrefix(host, "www.")

	return host
}
