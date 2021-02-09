package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ermanimer/progress_bar"
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
			//forever <- struct{}{}
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
