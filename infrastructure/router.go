package infrastructure

import (
	"github.com/go-chi/chi/v5"
	controllers "main/interfaces/api"
)

func AppendDownloaderRoute(controller controllers.DownloaderController) func(router chi.Router) {
	return func(router chi.Router) {
		router.Post("/", controller.CreateDownloader)
	}
}
