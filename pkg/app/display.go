package app

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
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

func (a *App) StartDownloadProgressLoop(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			downloadProgress := ""
			for _, mf := range a.FilesToDownload {
				if !mf.Done {
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
