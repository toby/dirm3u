package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Server struct {
	port     int
	hostname string
	files    map[string]bool
}

func NewServer(p int, h string) Server {
	return Server{
		port:     p,
		hostname: h,
		files:    make(map[string]bool),
	}
}

func (me *Server) m3uHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Content-Type", "application/mpegurl")
	for f, _ := range me.files {
		fmt.Fprintf(w, "http://%s:%d/media/%s\n", me.hostname, me.port, f)
	}
}

func (me *Server) mediaHandler(w http.ResponseWriter, r *http.Request) {
	f := path.Base(r.URL.Path)
	_, ok := me.files[f]
	if ok {
		http.ServeFile(w, r, f)
	} else {
		http.NotFound(w, r)
	}
}

func (me *Server) Start() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Adding file: %s\n", path)
		me.files[path] = true
		return nil
	})
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/playlist.m3u8", me.m3uHandler)
	http.HandleFunc("/media/", me.mediaHandler)
	fmt.Printf("Playlist: http://%s:%d/playlist.m3u8\n", me.hostname, me.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", me.port), nil))
}

func main() {
	p := flag.Int("p", 20202, "http server port")
	h := flag.String("h", "localhost", "hostname")
	flag.Parse()
	s := NewServer(*p, *h)
	s.Start()
}
