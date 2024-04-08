package usecase

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

type IDownloaderService interface {
	ProcessDownload(url []string, dirPath string)
}

type EventServiceImpl struct {
}

func (impl *EventServiceImpl) ProcessDownload(urls []string, dirPath string) {

	done := make(chan bool, len(urls))
	errch := make(chan error, len(urls))

	for _, URL := range urls {

		go func(URL string) {
			b, err := download(URL, dirPath)
			if err != nil {
				errch <- err
				done <- false
				return
			}
			done <- b
			errch <- nil
		}(URL)
	}
}

func download(url string, dirPath string) (bool, error) {
	fileName, err := extractFileNameFromURL(url)
	if err != nil {
		return false, err
	}

	filePath := path.Join(dirPath, fileName)

	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, errors.New(resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return false, err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}
	return true, nil
}

func extractFileNameFromURL(fileURL string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}
	return path.Base(parsedURL.Path), nil
}
