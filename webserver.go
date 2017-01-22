package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/tayza1048/simplefileserver/filestore"
)

var (
	hostname string
	port     string
)

func handleDefault(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "I am a simple file server.\n")
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

func download(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		path := strings.Split(req.URL.Path, "/")
		if len(path) == 3 {
			log.Printf("Handling download from user %s for file %s ...\n", path[1], path[2])
			data, err := filestore.Retrieve(path[1], path[2])
			if err != nil {
				http.NotFound(w, req)
				return
			}

			w.Header().Set("Content-Type", http.DetectContentType(data))
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			w.Write(data)
		}
	} else {
		http.Error(w, "Please use Get requests to retrieve files.", http.StatusInternalServerError)
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
	http.HandleFunc("/", download)
	http.HandleFunc("/upload", upload)
}
