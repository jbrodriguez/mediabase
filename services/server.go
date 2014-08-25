package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

const apiVersion string = "/api/v1"
const docPath string = "build/"

type Server struct {
	Bus    *bus.Bus
	Config *helper.Config
	r, s   *mux.Router
}

func (self *Server) static(res http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	http.ServeFile(res, req, docPath+req.URL.Path)
}

func (self *Server) notFound(res http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	http.ServeFile(res, req, docPath+"404.html")
}

func (self *Server) status(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}

func (self *Server) scanMovies(w http.ResponseWriter, req *http.Request) {
	log.Println("you know .. i got here")
	// data := struct {
	// 	Code        int8
	// 	Description string
	// }{0, "all is good"}
	// helper.WriteJson(w, 200, &data)

	msg := message.ScanMovies{make(chan string)}
	self.Bus.ScanMovies <- &msg
	reply := <-msg.Reply

	log.Printf("response is: %s", reply)

	helper.WriteJson(w, 200, &helper.StringMap{"message": reply})
}

func (self *Server) pruneMovies(w http.ResponseWriter, req *http.Request) {
	log.Println("pruning .. i got here")
	// data := struct {
	// 	Code        int8
	// 	Description string
	// }{0, "all is good"}
	// helper.WriteJson(w, 200, &data)

	msg := message.PruneMovies{make(chan string)}
	self.Bus.PruneMovies <- &msg
	reply := <-msg.Reply

	log.Printf("response is: %s", reply)

	helper.WriteJson(w, 200, &helper.StringMap{"message": reply})
}

func (self *Server) getMovies(w http.ResponseWriter, req *http.Request) {
	msg := message.GetMovies{make(chan []*message.Movie)}
	self.Bus.GetMovies <- &msg
	reply := <-msg.Reply

	// log.Printf("response is: %s", reply)

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) listMovies(w http.ResponseWriter, req *http.Request) {
	msg := message.ListMovies{make(chan []*message.Movie)}
	self.Bus.ListMovies <- &msg
	reply := <-msg.Reply

	// log.Printf("response is: %s", reply)

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) showDuplicates(w http.ResponseWriter, req *http.Request) {
	msg := message.Movies{make(chan []*message.Movie)}
	self.Bus.ShowDuplicates <- &msg
	reply := <-msg.Reply
	log.Printf("never returned")

	// log.Printf("response is: %s", reply)

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) listByRuntime(w http.ResponseWriter, req *http.Request) {
	msg := message.Movies{make(chan []*message.Movie)}
	self.Bus.ListByRuntime <- &msg
	reply := <-msg.Reply

	// log.Printf("response is: %s", reply)

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) searchMovies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	term, ok := vars["term"]
	if !ok {
		// do some error handling
		return
	}

	log.Printf("the mother is: %s", term)

	msg := message.SearchMovies{term, make(chan []*message.Movie)}
	self.Bus.SearchMovies <- &msg
	reply := <-msg.Reply

	// log.Printf("response is: %s", reply)

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) fixMovies(w http.ResponseWriter, req *http.Request) {
	self.Bus.FixMovies <- 1
	helper.WriteJson(w, 200, "ok")
}

func (self *Server) testScan() {
	msg := message.ScanMovies{make(chan string)}
	self.Bus.ScanMovies <- &msg
	// reply := <-msg.Reply
}

// func (self *Server) postLogin(w http.ResponseWriter, req *http.Request) {
// 	log.Println("life's rich")
// 	user := &model.UserAuthReq{}
// 	if !helper.ReadJson(w, req, user) {
// 		data := struct {
// 			Code        int8
// 			Description string
// 		}{0, "not authorized"}
// 		helper.WriteJson(w, 304, &data)
// 		return
// 	}

// 	log.Printf("email: %s", user.Email)
// 	log.Printf("password: %s", user.Password)

// 	if user.Email == "" || user.Password == "" {
// 		helper.WriteJson(w, 400, &helper.StringMap{"error": "Invalid body"})
// 		return
// 	}

// 	msg := message.UserAuth{user, make(chan *model.UserAuthRep)}
// 	self.Bus.UserAuth <- &msg
// 	reply := <-msg.Reply

// 	helper.WriteJson(w, 200, &reply)
// }

// func (self *Server) getEvents(w http.ResponseWriter, req *http.Request) {
// 	io.WriteString(w, "Nothing to see")
// }

func (self *Server) Start() {
	log.Printf("starting server service")

	self.r = mux.NewRouter()

	self.r.PathPrefix("/" + docPath).Handler(http.StripPrefix("/"+docPath, http.FileServer(http.Dir("./"+docPath))))

	self.s = self.r.PathPrefix(apiVersion).Subrouter()
	self.s.HandleFunc("/", self.status).Methods("GET")
	self.s.HandleFunc("/movies", self.getMovies).Methods("GET")
	self.s.HandleFunc("/movies/all", self.listMovies).Methods("GET")
	self.s.HandleFunc("/movies/runtime", self.listByRuntime).Methods("GET")
	self.s.HandleFunc("/movies/duplicates", self.showDuplicates).Methods("GET")
	self.s.HandleFunc("/movies/scan", self.scanMovies).Methods("GET")
	self.s.HandleFunc("/movies/prune", self.pruneMovies).Methods("GET")
	self.s.HandleFunc("/movies/search&q={term}", self.searchMovies).Methods("GET")
	self.s.HandleFunc("/movies/fix", self.fixMovies).Methods("GET")

	// self.s.HandleFunc("/login", self.postLogin).Methods("POST")
	// self.s.HandleFunc("/events", self.getEvents).Methods("GET")

	self.r.Handle("/", http.RedirectHandler(docPath+"index.html", 302))

	// log.Printf("start listening on %s:%s", self.Config.Host, self.Config.Port)
	// go http.ListenAndServe(fmt.Sprintf("%s:%s", self.Config.Host, self.Config.Port), self.r)
	log.Printf("start listening on :%s", self.Config.Port)
	go http.ListenAndServe(fmt.Sprintf(":%s", self.Config.Port), self.r)

	// go self.testScan()
}

func (self *Server) Stop() {
	log.Printf("server service stopped")
	// nothing here
}
