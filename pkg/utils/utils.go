package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ermanimer/progress_bar"
	"github.com/minio/minio/pkg/disk"
	"github.com/rivo/tview"
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

// func ChooseFolder() (string, error) {
// 	result, err := cfdutil.ShowPickFolderDialog(cfd.DialogConfig{
// 		Title:  "Pick Folder",
// 		Role:   "PickFolderExample",
// 		Folder: "C:\\",
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	log.Printf("Chosen folder: %s\n", result)
// 	return result, nil
// }

func GetDiskUsage(diskName string) (string, error) {
	di, err := disk.GetInfo(diskName)
	if err != nil {
		return "", err
	}
	freeSpace := int64(di.Free)
	return ByteCountIEC(freeSpace), nil
}

func FilterNodeChildren(children []*tview.TreeNode, key string) *tview.TreeNode {

	key = strings.ToLower(key)
	previousNode := children[0]
	for _, node := range children {
		text := strings.ToLower(node.GetText())
		firstChar := text[:1]
		if firstChar == key {
			return node
		} else if firstChar > key {
			return previousNode
		}
		previousNode = node
	}
	return children[len(children)-1]
}
