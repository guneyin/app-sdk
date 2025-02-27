package router

type HTTPMethod string

const (
	GET  HTTPMethod = "GET"
	POST HTTPMethod = "POST"
)

type HTTPHandler interface {
	Register(r *Server)
}

type HandlerFunc func(*Ctx) error

type Map map[string]interface{}
