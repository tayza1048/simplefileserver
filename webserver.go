package main

import (
	"io"
	"log"
	"net/http"

	"github.com/tayza1048/simplefileserver/filestore"
)

var (
	hostname string
	port     string
)

func handleDefault(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Please call /hello")
}

func handleHello(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello. I am a simple file server.\n")
}

func upload(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		req.ParseMultipartForm(32 << 20)

		username := req.Form["username"]
		if len(username) == 0 {
			http.Error(w, "Please provide your api username.", http.StatusInternalServerError)
			return
		}

		file, handler, err := req.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Print(handler.Header)
		io.WriteString(w, "http://"+hostname+":"+port+"/"+filestore.Upload(username[0], handler.Filename, file))
	} else {
		http.Error(w, "Please use POST or PUT requests to upload files.", http.StatusInternalServerError)
	}
}

func main() {
	loadSettings()
	initializeHandlers()

	// start web server
	err := http.ListenAndServe(":6061", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func loadSettings() {
	hostname = "localhost"
	port = "6061"
	filestore.StorageOption = filestore.StorageOptionMemory
}

func initializeHandlers() {
	http.HandleFunc("/", handleDefault)
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/upload", upload)
}
