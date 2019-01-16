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

type File struct {
	info os.FileInfo
	path string
}

type Files []File

type Server struct {
	port     int
	hostname string
	files    Files
}

func NewServer(p int, h string) Server {
	return Server{
		port:     p,
		hostname: h,
		files:    make([]File, 0),
	}
}

func (me *Server) loadFiles() {
	me.files = make([]File, 0)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if SupportedType(path) {
			fmt.Printf("Adding file: %s\n", path)
			me.files = append(me.files, File{info, path})
		} else {
			fmt.Printf("Skipping file: %s\n", path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (me *Server) m3uHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/mpegurl")
	for _, f := range me.files {
		fmt.Fprintf(w, "http://%s:%d/media/%s\n", me.hostname, me.port, f)
	}
}

func (me *Server) mediaHandler(w http.ResponseWriter, r *http.Request) {
	f := path.Base(r.URL.Path)
	if me.files.ContainsPath(f) {
		http.ServeFile(w, r, f)
	} else {
		http.NotFound(w, r)
	}
}

func (me *Server) reloadHandler(w http.ResponseWriter, r *http.Request) {
	me.loadFiles()
	fmt.Fprintf(w, "Reloaded %d files\n", len(me.files))
	for _, f := range me.files {
		fmt.Fprintf(w, "http://%s:%d/media/%s\n", me.hostname, me.port, f.path)
	}
}

func (me *Server) Start() {
	me.loadFiles()
	http.HandleFunc("/playlist.m3u", me.m3uHandler)
	http.HandleFunc("/reload", me.reloadHandler)
	http.HandleFunc("/media/", me.mediaHandler)
	fmt.Printf("Playlist: http://%s:%d/playlist.m3u\n", me.hostname, me.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", me.port), nil))
}

func (fs Files) ContainsPath(p string) bool {
	for _, f := range fs {
		if f.path == p {
			return true
		}
	}
	return false
}

func main() {
	p := flag.Int("p", 20202, "http server port")
	h := flag.String("h", "localhost", "hostname")
	flag.Parse()
	s := NewServer(*p, *h)
	s.Start()
}
