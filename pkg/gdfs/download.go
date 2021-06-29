package gdfs

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/api/drive/v3"
)

func (gdfs *GDFileSystem) DownloadFile(file *drive.File, path string) error {

	//file := gdfs.FilesToDownload[fileName]
	//path := filepath.Join(root, file.Name)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Unable to create file: %w", err)
	}

	// waitChan := make(chan struct{})
	// ctx, cancel := context.WithCancel(context.Background())
	// go utils.DisplayProgress(ctx, waitChan, f, file.Size)

	resp, err := gdfs.Service.Files.Get(file.Id).Download()
	if err != nil {
		return fmt.Errorf("Unable to download file: %w", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("Unable copy file: %w", err)
	}

	// cancel()
	// <-waitChan
	return nil
}
