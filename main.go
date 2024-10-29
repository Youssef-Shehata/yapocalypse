package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync/atomic"
)

type healthHandler struct{}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "server is all good\n")

}

type apiConfig struct {
	homeHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.homeHits.Add(1)
		log.Println("home been hit")
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.homeHits.Store(0)
	fmt.Fprintf(w, "hits been reset to : 0")

}

type page struct {
	Hits int32
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page{Hits: cfg.homeHits.Load()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type server struct {
    Opts
}
type optFunc func(* Opts)
type Opts struct {
	conns int
	id    string
	tls   bool
}
func defaultOpts() Opts{
    return Opts{
    	conns: 3,
    	id:    "",
    	tls:   false,
    }
}
func maxConns(n int) optFunc {
    return func(opts *Opts){
        opts.conns=n
    }
}

func withTls(opts * Opts) {
    opts.tls = true
}
func newServer(opts ...optFunc) *server{

    o := defaultOpts()
    for _, fn := range opts{
        fn(&o)
    }

    return &server{o}
}

func main() {

	mux := http.NewServeMux()

	cfg := &apiConfig{homeHits: atomic.Int32{}}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hitting /\n")

	})
	mux.Handle("/app", cfg.middlewareMetrics(handler))
	mux.Handle("GET /admin/health", healthHandler{})
	mux.HandleFunc("GET /admin/metrics", cfg.metrics)
	mux.HandleFunc("POST /admin/reset", cfg.reset)

	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path != "/" {
	//		w.WriteHeader(404)
	//		fmt.Fprintf(w, "404 Page Not Found\n")
	//		return
	//	}
	//	fmt.Fprintf(w, "welcome to the home page\n")

	//})

	server := http.Server{Handler: mux, Addr: "localhost:8080"}
	err := server.ListenAndServe()

	if err != nil {
		log.Print("error listeninig to request", err)

	}

}
