package server

import (
	"github.com/valyala/fasthttp"
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
		DisableKeepalive: true,
	}

	return s.httpServer.ListenAndServe(s.addr)
}

func (s *Server) Shutdown() error {
	return s.httpServer.Shutdown()
}
