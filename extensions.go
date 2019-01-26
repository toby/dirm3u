package main

import (
	"path/filepath"
	"strings"
)

var extensions = map[string]bool{
	".ASX":   true,
	".DTS":   true,
	".GXF":   true,
	".M2V":   true,
	".M3U":   true,
	".M4V":   true,
	".MPEG1": true,
	".MPEG2": true,
	".MTS":   true,
	".MXF":   true,
	".OGM":   true,
	".PLS":   true,
	".BUP":   true,
	".A52":   true,
	".AAC":   true,
	".B4S":   true,
	".CUE":   true,
	".DIVX":  true,
	".DV":    true,
	".FLV":   true,
	".M1V":   true,
	".M2TS":  true,
	".MKV":   true,
	".MOV":   true,
	".MPEG4": true,
	".OMA":   true,
	".SPX":   true,
	".TS":    true,
	".VLC":   true,
	".VOB":   true,
	".XSPF":  true,
	".DAT":   true,
	".BIN":   true,
	".IFO":   true,
	".PART":  true,
	".3G2":   true,
	".AVI":   true,
	".MPEG":  true,
	".MPG":   true,
	".FLAC":  true,
	".M4A":   true,
	".MP1":   true,
	".OGG":   true,
	".WAV":   true,
	".XM":    true,
	".3GP":   true,
	".SRT":   true,
	".WMV":   true,
	".AC3":   true,
	".ASF":   true,
	".MOD":   true,
	".MP2":   true,
	".MP3":   true,
	".MP4":   true,
	".WMA":   true,
	".MKA":   true,
	".M4P":   true,
	".WEBM":  true,
}

func SupportedType(p string) bool {
	ext := strings.ToUpper(filepath.Ext(p))
	_, ok := extensions[ext]
	return ok
}
