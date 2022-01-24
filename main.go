package main

import (
	"log"
	"net/http"
)

func main() {
	var mws []Middleware
	m1 := mid1{}
	m2 := mid2{}
	m3 := mid3{}

	mws = append(mws,
		m1.HandleMiddleware,
		m2.HandleMiddleware,
		m3.HandleMiddleware,
	)

	mw := chainMiddleware(mws...)
	srv := server{}
	s := &http.Server{
		Addr:    ":8085",
		Handler: mw(srv.server),
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type Middleware func(handlerFunc HandlerFunc) HandlerFunc

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		log.Println(err)
	}
}

func chainMiddleware(mw ...Middleware) Middleware {
	return func(final HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) error {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			return last(w, r)
		}
	}
}

type server struct {
}

func (m *server) server(w http.ResponseWriter, r *http.Request) error {
	log.Println("server")
	_, _ = w.Write([]byte("server"))
	return nil
}

type mid1 struct{}

func (m *mid1) HandleMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("middleware 1")
		next.ServeHTTP(w, r)
		log.Println("middleware 1 back")
		return nil
	}
}

type mid2 struct{}

func (m *mid2) HandleMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("middleware 2")
		next.ServeHTTP(w, r)
		log.Println("middleware 2 back")
		return nil
	}
}

type mid3 struct{}

func (m *mid3) HandleMiddleware(next HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		log.Println("middleware 3")
		next.ServeHTTP(w, r)
		log.Println("middleware 3 back")
		return nil
	}
}
