package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/message"
	"apertoire.net/mediabase/server/model"
	"fmt"
	"github.com/apertoire/mlog"
	"github.com/looplab/fsm"
)

type Core struct {
	Bus     *bus.Bus
	Config  *model.Config
	fsm     *fsm.FSM
	context message.Context
}

func (self *Core) Start() {
	mlog.Info("starting core service ...")

	// some initialization
	self.fsm = fsm.NewFSM(
		"idle",
		fsm.Events{
			{Name: "import", Src: []string{"idle", "scanning"}, Dst: "scanning"},
			{Name: "found", Src: []string{"scanning"}, Dst: "scanning"},
			{Name: "scraped", Src: []string{"scanning"}, Dst: "scanning"},
			{Name: "status", Src: []string{"idle", "scanning"}, Dst: "scanning"},
			{Name: "finish", Src: []string{"scanning"}, Dst: "idle"},
		},
		fsm.Callbacks{
			"import":  self.importer,
			"found":   self.found,
			"scraped": self.scraped,
			"finish":  self.finish,
		},
	)

	self.context = message.Context{Message: "Idle", Backdrop: "/mAwd34SAC8KqBKRm2MwHPLhLDU5.jpg", Completed: false}

	go self.react()

	mlog.Info("core service started")
}

func (self *Core) Stop() {
	// some deinitialization
	mlog.Info("core service stopped")
}

func (self *Core) react() {
	for {
		select {
		case msg := <-self.Bus.GetConfig:
			go self.doGetConfig(msg)
		case msg := <-self.Bus.SaveConfig:
			go self.doSaveConfig(msg)
		case msg := <-self.Bus.ImportMovies:
			go self.doImportMovies(msg)
		case msg := <-self.Bus.MovieFound:
			go self.doMovieFound(msg)
		case msg := <-self.Bus.MovieScraped:
			go self.doMovieScraped(msg)
		case msg := <-self.Bus.MovieRescraped:
			go self.doMovieRescraped(msg)
		// case msg := <-self.Bus.FixMovies:
		// 	go self.doFixMovies(msg)
		case msg := <-self.Bus.FixMovie:
			go self.doFixMovie(msg)

		case msg := <-self.Bus.ImportMoviesFinished:
			go self.doImportMoviesFinished(msg)
		}
	}
}

func (self *Core) doGetConfig(msg *message.GetConfig) {
	msg.Reply <- self.Config
}

func (self *Core) doSaveConfig(msg *message.SaveConfig) {
	self.Config = msg.Config
	self.Config.Save()

	msg.Reply <- true
}

func (self *Core) doScrape(e *fsm.Event) {
	movie, _ := e.Args[0].(*message.Movie)
	self.context.Message = fmt.Sprintf("Scraping %s ", movie.Location)
}

func (self *Core) doImportMovies(status *message.Status) {
	if err := self.fsm.Event("import", status); err != nil {
		mlog.Info("error trying to trigger import event: %s", err)
	}
}

func (self *Core) importer(e *fsm.Event) {
	if e.Src == "idle" {
		msg := message.ScanMovies{Reply: make(chan string)}
		self.Bus.ScanMovies <- &msg
		reply := <-msg.Reply

		self.context.Message = reply
	}

	mlog.Info("Before sending some answers %s: %s", e.Event, e.FSM.Current())

	status, _ := e.Args[0].(*message.Status)
	status.Reply <- &self.context

	mlog.Info("Event %s was fired, currently in state %s", e.Event, e.FSM.Current())
}

func (self *Core) doMovieFound(movie *message.Movie) {
	c := make(chan bool)

	self.Bus.CheckMovie <- &message.CheckMovie{Movie: movie, Result: c}
	exists := <-c

	var text string
	if exists {
		text = fmt.Sprintf("SKIPPED: present in db [%s] (%s)", movie.Title, movie.Location)
		mlog.Info(text)
	} else {
		text = fmt.Sprintf("FOUND: [%s] (%s)", movie.Title, movie.Location)
		self.Bus.ScrapeMovie <- movie
	}

	if err := self.fsm.Event("found", text); err != nil {
		mlog.Info("error trying to trigger found event: %s", err)
	}
}

func (self *Core) found(e *fsm.Event) {
	text, _ := e.Args[0].(string)
	self.context.Message = text
}

func (self *Core) doMovieScraped(media *message.Media) {
	go func() {
		mlog.Info("STORING MOVIE [%s]", media.Movie.Title)
		self.Bus.StoreMovie <- media.Movie
	}()

	go func() {
		mlog.Info("CACHING MEDIA [%s]", media.Movie.Title)
		media.BasePath = self.Config.AppDir
		self.Bus.CacheMedia <- media

		self.fsm.Event("scrape", media.Movie)
	}()
}

func (self *Core) scraped(e *fsm.Event) {
	movie, _ := e.Args[0].(*message.Movie)
	self.context.Message = fmt.Sprintf("Scraped %s ", movie.Location)
}

func (self *Core) doMovieRescraped(media *message.Media) {
	go func() {
		mlog.Info("UPDATING MOVIE [%s]", media.Movie.Title)
		self.Bus.UpdateMovie <- media.Movie
	}()

	go func() {
		mlog.Info("CACHING MEDIA [%s]", media.Movie.Title)
		media.BasePath = self.Config.AppDir
		self.Bus.CacheMedia <- media
	}()
}

// func (self *Core) doFixMovies(flag int) {
// 	msg := message.Movies{Reply: make(chan *message.MoviesDTO)}
// 	self.Bus.GetMoviesToFix <- &msg

// 	mlog.Info("AFTER GET MOVIES TO FIX [%v]", msg.Reply)

// 	reply := <-msg.Reply

// 	mlog.Info("WAITING FOR REPLY [%v]", reply)

// 	self.Bus.RescrapeMovies <- reply
// }

func (self *Core) doFixMovie(msg *message.SingleMovie) {
	movies := make([]*message.Movie, 0)
	movies = append(movies, msg.Movie)

	self.Bus.RescrapeMovies <- &message.MoviesDTO{Count: 1, Movies: movies}
}

func (self *Core) doImportMoviesFinished(completed int) {
	self.fsm.Event("finish")
}

func (self *Core) finish(e *fsm.Event) {
	self.context.Message = "Import completed"
	self.context.Completed = true
}
