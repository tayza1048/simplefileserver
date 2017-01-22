package filestore

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
)

const (
	// StorageOptionMemory means files will be stored in memory
	StorageOptionMemory = 1
)

var (
	// StorageOption represents setting for storage option
	StorageOption int
	memory        = make(map[string]map[string]multipart.File)
)

// Upload handles file upload based on the storage option
func Upload(username string, filename string, file multipart.File) string {
	log.Printf("Handling file upload from user: %s ...\n", username)

	var path string

	if StorageOption == StorageOptionMemory {
		addToMemory(username, filename, file)
		log.Println("Done storing file in memory.")
		path = fmt.Sprintf("%s/%s", username, filename)
	}

	return path
}

// Retrieve returns the data uploaded by users
func Retrieve(username string, filename string) ([]byte, error) {
	if StorageOption == StorageOptionMemory {
		file := memory[username][filename]
		if file != nil {
			return ioutil.ReadAll(file)
		}
	}

	return nil, errors.New("no such file")
}

func addToMemory(username string, filename string, file multipart.File) {
	mm, ok := memory[username]
	if !ok {
		mm = make(map[string]multipart.File)
		memory[username] = mm
	}
	mm[filename] = file
}
