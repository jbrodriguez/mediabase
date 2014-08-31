package services

import (
	"apertoire.net/mediabase/server/bus"
	"apertoire.net/mediabase/server/helper"
	"apertoire.net/mediabase/server/message"
	"github.com/apertoire/mlog"
	"github.com/goinggo/workpool"
	"github.com/nfnt/resize"
	"image/jpeg"
	"os"
	"path/filepath"
)

type Cache struct {
	Bus      *bus.Bus
	Config   *helper.Config
	workpool *workpool.WorkPool
}

func (self *Cache) Start() {
	mlog.Info("starting cache service ...")

	self.workpool = workpool.New(4, 2000)

	go self.react()

	mlog.Info("cache service started")
}

func (self *Cache) Stop() {
	self.workpool.Shutdown("cache")
	mlog.Info("cache service stopped")
}

func (self *Cache) react() {
	for {
		select {
		case msg := <-self.Bus.CacheMedia:
			go self.requestWork(msg)
		}
	}
}

func (self *Cache) requestWork(media *message.Media) {
	mlog.Info("CACHE MEDIA REQUESTED [%s]", media.Movie.Title)

	gig := &CacheGig{
		media,
		self.Config.AppDir,
	}

	self.workpool.PostWork("cachegig", gig)
}

type CacheGig struct {
	media  *message.Media
	appDir string
}

func (self *CacheGig) DoWork(workRoutine int) {
	coverPath := filepath.Join(self.appDir, "/web/img/p", self.media.Movie.Cover)
	if _, err := os.Stat(coverPath); err == nil && !self.media.Forced {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		mlog.Info("COVER DOWNLOAD SKIPPED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover)
	} else {
		helper.Download(self.media.SecureBaseUrl+"original"+self.media.Movie.Cover, coverPath)
		mlog.Info("COVER DOWNLOADED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover)
	}

	thumbPath := filepath.Join(self.appDir, "/web/img/t", self.media.Movie.Cover)
	if _, err := os.Stat(thumbPath); err == nil && !self.media.Forced {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		mlog.Info("THUMB GENERATION SKIPPED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover)
	} else {
		doResize(coverPath, thumbPath)
		mlog.Info("THUMB CREATED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover)
	}

	backdropPath := filepath.Join(self.appDir, "/web/img/b", self.media.Movie.Backdrop)
	if _, err := os.Stat(backdropPath); err == nil && !self.media.Forced {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		mlog.Info("BACKDROP DOWNLOAD SKIPPED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Backdrop)
	} else {
		helper.Download(self.media.SecureBaseUrl+"original"+self.media.Movie.Backdrop, backdropPath)
		mlog.Info("BACKDROP DOWNLOADED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Backdrop)
	}

}

func doResize(src, dst string) {
	// open "test.jpg"
	file, err := os.Open(src)
	if err != nil {
		mlog.Error(err)
		return
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		mlog.Error(err)
		return
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(80, 0, img, resize.Lanczos3)

	out, err := os.Create(dst)
	if err != nil {
		mlog.Error(err)
		return
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
