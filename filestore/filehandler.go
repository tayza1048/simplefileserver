package filestore

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
)

type storage interface {
	save(username string, filename string, data []byte) error
	retrieve(username string, filename string) ([]byte, error)
}

const (
	// StorageOptionMemory means files will be stored in memory
	StorageOptionMemory = 1

	// StorageOptionFileSystem means files will be kept inside file system
	StorageOptionFileSystem = 2
)

var (
	// StorageOption represents setting for storage option
	StorageOption int
	// memory        = make(map[string]map[string]multipart.File)
	currentStorage storage
)

// Upload handles file upload based on the storage option
// func Upload(username string, filename string, file multipart.File) string {
// 	log.Printf("Handling file upload from user %s for file %s ...\n", username, filename)

// 	var path string

// 	switch StorageOption {
// 	case StorageOptionMemory:
// 		addToMemory(username, filename, file)
// 		log.Println("Done storing file in memory.")
// 		path = fmt.Sprintf("%s/%s", username, filename)
// 	case StorageOptionFileSystem:
// 		log.Println("Done storing file in file system.")
// 	}

// 	return path
// }

// // Retrieve returns the data uploaded by users
// func Retrieve(username string, filename string) ([]byte, error) {
// 	if StorageOption == StorageOptionMemory {
// 		file := memory[username][filename]
// 		if file != nil {
// 			return ioutil.ReadAll(file)
// 		}
// 	}

// 	return nil, errors.New("no such file")
// }

// func addToMemory(username string, filename string, file multipart.File) {
// 	mm, ok := memory[username]
// 	if !ok {
// 		mm = make(map[string]multipart.File)
// 		memory[username] = mm
// 	}
// 	mm[filename] = file
// }

// Upload handles file upload based on the storage option
func Upload(username string, filename string, file *multipart.File) (string, error) {
	log.Printf("Handling file upload from user %s for file %s ...\n", username, filename)

	// path
	path := fmt.Sprintf("%s/%s", username, filename)

	// convert multipart to byte array
	data, err := ioutil.ReadAll((*file))
	if err != nil {
		return path, err
	}

	saveErr := getStorage().save(username, filename, data)
	return fmt.Sprintf("%s/%s", username, filename), saveErr
}

// Retrieve returns the data uploaded by users
func Retrieve(username string, filename string) ([]byte, error) {
	return getStorage().retrieve(username, filename)
}

func getStorage() storage {
	switch StorageOption {
	case StorageOptionMemory:
		if currentStorage == nil {
			currentStorage = memorystore{
				memory: make(map[string]map[string][]byte),
			}
		}
	case StorageOptionFileSystem:
		if currentStorage == nil {
			currentStorage = filesystemstore{}
		}
	}

	return currentStorage
}
