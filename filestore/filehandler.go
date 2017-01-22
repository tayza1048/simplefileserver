package filestore

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"mime/multipart"

	"github.com/nfnt/resize"
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
	StorageOption  int
	currentStorage storage
)

// Upload handles file upload based on the storage option
func Upload(username string, filename string, file *multipart.File, width uint, height uint) (string, error) {
	log.Printf("Handling file upload from user %s for file %s ...\n", username, filename)

	// path
	path := fmt.Sprintf("%s/%s", username, filename)

	// convert multipart to byte array
	data, err := ioutil.ReadAll((*file))
	if err != nil {
		return path, err
	}

	// resize image if required
	data = resizeImage(data[0:], width, height)

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

func resizeImage(data []byte, width uint, height uint) []byte {
	if image, format, err := image.Decode(bytes.NewReader(data)); err == nil && width > 0 && height > 0 {
		isJpeg := format == "jpeg"
		isPng := format == "png"

		if isJpeg || isPng {
			resizeImage := resize.Resize(width, height, image, resize.Lanczos3)
			buf := new(bytes.Buffer)

			if isJpeg {
				if err := jpeg.Encode(buf, resizeImage, nil); err == nil {
					log.Println("jpeg resized")
					return buf.Bytes()
				}
			} else {
				if err := png.Encode(buf, resizeImage); err == nil {
					log.Println("png resized")
					return buf.Bytes()
				}
			}
		}
	}
	return data
}
