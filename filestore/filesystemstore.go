package filestore

import (
	"io/ioutil"
	"log"
	"os"
)

type filesystemstore struct {
}

func (fs filesystemstore) save(username string, filename string, data []byte, contentType string) error {
	if _, err := os.Stat(username); err != nil {
		os.Mkdir(username, 0700)
	}

	file, err := os.OpenFile(username+"/"+filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)
	log.Println("Done storing file in file system.")

	return nil
}

func (fs filesystemstore) retrieve(username string, filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(username + "/" + filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}
