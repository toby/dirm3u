package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//go:generate go-bindata -nomemcopy index.tmpl

type File struct {
	Info os.FileInfo
	Path string
}

type Files []*File

type FileDB struct {
	files Files
	tags  map[string]Files
}

type Server struct {
	Hostname      string
	Port          int
	Limit         int
	Pages         []Files
	indexTemplate *template.Template
	db            *FileDB
}

type Page struct {
	Server *Server
	Index  int
}

func NewFileDB() FileDB {
	return FileDB{
		tags: make(map[string]Files),
	}
}

func (db *FileDB) loadFiles() {
	db.files = make([]*File, 0)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		ts, err := FileTags(path)
		if err != nil {
			fmt.Printf("Skipping file: %s\n", path)
		} else {
			fmt.Printf("Adding file: %s\n", path)
			f := &File{info, path}
			db.files = append(db.files, f)
			for _, t := range ts {
				db.tags[t] = append(db.tags[t], f)
			}
		}
		return nil
	})
	for k, _ := range db.tags {
		fs, _ := db.tags[k]
		sort.Sort(fs)
		db.tags[k] = fs
		fmt.Printf("File Tag: '%s' has %d files\n", k, len(db.tags[k]))
	}
	if err != nil {
		panic(err)
	}
}

func (db *FileDB) Files() Files {
	return db.files
}

func (db *FileDB) ForTag(t string) (Files, bool) {
	fs, err := db.tags[t]
	return fs, err
}

func NewServer(p int, h string, l int) Server {
	data, err := Asset("index.tmpl")
	if err != nil {
		panic(err)
	}
	fm := template.FuncMap{
		"inc":   func(i int) int { return i + 1 },
		"image": IsImage,
	}
	tmpl, err := template.New("index.tmpl").Funcs(fm).Parse(string(data[:]))
	if err != nil {
		panic(err)
	}
	db := NewFileDB()
	s := Server{
		Port:          p,
		Hostname:      h,
		Limit:         l,
		indexTemplate: tmpl,
		db:            &db,
	}
	return s
}

func (me *Server) makePages() {
	var p Files
	me.db.loadFiles()
	ps := make([]Files, 0)
	fsv, _ := me.db.ForTag("web-video")
	fsi, _ := me.db.ForTag("web-image")
	fs := append(fsv, fsi...)
	sort.Sort(fs)
	if len(fs) == 0 {
		return
	}
	for n, f := range fs {
		if n%me.Limit == 0 {
			if n > 0 {
				ps = append(ps, p)
			}
			p = make([]*File, 0)
		}
		p = append(p, f)
	}
	if len(p) > 0 {
		ps = append(ps, p)
	}
	me.Pages = ps
}

func (me *Server) m3uHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/mpegurl")
	fs, ok := me.db.ForTag("vlc")
	if !ok {
		http.Error(w, "Page not found", 404)
		return
	}
	for _, f := range fs {
		fmt.Fprintf(w, "http://%s/media/%s\n", me.HostPort(), f.Path)
	}
}

func (me *Server) indexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := &Page{Server: me}
	n := ps.ByName("page")
	var i int
	var err error
	if n == "" {
		i = 1
	} else {
		i, err = strconv.Atoi(n)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if i > len(me.Pages) {
		http.Error(w, "Page not found", 404)
		return
	}
	p.Index = i - 1
	err = me.indexTemplate.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (me *Server) mediaHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	f := ps.ByName("path")
	if len(f) < 1 {
		http.NotFound(w, r)
	}
	f = f[1:]
	if me.db.Files().ContainsPath(f) {
		http.ServeFile(w, r, f)
	} else {
		http.NotFound(w, r)
	}
}

func (me *Server) reloadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	me.makePages()
	fmt.Fprintf(w, "Reloaded %d files\n", len(me.db.Files()))
	for _, f := range me.db.Files() {
		fmt.Fprintf(w, "http://%s/media/%s\n", me.HostPort(), url.PathEscape(f.Path))
	}
}

func (me *Server) Start() {
	me.makePages()
	router := httprouter.New()
	router.GET("/", me.indexHandler)
	router.GET("/page/:page", me.indexHandler)
	router.GET("/playlist.m3u", me.m3uHandler)
	router.GET("/reload", me.reloadHandler)
	router.GET("/media/*path", me.mediaHandler)
	fmt.Printf("Files: %d\n", len(me.db.Files()))
	fmt.Printf("Pages: %d\n", len(me.Pages))
	fmt.Printf("Serving: http://%s\n", me.HostPort())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", me.Port), router))
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
	l := flag.Int("l", 5, "limit results per page in web view")
	flag.Parse()
	s := NewServer(*p, *h, *l)
	s.Start()
}
