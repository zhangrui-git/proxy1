package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func main() {
	var host string
	var listen uint
	flag.StringVar(&host, "h", "http://127.0.0.1", "web server url")
	flag.UintVar(&listen, "l", 8888, "localhost listen port")
	flag.Parse()

	log.Printf("service start http://127.0.0.1:%d", listen)
	url := Url{Host: host}
	http.Handle("/", http.FileServer(url))
	addr := ":" + strconv.Itoa(int(listen))
	log.Fatal(http.ListenAndServe(addr, nil))
}

type Url struct {
	Host string
}

func (u Url) Open(name string) (http.File, error) {
	log.Printf("[request name] %s", name)
	basePath := "./file"
	filePath := basePath + name
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		log.Printf("[file name] %s", fileInfo.Name())
		log.Printf("[open path] %s", filePath)
		file, err := os.Open(filePath)
		return file, err
	}

	url := u.Host + name
	log.Printf("[download url] %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	d,_ := path.Split(filePath)
	err = os.MkdirAll(d, os.ModePerm)
	if err != nil {
		log.Panicln(err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		 log.Panicln(err)
	}

	log.Printf("[save file] %s", filePath)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	return out, err
}
