package gcli

import (
	"context"
	"fmt"
	"gdrivecli/pkg/utils"
	"io"
	"log"
	"os"
)

func (gcli *GDriveCLI) DownloadFile(rootPath, fileName string) {

	file := gcli.FilesToDownload[fileName]
	fmt.Printf("Downloading %s...\n", file.Name)
	path := rootPath + file.Name
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
	}

	waitChan := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go utils.DisplayProgress(ctx, waitChan, f, file.Size)

	resp, err := gcli.Service.Files.Get(file.Id).Download()
	if err != nil {
		log.Fatalf("Unable to download file: %v", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatalf("Unable copy file: %v", err)
	}

	cancel()
	<-waitChan
}
