package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"main/usecase"
	"net/http"
)

type DownloaderController interface {
	CreateDownloader(w http.ResponseWriter, r *http.Request)
}

type ControllerDownloader struct {
	service usecase.IDownloaderService
}

func (impl *ControllerDownloader) CreateDownloader(w http.ResponseWriter, r *http.Request) {
	var urls []string
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	if err := json.Unmarshal(bodyBytes, urls); err != nil {
		_ = fmt.Errorf("error to unmarshal body %d", err)
	}
	impl.service.ProcessDownload(urls, "/test")
}

func NewDownloaderController(service usecase.IDownloaderService) DownloaderController {
	return &ControllerDownloader{
		service: service,
	}
}
