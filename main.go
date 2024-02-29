package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"forum/forum/business"
	"forum/forum/handlers"
	"forum/forum/middleware"
	"forum/forum/repo"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.Parse()
	lg := LoggingMiddleware(*log.Default())
	rep, err := repo.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	bus, err := business.NewBusiness(rep)
	if err != nil {
		log.Fatal(err)
	}
	hand, err := handlers.NewHandler(bus)
	rateLimiter := middleware.NewRateLimiter(2)
	mux := http.NewServeMux()
	customHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow() {
			http.Error(w, "Too Many Requests. ", http.StatusTooManyRequests)
			return
		}

		// Your existing handler logic goes here.
		hand.ServeHTTP(w, r)
	})
	fileServer := http.FileServer(http.Dir("./forum/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", customHandler)
	addr := fmt.Sprintf(":%d", port)

	fmt.Printf("Server listening on http://localhost%s\n", addr)
	err = http.ListenAndServe(addr, lg(mux))
	if err != nil {
		fmt.Printf("Cannot open the port %s\n", addr)
		return
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{w, http.StatusOK}

			next.ServeHTTP(rw, r)

			defer func() {
				logger.Println(r.URL.RequestURI(), r.Method, rw.statusCode)
			}()
		}

		return http.HandlerFunc(fn)
	}
}
