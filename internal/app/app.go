package app

import "github.com/maraero/image-previewer/internal/resizesrv"

type App struct {
	ResizeSrv *resizesrv.ResizeSrv
}

func New(resizeSrv *resizesrv.ResizeSrv) *App {
	return &App{resizeSrv}
}
