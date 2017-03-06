package filestore

import "testing"

func TestSaveRetrieve(t *testing.T) {
	StorageOption = StorageOptionMemory

	username := "tayza"
	filename := "cat.txt"
	content := []byte("Here is a string....")
	contentType := "text/plain"

	storage := &memorystore{
		memory: make(map[string]map[string][]byte),
	}

	storage.save(username, filename, content, contentType)
	data, err := storage.retrieve(username, filename)
	if err != nil {
		t.Error(err)
	}

	if len(data) != len(content) {
		t.Error("Expected same length", len(data))
	}
}
