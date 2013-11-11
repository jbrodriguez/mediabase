package services

import (
	"apertoire.net/moviebase/bus"
	"apertoire.net/moviebase/helper"
	"apertoire.net/moviebase/message"
	"apertoire.net/moviebase/model"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

const apiVersion string = "/v1"

type Server struct {
	Bus    *bus.Bus
	Config helper.Config
	r, s   *mux.Router
}

func (self *Server) static(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "public"+req.URL.Path)
}

func (self *Server) notFound(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "public/404.html")
}

func (self *Server) status(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func (self *Server) postLogin(w http.ResponseWriter, req *http.Request) {
	log.Println("life's rich")
	user := &model.UserAuthReq{}
	if !helper.ReadJson(w, req, user) {
		data := struct {
			Code        int8
			Description string
		}{0, "not authorized"}
		helper.WriteJson(w, 304, &data)
		return
	}

	log.Printf("email: %s", user.Email)
	log.Printf("password: %s", user.Password)

	if user.Email == "" || user.Password == "" {
		helper.WriteJson(w, 400, &helper.StringMap{"error": "Invalid body"})
		return
	}

	msg := message.UserAuth{user, make(chan *model.UserAuthRep)}
	self.Bus.UserAuth <- &msg
	reply := <-msg.Reply

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) getEvents(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Nothing to see")
}

func (self *Server) Start() {
	log.Printf("starting server service")

	self.r = mux.NewRouter()
	self.s = self.r.PathPrefix(apiVersion).Subrouter()

	self.s.HandleFunc("/", self.status).Methods("GET")
	self.s.HandleFunc("/login", self.postLogin).Methods("POST")
	self.s.HandleFunc("/events", self.getEvents).Methods("GET")

	log.Printf("start listening on %s:%s", self.Config.Host, self.Config.Port)
	go http.ListenAndServe(fmt.Sprintf("%s:%s", self.Config.Host, self.Config.Port), self.s)
}

func (self *Server) Stop() {
	// nothing here
}
