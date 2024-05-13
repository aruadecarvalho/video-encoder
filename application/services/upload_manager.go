package services

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) LoadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.LoadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.UploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for i := 0; i < len(vu.Paths); i++ {
			in <- i
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) UploadWorker(in chan int, returnChan chan string, uploadClient *storage.Client, ctx context.Context) {
	for i := range in {
		path := vu.Paths[i]
		err := vu.UploadObject(path, uploadClient, ctx)
		if err != nil {
			vu.Errors = append(vu.Errors, path)
			log.Printf("error during upload: %v. Error: %v", path, err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload completed"
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
