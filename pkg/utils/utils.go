package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ermanimer/progress_bar"
	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"
	"github.com/minio/minio/pkg/disk"
)

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func DisplayProgress(ctx context.Context, waitChan chan struct{}, f *os.File, totalSize int64) {

	//create new progress bar
	pb := progress_bar.DefaultProgressBar(float64(totalSize) / 1024 / 1024)
	//start
	err := pb.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//update
	for {
		select {
		case <-ctx.Done():
			pb.Update(float64(totalSize) / 1024 / 1024)
			pb.Stop()
			waitChan <- struct{}{}
			return
		default:
			fi, err := f.Stat()
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			currentFileSize := fi.Size()
			err = pb.Update(float64(currentFileSize) / 1024 / 1024)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func ChooseFolder() (string, error) {
	result, err := cfdutil.ShowPickFolderDialog(cfd.DialogConfig{
		Title:  "Pick Folder",
		Role:   "PickFolderExample",
		Folder: "C:\\",
	})
	if err != nil {
		return "", err
	}
	log.Printf("Chosen folder: %s\n", result)
	return result, nil
}

func GetDiskUsage(diskName string) (string, error) {
	di, err := disk.GetInfo(diskName)
	if err != nil {
		return "", err
	}
	freeSpace := int64(di.Free)
	return ByteCountIEC(freeSpace), nil
}
