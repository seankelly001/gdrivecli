package myfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/schollz/progressbar/v3"
	drive "google.golang.org/api/drive/v3"
)

type MyFile struct {
	Name        string
	GFile       *drive.File
	Path        string
	ProgressBar *progressbar.ProgressBar
	Done        bool
}

func NewFile(gf *drive.File, path string) *MyFile {

	myFile := &MyFile{
		Name:  gf.Name,
		GFile: gf,
		Path:  path,
		Done:  false,
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
			myFile.Done = true
		}),
	)

	myFile.ProgressBar = bar

	return myFile
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
