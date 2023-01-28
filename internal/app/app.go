package app

import "github.com/maraero/image-previewer/internal/imagesrv"

type App struct {
	ImageSrv imagesrv.ImageService
}

func New(imageSrv *imagesrv.ImageSrv) *App {
	return &App{imageSrv}
}
