package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port     int
	hostname string
}

func (me *Server) m3uHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "TODO m3u")
}

func (me *Server) Start() {
	http.HandleFunc("/playlist.m3u", me.m3uHandler)
	fmt.Printf("Playlist: http://%s:%d/playlist.m3u\n", me.hostname, me.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", me.port), nil))
}

func main() {
	p := flag.Int("p", 20202, "http server port")
	h := flag.String("h", "localhost", "hostname")
	flag.Parse()
	s := Server{
		port:     *p,
		hostname: *h,
	}
	s.Start()
}
