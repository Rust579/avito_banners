package server

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

type Server struct {
	httpServer *fasthttp.Server
	addr       string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Run(handler fasthttp.RequestHandler) error {
	s.httpServer = &fasthttp.Server{
		Handler:          handler,
		ReadTimeout:      10 * time.Second,
		WriteTimeout:     10 * time.Second,
		DisableKeepalive: true, // set "connection: close" header in each response
	}

	return s.httpServer.ListenAndServe(s.addr)
}

func (s *Server) Shutdown() error {
	return s.httpServer.Shutdown()
}

// TODO swaggo
func ServDocs(host, port string) error {
	router := mux.NewRouter()
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)
	return http.ListenAndServe(host+":"+port, router)
}
