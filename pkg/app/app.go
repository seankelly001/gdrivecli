package app

import (
	"context"
	"fmt"
	"gdrivecli/pkg/config"
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
	DownloadProgress *tview.TextView
	FSTree           *tview.TreeView
	GDFSTree         *tview.TreeView
	GDFS             *gdfs.GDFileSystem
	FilesToDownload  []*myfile.MyFile
}

func NewApp(gdfs *gdfs.GDFileSystem) (*App, error) {

	app := &App{
		GDFS: gdfs,
	}

	tvApp := tview.NewApplication()
	inputView := tview.NewInputField()
	inputView.SetLabelColor(config.OUTPUT_COLOUR).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(config.INPUT_COLOUR).
		//SetColor
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

	downloadProgressView := tview.NewTextView()
	downloadProgressView.
		SetTextColor(config.OUTPUT_COLOUR).
		SetBorderColor(config.BORDER_COLOUR).
		SetBorder(true).
		SetTitle("download progress").
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
		if event.Key() == tcell.KeyTab {
			if gdfsTree.HasFocus() {
				tvApp.SetFocus(fsTree)
			} else {
				tvApp.SetFocus(gdfsTree)
			}
		}
		return event
	})

	//Flex for text views
	textFlex := tview.NewFlex().
		AddItem(inputView, 0, 1, false).
		AddItem(outputView, 0, 3, false)

	outerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(treeFlex, 0, 10, true).
		AddItem(textFlex, 0, 2, false).
		AddItem(downloadProgressView, 0, 4, false)

	tvApp.SetRoot(outerFlex, true).SetFocus(outerFlex)

	app.Input = inputView
	app.Output = outputView
	app.DownloadProgress = downloadProgressView
	app.tvApp = tvApp
	app.FSTree = fsTree
	app.GDFSTree = gdfsTree
	return app, nil
}

func (a *App) Run() {

	ctx, cancel := context.WithCancel(context.Background())
	go a.StartDownloadProgressLoop(ctx)

	if err := a.tvApp.Run(); err != nil {
		log.Fatalf("Unable to start application: %v", err)
	}
	cancel()
}
