package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

var webVideoExtensions = map[string]bool{
	".webm": true,
	".mp4":  true,
	".ogg":  true,
}

var webImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".png":  true,
}

var vlcExtensions = map[string]bool{
	".asx":   true,
	".dts":   true,
	".gxf":   true,
	".m2v":   true,
	".m3u":   true,
	".m4v":   true,
	".mpeg1": true,
	".mpeg2": true,
	".mts":   true,
	".mxf":   true,
	".ogm":   true,
	".pls":   true,
	".bup":   true,
	".a52":   true,
	".aac":   true,
	".b4s":   true,
	".cue":   true,
	".divx":  true,
	".dv":    true,
	".flv":   true,
	".m1v":   true,
	".m2ts":  true,
	".mkv":   true,
	".mov":   true,
	".mpeg4": true,
	".oma":   true,
	".spx":   true,
	".ts":    true,
	".vlc":   true,
	".vob":   true,
	".xspf":  true,
	".dat":   true,
	".bin":   true,
	".ifo":   true,
	".part":  true,
	".3g2":   true,
	".avi":   true,
	".mpeg":  true,
	".mpg":   true,
	".flac":  true,
	".m4a":   true,
	".mp1":   true,
	".ogg":   true,
	".wav":   true,
	".xm":    true,
	".3gp":   true,
	// ".srt":   true,
	".wmv":  true,
	".ac3":  true,
	".asf":  true,
	".mod":  true,
	".mp2":  true,
	".mp3":  true,
	".mp4":  true,
	".wma":  true,
	".mka":  true,
	".m4p":  true,
	".webm": true,
}

func Extension(p string) string {
	return strings.ToLower(filepath.Ext(p))
}

func IsImage(p string) bool {
	_, ok := webImageExtensions[Extension(p)]
	return ok
}

func FileTags(p string) ([]string, error) {
	ts := make([]string, 0)
	ext := Extension(p)
	_, ok := webVideoExtensions[ext]
	if ok {
		ts = append(ts, "web-video")
	}
	_, ok = webImageExtensions[ext]
	if ok {
		ts = append(ts, "web-image")
	}
	_, ok = vlcExtensions[ext]
	if ok {
		ts = append(ts, "vlc")
	}
	if len(ts) == 0 {
		return ts, fmt.Errorf("file type %s not supported", ext)
	}
	return ts, nil
}
