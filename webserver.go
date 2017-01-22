package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/tayza1048/simplefileserver/filestore"
)

var (
	hostname string
	port     string
)

const (
	defaultHostname = "localhost"
	defaultPort     = "6061"
)

func main() {
	initializeHandlers()
	loadSettings()

	// start web server
	err := http.ListenAndServe(":6061", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func loadSettings() {
	file, err := os.Open("settings.json")
	if err != nil {
		hostname = defaultPort
		port = defaultPort
		filestore.StorageOption = filestore.StorageOptionMemory
	} else {
		defer file.Close()

		dec := json.NewDecoder(file)
		for {
			var v map[string]interface{}
			if err := dec.Decode(&v); err != nil {
				log.Println(err)
				return
			}
			for k, value := range v {
				switch k {
				case "hostname":
					if str, ok := value.(string); ok {
						hostname = str
					}
				case "port":
					if str, ok := value.(string); ok {
						port = str
					}
				case "storageOption":
					var storage string
					if str, ok := value.(string); ok {
						storage = str
						switch storage {
						case "memory":
							filestore.StorageOption = filestore.StorageOptionMemory
						case "filesystem":
							filestore.StorageOption = filestore.StorageOptionFileSystem
						default:
							log.Fatal("Invalid settings")
						}
					}

				}
			}
		}
	}
}

func initializeHandlers() {
	http.HandleFunc("/", download)
	http.HandleFunc("/upload", upload)
}

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

		width := parseParamInt(req.Form["width"])
		height := parseParamInt(req.Form["height"])

		path, err := filestore.Upload(username[0], handler.Filename, &file, width, height)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		io.WriteString(w, "http://"+hostname+":"+port+"/"+path)
	} else {
		http.Error(w, "Please use POST or PUT requests to upload files.", http.StatusInternalServerError)
	}
}

func parseParamInt(formArr []string) uint {
	if len(formArr) != 0 {
		if u64, err := strconv.ParseUint(formArr[0], 10, 32); err == nil {
			return uint(u64)
		}
	}

	return 0
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
