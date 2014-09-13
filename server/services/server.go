package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/helper"
	"apertoire.net/mediabase/server/message"
	"apertoire.net/mediabase/server/static"
	"fmt"
	"github.com/apertoire/mlog"
	"github.com/gin-gonic/gin"
	// "io"
	"net/http"
)

const apiVersion string = "/api/v1"
const docPath string = ""

type Server struct {
	Bus    *bus.Bus
	Config *helper.Config
	r, s   *gin.Engine
}

func (self *Server) Start() {
	mlog.Info("starting server service")

	self.r = gin.New()

	self.r.Use(gin.Recovery())
	self.r.Use(helper.Logging())

	self.r.Use(static.Serve("./"))
	self.r.NoRoute(static.Serve("./"))

	api := self.r.Group(apiVersion)
	{
		api.GET("/movies", self.getMovies)
		api.GET("/all", self.listMovies)
		api.GET("/import", self.importMovies)
		api.GET("/search/:term", self.searchMovies)

		api.POST("/movie/watched", self.watchedMovie)
		api.POST("/movie/fix", self.fixMovie)
	}

	mlog.Info("service started listening on %s:%s", self.Config.Host, self.Config.Port)

	go self.r.Run(fmt.Sprintf("%s:%s", self.Config.Host, self.Config.Port))
}

func (self *Server) Stop() {
	mlog.Info("server service stopped")
	// nothing here
}

func (self *Server) getMovies(c *gin.Context) {
	msg := message.GetMovies{Reply: make(chan []*message.Movie)}
	self.Bus.GetMovies <- &msg
	reply := <-msg.Reply

	// mlog.Info("response is: %s", reply)

	// helper.WriteJson(w, 200, &reply)
	c.JSON(200, &reply)
}

func (self *Server) listMovies(c *gin.Context) {
	msg := message.ListMovies{Reply: make(chan []*message.Movie)}
	self.Bus.ListMovies <- &msg
	reply := <-msg.Reply

	// mlog.Info("response is: %s", reply)

	c.JSON(200, &reply)
}

func (self *Server) importMovies(c *gin.Context) {
	mlog.Info("importMovies: you know .. i got here")

	msg := message.Status{Reply: make(chan *message.Context)}
	self.Bus.ImportMovies <- &msg
	reply := <-msg.Reply

	// msg := message.ScanMovies{Reply: make(chan string)}
	// self.Bus.ScanMovies <- &msg
	// reply := <-msg.Reply

	// mlog.Info("response is: %+v", reply)

	// helper.WriteJson(w, 200, &helper.StringMap{"message": reply})
	c.JSON(200, &reply)
}

func (self *Server) searchMovies(c *gin.Context) {
	mlog.Info("searchMovies: are you a head honcho ?")
	term := c.Params.ByName("term")

	msg := message.SearchMovies{Term: term, Reply: make(chan []*message.Movie)}
	self.Bus.SearchMovies <- &msg
	reply := <-msg.Reply

	// mlog.Info("%s", reply)

	c.JSON(200, &reply)
}

func (self *Server) pruneMovies(w http.ResponseWriter, req *http.Request) {
	mlog.Info("pruning .. i got here")
	// data := struct {
	// 	Code        int8
	// 	Description string
	// }{0, "all is good"}
	// helper.WriteJson(w, 200, &data)

	msg := message.PruneMovies{Reply: make(chan string)}
	self.Bus.PruneMovies <- &msg
	reply := <-msg.Reply

	mlog.Info("response is: %s", reply)

	helper.WriteJson(w, 200, &helper.StringMap{"message": reply})
}

func (self *Server) showDuplicates(w http.ResponseWriter, req *http.Request) {
	msg := message.Movies{Reply: make(chan []*message.Movie)}
	self.Bus.ShowDuplicates <- &msg
	reply := <-msg.Reply
	mlog.Info("never returned")

	// mlog.Info("response is: %s", reply)

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) listByRuntime(w http.ResponseWriter, req *http.Request) {
	msg := message.Movies{Reply: make(chan []*message.Movie)}
	self.Bus.ListByRuntime <- &msg
	reply := <-msg.Reply

	// mlog.Info("response is: %s", reply)

	helper.WriteJson(w, 200, &reply)
}

func (self *Server) watchedMovie(c *gin.Context) {
	var movie message.Movie

	c.Bind(&movie)
	// mlog.Info("%+v", movie)

	msg := message.SingleMovie{Movie: &movie, Reply: make(chan bool)}
	self.Bus.WatchedMovie <- &msg
	reply := <-msg.Reply

	data := struct {
		Status bool `json:"status"`
	}{Status: reply}

	c.JSON(200, &data)
}

func (self *Server) fixMovies(w http.ResponseWriter, req *http.Request) {
	self.Bus.FixMovies <- 1
	helper.WriteJson(w, 200, "ok")
}

func (self *Server) fixMovie(c *gin.Context) {
	var movie message.Movie

	c.Bind(&movie)
	mlog.Info("%+v", movie)

	msg := message.SingleMovie{Movie: &movie, Reply: make(chan bool)}
	self.Bus.FixMovie <- &msg

	data := struct {
		Status bool `json:"status"`
	}{Status: true}

	c.JSON(200, &data)
}
