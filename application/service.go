package application

import "github.com/guneyin/app-sdk/router"

type Service interface {
	RegisterHandlers(r *router.Server)
}

func (app *App) RegisterService(s Service) {
	s.RegisterHandlers(app.router)
}
