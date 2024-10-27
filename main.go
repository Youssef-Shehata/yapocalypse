package main

import (
	"fmt"
	"log"
	"net/http"
)


type healthHandler struct {}
func (h healthHandler) ServeHTTP( w http.ResponseWriter , r *http.Request){
    if r.URL.Path != "/health"{
        w.WriteHeader(404)
    }
        w.WriteHeader(200)
    w.Header().Add("Content-Type" , "text/plain")
    fmt.Fprintf(w , "server is all good\n")

}

func main() {
	mux := http.NewServeMux()

    mux.Handle("/app", http.StripPrefix("/app",http.FileServer(http.Dir("."))) )
    mux.Handle("/files", http.FileServer(http.Dir("."))) 
    mux.Handle("/health", healthHandler {})

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
