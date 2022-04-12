package gdfs

import (
	"context"
	"fmt"
	"gdrivecli/pkg/myfile"
	"io"
	"os"

	"google.golang.org/api/drive/v3"
)

func (gdfs *GDFileSystem) DownloadFile(file *drive.File, path string) error {

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Unable to create file: %w", err)
	}

	resp, err := gdfs.Service.Files.Get(file.Id).Download()
	if err != nil {
		return fmt.Errorf("Unable to download file: %w", err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("Unable copy file: %w", err)
	}
	return nil
}

func (gdfs *GDFileSystem) UploadFile(mf *myfile.MyFile) error {

	f, err := os.Open(mf.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	fStat, err := f.Stat()
	if err != nil {
		return err
	}

	_, err = gdfs.Service.Files.
		Create(mf.GFile).
		ResumableMedia(context.Background(), f, fStat.Size(), "").
		ProgressUpdater(func(current, total int64) {
			mf.ProgressBar.Set64(current)
		}).
		Do()
	return err
}
