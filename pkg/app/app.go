package app

import (
	"context"
	"fmt"
	"gdrivecli/pkg/config"
	"gdrivecli/pkg/file_system"
	"gdrivecli/pkg/gdfs"
	"gdrivecli/pkg/myfile"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	tvApp            *tview.Application
	Input            *tview.InputField
	Output           *tview.TextView
	Info             *tview.TextView
	DownloadProgress *tview.TextView
	UploadProgress   *tview.TextView
	FSTree           *tview.TreeView
	GDFSTree         *tview.TreeView
	FS               *file_system.FileSystem
	GDFS             *gdfs.GDFileSystem
	FilesToDownload  map[string]*myfile.MyFile
	FilesToUpload    map[string]*myfile.MyFile
}

// type Tree interface {
// 	//NodeSelectedFun
// }

func NewApp(gdfs *gdfs.GDFileSystem) (*App, error) {

	fs := file_system.NewFS()
	app := &App{
		GDFS:            gdfs,
		FS:              fs,
		FilesToDownload: make(map[string]*myfile.MyFile),
		FilesToUpload:   make(map[string]*myfile.MyFile),
	}

	tvApp := tview.NewApplication()
	inputView := tview.NewInputField()
	inputView.SetLabelColor(config.OUTPUT_COLOUR).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(config.INPUT_COLOUR).
		SetBorderColor(config.BORDER_COLOUR).
		SetBorder(true).
		SetTitle("input").
		SetTitleColor(config.TITLE_COLOUR)

	outputView := tview.NewTextView()
	outputView.SetText("Welcome to GDrive CLI").
		SetTextColor(config.OUTPUT_COLOUR).
		SetWrap(true).
		SetBorderColor(config.BORDER_COLOUR).
		SetBorder(true).
		SetTitle("output").
		SetTitleColor(config.TITLE_COLOUR)

	infoView := tview.NewTextView()
	infoView.SetText("").
		SetTextColor(config.OUTPUT_COLOUR).
		SetWrap(true).
		SetBorderColor(config.BORDER_COLOUR).
		SetBorder(true).
		SetTitle("info").
		SetTitleColor(config.TITLE_COLOUR)

	downloadProgressView := tview.NewTextView()
	downloadProgressView.
		SetTextColor(config.OUTPUT_COLOUR).
		SetBorderColor(config.BORDER_COLOUR).
		SetBorder(true).
		SetTitle("download progress").
		SetTitleColor(config.TITLE_COLOUR)

	uploadProgressView := tview.NewTextView()
	uploadProgressView.
		SetTextColor(config.OUTPUT_COLOUR).
		SetBorderColor(config.BORDER_COLOUR).
		SetBorder(true).
		SetTitle("upload progress").
		SetTitleColor(config.TITLE_COLOUR)

	fsTree, err := NewFSTree(app)
	if err != nil {
		return nil, fmt.Errorf("err creating FS tree: %w", err)
	}

	gdfsTree := NewGDFSTree(app)

	//Flex for the 2 tree views
	treeFlex := tview.NewFlex().
		AddItem(gdfsTree, 0, 1, false).
		AddItem(fsTree, 0, 1, true)

	//Func to switch focus between the 2 trees
	treeFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyTab:
			if gdfsTree.HasFocus() {
				tvApp.SetFocus(fsTree)
			} else {
				tvApp.SetFocus(gdfsTree)
			}
		case tcell.KeyCtrlH:
			app.DisplayHelp()
		default:
			return event
		}
		return nil
	})

	//Flex for text views
	textFlex := tview.NewFlex().
		AddItem(inputView, 0, 1, false).
		AddItem(outputView, 0, 3, false).
		AddItem(infoView, 0, 1, false)

	outerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(treeFlex, 0, 10, true).
		AddItem(textFlex, 0, 3, false).
		AddItem(downloadProgressView, 0, 4, false).
		AddItem(uploadProgressView, 0, 3, false)

	tvApp.SetRoot(outerFlex, true).SetFocus(outerFlex)

	app.Input = inputView
	app.Output = outputView
	app.Info = infoView
	app.DownloadProgress = downloadProgressView
	app.UploadProgress = uploadProgressView
	app.tvApp = tvApp
	app.FSTree = fsTree
	app.GDFSTree = gdfsTree

	return app, nil
}

// Start the app
// Start download progres loop
// Start display info loop
func (a *App) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	go a.StartDownloadProgressLoop(ctx)
	go a.StartUploadProgressLoop(ctx)
	go a.StartDisplayInfoLoop(ctx)
	a.DisplayInfo()
	a.DisplayHelp()

	if err := a.tvApp.Run(); err != nil {
		log.Fatalf("Unable to start application: %v", err)
	}
	cancel()
}
