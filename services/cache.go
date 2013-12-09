package services

import (
	"apertoire.net/mediabase/bus"
	"apertoire.net/mediabase/helper"
	"apertoire.net/mediabase/message"
	"container/heap"
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

	reply chan int
	work  chan Request
	done  chan *Worker
	pool  Pool
}

type Request struct {
	media    *message.Media
	basePath string
	c        chan int
}

type Worker struct {
	requests chan Request
	pending  int
	index    int
}

func (w *Worker) work(done chan *Worker) {
	for {
		req := <-w.requests
		req.c <- cacheMedia(&req)
		done <- w
	}
}

type Pool []*Worker

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int {
	return len(p)
}

func (p *Pool) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	item := x.(*Worker)
	item.index = n
	a[n] = item
	*p = a
}

func (p *Pool) Pop() interface{} {
	a := *p
	fmt.Printf("Pop item %d\n", len(a)-1)
	n := len(a)
	item := a[n-1]
	item.index = -1
	*p = a[0 : n-1]
	return item
}

func (p Pool) Swap(i, j int) {
	// fmt.Printf("Swap(%d, %d) and pool length is %d\n", i, j, len(p))
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
	// fmt.Printf("Swap(%d, %d) and pool length is %d\n", i, j, len(p))
}

func (self *Cache) Start() {
	log.Println("starting cache service ...")

	self.reply = make(chan int)
	self.work = make(chan Request)
	self.done = make(chan *Worker)

	self.pool = make(Pool, 10)
	for k := 0; k < 10; k++ {
		self.pool[k] = &Worker{make(chan Request, 50), 0, 0}
		go self.pool[k].work(self.done)
	}

	go self.react()

	go self.balance()

	log.Println("cache service started")
}

func (self *Cache) Stop() {
	log.Printf("cache service stopped")
}

func (self *Cache) react() {
	for {
		select {
		case msg := <-self.Bus.CacheMedia:
			go self.requestor(msg)
		}
	}
}

func (self *Cache) requestor(media *message.Media) {
	self.work <- Request{media, self.Config.AppDir, self.reply}
	result := <-self.reply
	log.Printf("[code %d] work completed for %s", result, media.Movie.Title)
}

func (self *Cache) balance() {
	for {
		select {
		case req := <-self.work:
			self.dispatch(req)
		case w := <-self.done:
			self.completed(w)
		}
	}
}

func (self *Cache) dispatch(req Request) {
	fmt.Printf("dispatch request->%v\n", req)
	// Grab the least loaded worker...
	w := heap.Pop(&self.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in its work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&self.pool, w)
	fmt.Printf("end dispatch request->%v\n", req)
}

//Job is complete; update heap
func (self *Cache) completed(w *Worker) {
	fmt.Printf("completed worker->%s\n", w)
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	fmt.Printf("pool length is %d\n", len(self.pool))
	heap.Remove(&self.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&self.pool, w)
	fmt.Printf("done completed worker->%s\n", w)
}

func cacheMedia(req *Request) int {
	coverPath := filepath.Join(req.basePath, "/web/img/p", req.media.Movie.Cover)
	if _, err := os.Stat(coverPath); err == nil {
		// log.Printf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		// self.Bus.Log <- fmt.Sprintf("SKIP: picture in cache for [%s]: %s", picture.Name, picture.Id)
		tracelog.INFO("mb", "cache", fmt.Sprintf("SKIP: cover in cache for [%s]: %s", req.media.Movie.Title, req.media.Movie.Cover))
		return 1
	}

	// ext := filepath.Ext(picture.Path)
	// name := picture.Path[0:len(picture.Path)-len(ext)] + ".jpg"

	tracelog.INFO("mb", "cache", fmt.Sprintf("secureUrl is %s", req.media.SecureBaseUrl))

	helper.Download(req.media.SecureBaseUrl+"original"+req.media.Movie.Cover, coverPath)
	tracelog.INFO("mb", "cache", fmt.Sprintf("downloaded cover for %s", req.media.Movie.Title))
	// if err != nil {
	// 	// log.Printf("ERR: couldn't copy %s", name)
	// 	// self.Bus.Log <- fmt.Sprintf("couldn't copy %s", name)
	// 	tracelog.INFO("mb", "cache", fmt.Sprintf("for %s couldn't copy %s", picture.Name, name))

	// 	return
	// }

	// log.Printf("INFO: for [%s] copied image to %s", picture.Name, picture.Id)
	// self.Bus.Log <- fmt.Sprintf("INFO: for [%s] copied %s to %s", picture.Name, name, picture.Id)
	// tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] copied image to %s", picture.Name, picture.Id))

	doResize(coverPath, filepath.Join(req.basePath, "/web/img/t", req.media.Movie.Cover))
	tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] created thumb %s", req.media.Movie.Title, req.media.Movie.Cover))

	backdropPath := filepath.Join(req.basePath, "/web/img/b", req.media.Movie.Backdrop)
	helper.Download(req.media.SecureBaseUrl+"original"+req.media.Movie.Backdrop, backdropPath)
	tracelog.INFO("mb", "cache", fmt.Sprintf("INFO: for [%s] created backdrop %s", req.media.Movie.Title, req.media.Movie.Backdrop))

	return 0
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
