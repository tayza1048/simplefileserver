package filestore

import (
	"errors"
	"log"
)

type memorystore struct {
	memory map[string]map[string][]byte
}

func (ms memorystore) save(username string, filename string, data []byte) error {
	mm, ok := ms.memory[username]
	if !ok {
		mm = make(map[string][]byte)
		ms.memory[username] = mm
	}
	mm[filename] = data

	log.Println("Done storing file in memory.")
	return nil
}

func (ms memorystore) retrieve(username string, filename string) ([]byte, error) {
	if StorageOption == StorageOptionMemory {
		if data, ok := ms.memory[username][filename]; ok {
			return data, nil
		}
	}

	return nil, errors.New("no such file")
}
