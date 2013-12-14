package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"fmt"
	"github.com/goinggo/tracelog"
	"github.com/goinggo/workpool"
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

type Cache struct {
	Bus      *bus.Bus
	Config   *helper.Config
	workpool *workpool.WorkPool
}

func (self *Cache) Start() {
	log.Println("starting cache service ...")

	self.workpool = workpool.New(4, 2000)

	go self.react()

	log.Println("cache service started")
}

func (self *Cache) Stop() {
	self.workpool.Shutdown("cache")
	log.Printf("cache service stopped")
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
	tracelog.INFO("mb", "cache", "CACHE MEDIA REQUESTED [%s]", media.Movie.Title)

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
	if _, err := os.Stat(coverPath); err == nil {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		tracelog.TRACE("mb", "cache", fmt.Sprintf("COVER DOWNLOAD SKIPPED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover))
	} else {
		helper.Download(self.media.SecureBaseUrl+"original"+self.media.Movie.Cover, coverPath)
		tracelog.TRACE("mb", "cache", fmt.Sprintf("COVER DOWNLOADED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover))
	}

	// ext := filepath.Ext(picture.Path)
	// name := picture.Path[0:len(picture.Path)-len(ext)] + ".jpg"

	// tracelog.TRACE("mb", "cache", fmt.Sprintf("secureUrl is %s", self.media.SecureBaseUrl))

	// if err != nil {
	// 	// log.Printf("ERR: couldn't copy %s", name)
	// 	// self.Bus.Log <- fmt.Sprintf("couldn't copy %s", name)
	// 	tracelog.INFO("mb", "cache", fmt.Sprintf("for %s couldn't copy %s", picture.Name, name))

	// 	return
	// }

	// log.Printf("INFO: for [%s] copied image to %s", picture.Name, picture.Id)
	// self.Bus.Log <- fmt.Sprintf("INFO: for [%s] copied %s to %s", picture.Name, name, picture.Id)
	// tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] copied image to %s", picture.Name, picture.Id))

	thumbPath := filepath.Join(self.appDir, "/web/img/t", self.media.Movie.Cover)
	if _, err := os.Stat(thumbPath); err == nil {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		tracelog.TRACE("mb", "cache", fmt.Sprintf("THUMB GENERATION SKIPPED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover))
	} else {
		doResize(coverPath, thumbPath)
		tracelog.TRACE("mb", "cache", fmt.Sprintf("THUMB CREATED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Cover))
	}

	backdropPath := filepath.Join(self.appDir, "/web/img/b", self.media.Movie.Backdrop)
	if _, err := os.Stat(backdropPath); err == nil {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		tracelog.TRACE("mb", "cache", fmt.Sprintf("BACKDROP DOWNLOAD SKIPPED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Backdrop))
	} else {
		helper.Download(self.media.SecureBaseUrl+"original"+self.media.Movie.Backdrop, backdropPath)
		tracelog.TRACE("mb", "cache", fmt.Sprintf("BACKDROP DOWNLOADED [%s] (%s)", self.media.Movie.Title, self.media.Movie.Backdrop))
	}

	// coverPath := filepath.Join(self.Config.AppDir, "/web/img/p", media.Movie.Cover)
	// if _, err := os.Stat(coverPath); err == nil {
	// 	// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
	// 	// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
	// 	tracelog.INFO("mb", "cache", fmt.Sprintf("SKIP: cover in cache for [%s]: %s", media.Movie.Title, media.Movie.Cover))

	// 	return
	// }

	// // ext := filepath.Ext(picture.Path)
	// // name := picture.Path[0:len(picture.Path)-len(ext)] + ".jpg"

	// helper.Download(media.SecureBaseUrl+"original"+media.Movie.Cover, coverPath)
	// // if err != nil {
	// // 	// log.Printf("ERR: couldn't copy %s", name)
	// // 	// self.Bus.Log <- fmt.Sprintf("couldn't copy %s", name)
	// // 	tracelog.INFO("mb", "cache", fmt.Sprintf("for %s couldn't copy %s", picture.Name, name))

	// // 	return
	// // }

	// // log.Printf("INFO: for [%s] copied image to %s", picture.Name, picture.Id)
	// // self.Bus.Log <- fmt.Sprintf("INFO: for [%s] copied %s to %s", picture.Name, name, picture.Id)
	// // tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] copied image to %s", picture.Name, picture.Id))

	// self.doResize(coverPath, filepath.Join(self.Config.AppDir, "/web/img/t", media.Movie.Cover))
	// // tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] created thumb to %s", filepath.Join(self.Config.AppDir, "/web/img/", "t_"+picture.Id)))

	// backdropPath := filepath.Join(self.Config.AppDir, "/web/img/b", media.Movie.Backdrop)
	// helper.Download(media.SecureBaseUrl+"original"+media.Movie.Backdrop, backdropPath)
}

func doResize(src, dst string) {
	// open "test.jpg"
	file, err := os.Open(src)
	if err != nil {
		log.Printf("[%s] unable to open %s", err, src)
		return
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Printf("[%s] unable to decode %s", err, src)
		return
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(80, 0, img, resize.Lanczos3)

	out, err := os.Create(dst)
	if err != nil {
		log.Printf("[%s] unable to create %s", err, dst)
		return
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
}
