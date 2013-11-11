package server

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

const apiVersion string = "/v1"

func static(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "public"+req.URL.Path)
}

func notFound(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "public/404.html")
}

func status(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func getEvents(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Nothing to see")
}

func Start() {
	log.Printf("starting server service")

	r := mux.NewRouter()
	s := r.PathPrefix(apiVersion).Subrouter()

	s.HandleFunc("/", status).Methods("GET")
	s.HandleFunc("/events", getEvents).Methods("GET")

	log.Printf("start listening on localhost:8080")
	http.ListenAndServe(":8080", s)
}
