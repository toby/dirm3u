# dirm3u

Simple HTTP media server. Turns a directory into a streamable playlist.

## Installing

`go get github.com/toby/dirm3u`

## Usage

`dirm3u` will recursively search the working directory for media files, host them and
generate a .m3u playlist for all compatible types (see [extensions.go](extensions.go)).

```
dirm3u [FLAG]
  -h string
    	hostname (default "localhost")
  -p int
    	http server port (default 20202)
```

## Endpoints

Once running, you can hit the following endpoints:

* [http://localhost:20202/playlist.m3u](http://localhost:20202/playlist.m3u) The auto-generated .m3u playlist.
* [http://localhost:20202/reload](http://localhost:20202/reload) Rescan the current directory for an update file list.
* [http://localhost:20202/media/MEDIA](http://localhost:20202/media/MEDIA) Streamable hosted media. Replace MEDIA with a valid file path.
