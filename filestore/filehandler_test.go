package filestore

import (
	"bytes"
	"image"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestResizeImage(t *testing.T) {
	resp, err := http.Get("https://upload.wikimedia.org/wikipedia/en/5/5f/Original_Doge_meme.jpg")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	dataResized := resizeImage(data, 100, 100)
	imageResizedConfig, _, err := image.DecodeConfig(bytes.NewReader(dataResized))
	if err != nil {
		t.Fatal(err)
	}

	if imageResizedConfig.Width != 100 || imageResizedConfig.Height != 100 {
		t.Errorf("Resize doesn't work for given dimensions: got %d, %d want 100, 100", imageResizedConfig.Width, imageResizedConfig.Height)
	}
}
