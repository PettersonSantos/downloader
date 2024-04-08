package usecase

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
)

type IDownloaderService interface {
	ProcessDownload(url []string, dirPath string)
}

type EventServiceImpl struct {
}

func (impl *EventServiceImpl) ProcessDownload(urls []string, dirPath string) {

	errch := make(chan error, len(urls))

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, URL := range urls {
		go func(URL string) {
			defer wg.Done()
			err := download(URL, dirPath)
			if err != nil {
				errch <- err
				return
			}
		}(URL)
	}

	var qtdErr int
	for err := range errch {
		if err != nil {
			qtdErr++
			fmt.Println("Erro ao baixar arquivo:", err)
		}
	}
	fmt.Printf("%d/%d downloads concluídos com sucesso\n", len(urls)-qtdErr, len(urls))

	wg.Wait()
	close(errch)
}

func download(url string, dirPath string) error {
	fileName, err := extractFileNameFromURL(url)
	if err != nil {
		return err
	}

	filePath := path.Join(dirPath, fileName)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	totalSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	err = showProgress(resp, totalSize, out, url)
	if err != nil {
		return err
	}

	return nil
}

func showProgress(resp *http.Response, totalSize int, out *os.File, url string) error {
	var downloaded int
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			downloaded += n
			out.Write(buffer[:n])
			percentComplete := float64(downloaded) / float64(totalSize) * 100
			fmt.Printf("\rDownloading %s: %.2f%% complete", url, percentComplete)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFileNameFromURL(fileURL string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}
	return path.Base(parsedURL.Path), nil
}
