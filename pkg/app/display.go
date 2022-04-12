package app

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"gdrivecli/pkg/utils"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/disk"
)

func (app *App) Prompt(label string, doneFunc func(key tcell.Key)) {
	app.Input.SetLabel(label)
	app.Input.SetDoneFunc(doneFunc)
	app.tvApp.SetFocus(app.Input)
}

func (a *App) WriteOutput(output string) {
	a.Output.SetText(fmt.Sprintf("%s%s", a.Output.GetText(false), output))
	a.Output.ScrollToEnd()
}

func (a *App) GetInputText() string {
	return a.Input.GetText()
}

func (a *App) ResetInput() {
	a.Input.SetLabel("")
	a.Input.SetText("")
}

func (a *App) DisplayHelp() {
	var b bytes.Buffer
	buf := bufio.NewWriter(&b)
	w := tabwriter.NewWriter(buf, 1, 1, 1, ' ', 0)

	fmt.Fprintln(w, "Google Drive\tCommon\tThis PC")
	fmt.Fprintln(w, "Enter: Select file for download\tTab: Switch between Google Drive and Local\tInsert: Create folder")
	fmt.Fprintln(w, "Ctrl + A: Select all files in folder\tCtrl + H: Display this help\tDelete: Delete folder")
	fmt.Fprintln(w, "Ctrl + U: Deselect all files in folder\t\tCtrl + S: Save all selected files to current folder")
	fmt.Fprintln(w, "Ctrl + N: Sort by name (reversible)\t\t")
	fmt.Fprintln(w, "Ctrl + N: Sort by date (reversible)\t\t")

	w.Flush()
	buf.Flush()
	a.WriteOutput(b.String())
}

func (a *App) DisplayInfo() {
	var msg string = ""
	partitions, err := disk.Partitions(false)
	if err != nil {
		msg = err.Error()
	} else {
		for _, partition := range partitions {

			diskName := partition.Device
			freeDiskSpace, err := utils.GetDiskUsage(diskName)
			if err != nil {
				msg = err.Error()
				break
			}
			msg = fmt.Sprintf("%s%s %s\n", msg, diskName, freeDiskSpace)
		}

		var numFilesSelected, numFilesInProgress, numFilesDownloaded int = 0, 0, 0
		var totalDownloadSizeSelected, totalDownloadSizeInProgress, totalDownloadSizeDownloaded int64 = 0, 0, 0
		for _, mf := range a.FilesToDownload {
			if mf.Done {
				numFilesDownloaded += 1
				totalDownloadSizeDownloaded += mf.GFile.Size
			} else if mf.InProgress {
				numFilesInProgress += 1
				totalDownloadSizeInProgress += mf.GFile.Size
			} else {
				numFilesSelected += 1
				totalDownloadSizeSelected += mf.GFile.Size
			}
		}
		msg = fmt.Sprintf("%sNum files selected: %d Total size: %s\n", msg, numFilesSelected, utils.ByteCountIEC(totalDownloadSizeSelected))
		msg = fmt.Sprintf("%sNum files in progress: %d Total size: %s\n", msg, numFilesInProgress, utils.ByteCountIEC(totalDownloadSizeInProgress))
		msg = fmt.Sprintf("%sNum files downloaded: %d Total size: %s\n", msg, numFilesDownloaded, utils.ByteCountIEC(totalDownloadSizeDownloaded))
	}
	a.Info.SetText(strings.TrimSpace(msg))

}

func (a *App) StartDownloadProgressLoop(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			downloadProgress := ""

			//Need to sort by key to stop downloads from jumping place as maps are unordered
			keys := make([]string, len(a.FilesToDownload))
			i := 0
			for k := range a.FilesToDownload {
				keys[i] = k
				i++
			}
			sort.Strings(keys)

			for _, mfk := range keys {
				mf := a.FilesToDownload[mfk]
				if mf.InProgress && !mf.Done {
					err := mf.UpdateProgress()
					if err != nil {
						a.WriteOutput(err.Error())
					}
				}
				downloadProgress = downloadProgress + mf.ProgressBar.String() + "\n"
			}
			a.DownloadProgress.SetText(downloadProgress)
			a.DownloadProgress.ScrollToEnd()
			a.tvApp.Draw()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func (a *App) StartUploadProgressLoop(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			uploadProgress := ""

			//Need to sort by key to stop downloads from jumping place as maps are unordered
			keys := make([]string, len(a.FilesToUpload))
			i := 0
			for k := range a.FilesToUpload {
				keys[i] = k
				i++
			}
			sort.Strings(keys)

			for _, mfk := range keys {
				mf := a.FilesToUpload[mfk]
				// if mf.InProgress && !mf.Done {
				// 	err := mf.UpdateProgress()
				// 	if err != nil {
				// 		a.WriteOutput(err.Error())
				// 	}
				// }
				uploadProgress = uploadProgress + mf.ProgressBar.String() + "\n"
			}
			a.UploadProgress.SetText(uploadProgress)
			a.UploadProgress.ScrollToEnd()
			a.tvApp.Draw()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func (a *App) StartDisplayInfoLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			a.DisplayInfo()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
