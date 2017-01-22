# Simple File Server

This is a simple project to experiment with Golang programming language and some third-party APIs.

In a nutshell, the server can handle:
- Uploading and retrieving files to and from the server using different storage options such as memory, file system and Amazon S3 storage.
- Resizing of image uploads according to the given dimension before storing the files.

### Usage
The project is built with Go 1.7. To avoid conflicts with your existing dependencies, please create a new directory and point your GOPATH to the directory.

```sh
$ cd $GOPATH
$ mkdir src
$ go get github.com/tayza1048/simplefileserver
$ cd src/github.com/tayza1048/simplefileserver
$ go run webserver.go
```

An alternative approach is to run "go install github.com/tayza1048/simplefileserver" and run the binary output file. But it also requires copying templates and setting files as an extra step.

Once the server starts, please go to http://{hostname}:{port}/upload and try out the functionalities.

### Settings
Modify settings.json file to configure server hostname and port. For storage options, please use one of the following:
- memory
- filesystem
- s3

If the option is "s3", it is required to configure Amazon S3 settings in settings_s3.json file.