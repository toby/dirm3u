package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

//go:generate go-bindata -nomemcopy index.tmpl

type File struct {
	Info os.FileInfo
	Path string
}

type Files []File

type Server struct {
	Port          int
	Hostname      string
	Files         Files
	indexTemplate *template.Template
}

func NewServer(p int, h string) Server {
	data, err := Asset("index.tmpl")
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("index.tmpl").Parse(string(data[:]))
	if err != nil {
		panic(err)
	}
	s := Server{
		Port:          p,
		Hostname:      h,
		Files:         make([]File, 0),
		indexTemplate: tmpl,
	}
	return s
}

func (me *Server) loadFiles() {
	me.Files = make([]File, 0)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if SupportedType("web", path) {
			fmt.Printf("Adding file: %s\n", path)
			me.Files = append(me.Files, File{info, path})
		} else {
			fmt.Printf("Skipping file: %s\n", path)
		}
		return nil
	})
	sort.Sort(me.Files)
	if err != nil {
		panic(err)
	}
}

func (me *Server) m3uHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/mpegurl")
	for _, f := range me.Files {
		fmt.Fprintf(w, "http://%s/media/%s\n", me.HostPort(), f.Path)
	}
}

func (me *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	err := me.indexTemplate.Execute(w, me)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (me *Server) mediaHandler(w http.ResponseWriter, r *http.Request) {
	// trim `/media/`
	f := r.URL.Path[7:]
	if me.Files.ContainsPath(f) {
		http.ServeFile(w, r, f)
	} else {
		http.NotFound(w, r)
	}
}

func (me *Server) reloadHandler(w http.ResponseWriter, r *http.Request) {
	me.loadFiles()
	fmt.Fprintf(w, "Reloaded %d files\n", len(me.Files))
	for _, f := range me.Files {
		fmt.Fprintf(w, "http://%s/media/%s\n", me.HostPort(), f.Path)
	}
}

func (me *Server) Start() {
	me.loadFiles()
	http.HandleFunc("/", me.indexHandler)
	http.HandleFunc("/playlist.m3u", me.m3uHandler)
	http.HandleFunc("/reload", me.reloadHandler)
	http.HandleFunc("/media/", me.mediaHandler)
	fmt.Printf("Serving: http://%s\n", me.HostPort())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", me.Port), nil))
}

func (me *Server) HostPort() string {
	return fmt.Sprintf("%s:%d", me.Hostname, me.Port)
}

func (f File) Type() string {
	return filepath.Ext(f.Path)[1:]
}

func (f File) Base() string {
	return filepath.Base(f.Path)
}

func (fs Files) ContainsPath(p string) bool {
	for _, f := range fs {
		if f.Path == p {
			return true
		}
	}
	return false
}

func (fs Files) Len() int {
	return len(fs)
}

func (fs Files) Less(i int, j int) bool {
	return fs[i].Info.ModTime().After(fs[j].Info.ModTime())
}

func (fs Files) Swap(i int, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

func main() {
	p := flag.Int("p", 20202, "http server port")
	h := flag.String("h", "localhost", "hostname")
	flag.Parse()
	s := NewServer(*p, *h)
	s.Start()
}
