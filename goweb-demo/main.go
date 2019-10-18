package main

import (
	"log"
	"net/http"
	"time"
)
type HandlerFunc func(*Context)
type MiddleWare func(HandlerFunc)HandlerFunc
type MyServer struct {
	router map[string]map[string]HandlerFunc
	chain []MiddleWare
}

type Context struct {
	Rw  http.ResponseWriter
	R *http.Request
}

func timeMiddleware(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx *Context) {
		start := time.Now()
		next(ctx)
		elapsed := time.Since(start)
		log.Println("time elapsed",elapsed)
	})
}

func (ctx *Context)SayHello()  {
	ctx.Rw.WriteHeader(200)
	ctx.Rw.Write([]byte("hello world"))
}

func (s *MyServer)ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if  _,found := s.router[r.Method] ;found {
		if fn,found :=s.router[r.Method][r.URL.Path];found {
			fn(&Context{Rw:rw,R:r})
			return
		}
	}
	rw.WriteHeader(404)
	rw.Write([]byte("page not found"))
}

func (s *MyServer)Use(middleware ...MiddleWare)  {
	for _,m:= range middleware {
		s.chain = append(s.chain,m)
	}
}


func (s *MyServer)Get(path string,fn HandlerFunc)  {
	if s.router["GET"] == nil {
		s.router["GET"] = make(map[string]HandlerFunc)
	}
	handler := fn
	for i := len(s.chain) - 1; i >= 0; i-- {
		handler = s.chain[i](handler)
	}
	s.router["GET"][path] = handler
}

func NewServer() *MyServer  {
	return &MyServer{
		router: make(map[string]map[string]HandlerFunc),
	}
}
func main()  {
	s := NewServer()
	s.Use(timeMiddleware,timeMiddleware)
	s.Get("/get", func(ctx *Context) {
		ctx.SayHello()
	})
	http.ListenAndServe(":8888", s)
}
