package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"fmt"
	"github.com/goinggo/tracelog"
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

type Cache struct {
	Bus    *bus.Bus
	Config *helper.Config
}

func (self *Cache) Start() {
	log.Println("starting cache service ...")

	go self.react()

	log.Println("cache service started")
}

func (self *Cache) Stop() {

}

func (self *Cache) react() {
	for {
		select {
		case msg := <-self.Bus.CacheMedia:
			self.doCacheMedia(msg)
		}
	}
}

func (self *Cache) doCacheMedia(media *message.Media) {
	coverPath := filepath.Join(self.Config.AppDir, "/web/img/p", media.Movie.Cover)
	if _, err := os.Stat(coverPath); err == nil {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		tracelog.INFO("mb", "cache", fmt.Sprintf("SKIP: cover in cache for [%s]: %s", media.Movie.Title, media.Movie.Cover))

		return
	}

	// ext := filepath.Ext(picture.Path)
	// name := picture.Path[0:len(picture.Path)-len(ext)] + ".jpg"

	helper.Download(media.SecureBaseUrl+"original"+media.Movie.Cover, coverPath)
	// if err != nil {
	// 	// log.Printf("ERR: couldn't copy %s", name)
	// 	// self.Bus.Log <- fmt.Sprintf("couldn't copy %s", name)
	// 	tracelog.INFO("mb", "cache", fmt.Sprintf("for %s couldn't copy %s", picture.Name, name))

	// 	return
	// }

	// log.Printf("INFO: for [%s] copied image to %s", picture.Name, picture.Id)
	// self.Bus.Log <- fmt.Sprintf("INFO: for [%s] copied %s to %s", picture.Name, name, picture.Id)
	// tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] copied image to %s", picture.Name, picture.Id))

	self.doResize(coverPath, filepath.Join(self.Config.AppDir, "/web/img/t", media.Movie.Cover))
	// tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] created thumb to %s", filepath.Join(self.Config.AppDir, "/web/img/", "t_"+picture.Id)))

	backdropPath := filepath.Join(self.Config.AppDir, "/web/img/b", media.Movie.Backdrop)
	helper.Download(media.SecureBaseUrl+"original"+media.Movie.Backdrop, backdropPath)
}

func (self *Cache) doResize(src, dst string) {
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
