package main

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tayza1048/simplefileserver/filestore"
)

func TestEmptyDownloadHandler(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(download)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := ""
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUploadDownload(t *testing.T) {
	setup()

	// Get the image to upload and convert as temp file
	resp, err := http.Get("https://upload.wikimedia.org/wikipedia/en/5/5f/Original_Doge_meme.jpg")
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	// Prepare form data
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", "doge.jpg")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(fw, resp.Body); err != nil {
		t.Fatal(err)
	}

	fw, err = w.CreateFormField("username")
	fw.Write([]byte("tayza"))
	w.Close()

	// Create request to upload
	req, err := http.NewRequest("POST", "/upload", &b)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", w.FormDataContentType())

	// Create recorder and send request
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(upload)
	handler.ServeHTTP(rr, req)

	// Check uploading
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "http://localhost:6061/tayza/doge.jpg"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Create request to download
	reqDl, err := http.NewRequest("GET", "/tayza/doge.jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	rrDl := httptest.NewRecorder()
	handlerDl := http.HandlerFunc(download)
	handlerDl.ServeHTTP(rrDl, reqDl)

	// Check downloading
	if status := rrDl.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rrDl.Header().Get("Content-Type") != "image/jpeg" {
		t.Errorf("Wrong Content-Type: got %v want %v", rrDl.Header().Get("Content-Type"), "image/jpeg")
	}

	dataDl, err := ioutil.ReadAll(rrDl.Body)
	if err != nil {
		t.Fatal(err)
	}

	imageConfig, _, err := image.DecodeConfig(bytes.NewReader(dataDl))
	if err != nil {
		t.Fatal(err)
	}

	originalWidth := 369
	if imageConfig.Width != originalWidth {
		t.Errorf("Downloaded image has a different width: got %v want %v", imageConfig.Width, originalWidth)
	}

	originalHeight := 273
	if imageConfig.Height != originalHeight {
		t.Errorf("Downloaded image has a different height: got %v want %v", imageConfig.Height, originalHeight)
	}
}

func setup() {
	hostname = defaultHostname
	port = defaultPort
	filestore.StorageOption = filestore.StorageOptionMemory
}
