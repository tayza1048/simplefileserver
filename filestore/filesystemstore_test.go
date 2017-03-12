package filestore

import "testing"

func TestFileSystemStoreSaveRetrieve(t *testing.T) {
	StorageOption = StorageOptionFileSystem

	username := "tayza"
	filename := "cat.txt"
	content := []byte("Here is a string....")
	contentType := "text/plain"

	storage := filesystemstore{}

	err := storage.save(username, filename, content, contentType)
	if err != nil {
		t.Error(err)
	}

	data, err := storage.retrieve(username, filename)
	if err != nil {
		storage.delete(username, filename)
		t.Error(err)
	}
	storage.delete(username, filename)

	if len(data) != len(content) {
		t.Error("Expected same length", len(data))
	}
}
