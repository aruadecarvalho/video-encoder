package services

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"video-encoder/application/repositories"
	"video-encoder/domain"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	v_path := getVideoPath(v)
	f, err := os.Create(v_path + ".mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	log.Printf("video %v has been stored.", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {
	v_path := getVideoPath(v)
	err := os.Mkdir(v_path, os.ModePerm)
	if err != nil {
		return err
	}

	source := v_path + ".mp4"
	target := v_path + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printBytesOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	v_path := getVideoPath(v)
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, v_path+".frag")
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, v_path)
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printBytesOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	v_path := getVideoPath(v)

	err := os.Remove(v_path + ".mp4")
	if err != nil {
		log.Println("error removing mp4", v.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(v_path + ".frag")
	if err != nil {
		log.Println("error removing frag", v.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(v_path)
	if err != nil {
		log.Println("error removing all paths from", v_path)
		return err
	}

	log.Println("files have been removed: ", v.Video.ID)

	return nil

}

func (v VideoService) InsertVideo() error {
	_, err := v.VideoRepository.Insert(v.Video)
	if err != nil {
		return err
	}

	return nil
}

func getVideoPath(v *VideoService) string {
	v_path := os.Getenv("localStoragePath") + "/" + v.Video.ID

	return v_path
}

func printBytesOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("======> Output: %s\n", string(out))
	}

}
