package myfile

import (
	"fmt"
	"gdrivecli/pkg/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/rivo/tview"
	"github.com/schollz/progressbar/v3"
	drive "google.golang.org/api/drive/v3"
)

//This will be part of NodeRef
type MyFile struct {
	Name        string
	GFile       *drive.File
	Path        string
	ProgressBar *progressbar.ProgressBar
	Done        bool
	InProgress  bool
	Node        *tview.TreeNode
}

//func NewDownloadFile(gf *drive.File, path string) *MyFile {
func NewDownloadFile(gf *drive.File, node *tview.TreeNode) *MyFile {

	myFile := &MyFile{
		Name:       gf.Name,
		GFile:      gf,
		Node:       node,
		Done:       false,
		InProgress: false,
	}

	cols, _ := consolesize.GetConsoleSize()
	var maxLen, pbWidth int

	if cols > 100 {
		maxLen = 50
		pbWidth = cols - maxLen - 53
	} else {
		maxLen = 30
		pbWidth = 10
	}

	desc := fmt.Sprintf("%-*v", maxLen, gf.Name)
	desc = desc[:maxLen]
	bar := progressbar.NewOptions64(
		gf.Size,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(ioutil.Discard),
		progressbar.OptionSetWidth(pbWidth),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionOnCompletion(func() {
			myFile.InProgress = false
			myFile.Done = true
			myFile.Node.SetColor(config.TREE_DOWNLOADED_COLOUR)
		}),
	)
	myFile.ProgressBar = bar
	return myFile
}

//func NewDownloadFile(gf *drive.File, path string) *MyFile {
func NewUploadFile(path string, node *tview.TreeNode) (*MyFile, error) {

	myFile := &MyFile{
		Name:       filepath.Base(path),
		Path:       path,
		Node:       node,
		Done:       false,
		InProgress: false,
	}

	cols, _ := consolesize.GetConsoleSize()
	var maxLen, pbWidth int

	if cols > 100 {
		maxLen = 50
		pbWidth = cols - maxLen - 53
	} else {
		maxLen = 30
		pbWidth = 10
	}

	file, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	desc := fmt.Sprintf("%-*v", maxLen, path)
	desc = desc[:maxLen]
	bar := progressbar.NewOptions64(
		file.Size(),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(ioutil.Discard),
		progressbar.OptionSetWidth(pbWidth),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionOnCompletion(func() {
			myFile.InProgress = false
			myFile.Done = true
			myFile.Node.SetColor(config.TREE_DOWNLOADED_COLOUR)
		}),
	)
	myFile.ProgressBar = bar
	return myFile, nil
}

func (mf *MyFile) UpdateProgress() error {

	fi, err := os.Stat(mf.Path)
	if err != nil {
		return err
	}
	currentFileSize := fi.Size()
	mf.ProgressBar.Set64(currentFileSize)
	return nil
}
