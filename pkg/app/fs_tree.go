package app

import (
	"fmt"
	"gdrivecli/pkg/config"
	"gdrivecli/pkg/file_system"
	fs "gdrivecli/pkg/file_system"
	"gdrivecli/pkg/myfile"
	"gdrivecli/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/drive/v3"
)

func NewFSTree(a *App) (*tview.TreeView, error) {
	root := tview.NewTreeNode("This PC").SetColor(config.TREE_DIR_COLOUR)
	root.SetReference(fs.FileReference{
		FilePath: "",
		FileName: "",
	})

	pcName, err := os.Hostname()
	if err != nil {
		pcName = "This PC"
	}
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetTitle(pcName)
	tree.SetTitleColor(config.TITLE_COLOUR)
	tree.SetGraphicsColor(config.TREE_GRAPHICS_COLOUR)
	tree.SetBorderColor(config.BORDER_COLOUR)
	tree.SetSelectedFunc(a.FSNodeSelectedFunc)
	tree.SetBorder(true)
	tree.SetInputCapture(a.FSInputFunc)

	partitionNodes, err := fs.GetPartitionNodes()
	if err != nil {
		return nil, err
	}
	for _, node := range partitionNodes {
		root.AddChild(node)
	}

	return tree, nil
}

func (a *App) FSNodeSelectedFunc(node *tview.TreeNode) {

	ref := node.GetReference()
	nodeRef, ok := ref.(fs.FileReference)
	if !ok {
		a.WriteOutput("error casting")
		return
	}
	f, err := os.Stat(nodeRef.FilePath)
	if err != nil {
		a.WriteOutput(err.Error())
		return
	}
	if f.IsDir() {
		if len(node.GetChildren()) == 0 {
			err := fs.SetFSChildren(node)
			if err != nil {
				a.WriteOutput(err.Error())
			}
		}
		node.SetExpanded(!node.IsExpanded())
	} else {
		//TODO
	}
}

//Determine what to do when a node has input from keyboard
func (a *App) FSInputFunc(event *tcell.EventKey) *tcell.EventKey {

	cur := a.FSTree.GetCurrentNode()
	switch event.Key() {
	case tcell.KeyInsert:
		a.CreateDirPrompt(cur)
	case tcell.KeyDelete:
		a.DeletePrompt(cur)
	case tcell.KeyCtrlS:
		a.SavePrompt(cur)
	}

	return event
}

func (a *App) CreateDirPrompt(node *tview.TreeNode) error {

	reference := node.GetReference()
	nodeReference, ok := reference.(fs.FileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	path := nodeReference.FilePath
	prompt := "Enter dir name: "

	doneFunc := func(key tcell.Key) {
		if key == tcell.KeyEnter {
			//Create directory with name from input, write message, reset nodes children and expand node
			newDirPath := filepath.Join(path, a.GetInputText())
			msg := file_system.CreateDir(newDirPath)
			a.WriteOutput(msg)
			cur := a.FSTree.GetCurrentNode()
			fs.SetFSChildren(cur)
			cur.Expand()
		}
		//Reset
		a.ResetInput()
		a.tvApp.SetFocus(a.FSTree)
	}
	a.Prompt(prompt, doneFunc)
	return nil
}

func (a *App) DeletePrompt(node *tview.TreeNode) error {
	reference := node.GetReference()
	nodeReference, ok := reference.(fs.FileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	path := nodeReference.FilePath
	fileName := nodeReference.FileName
	isDir, err := file_system.IsDir(path)
	if err != nil {
		return err
	}
	var prompt string
	if isDir {
		prompt = fmt.Sprintf("Delete directory %s ? (y/n): ", fileName)
	} else {
		prompt = fmt.Sprintf("Delete file %s ? (y/n): ", fileName)
	}

	doneFunc := func(key tcell.Key) {
		if key == tcell.KeyEnter {
			//Create directory with name from input, write message, reset nodes children and expand node
			answer := strings.ToLower(a.GetInputText())
			if answer == "y" {
				//delete dir
				err := file_system.Delete(path)
				if err != nil {
					msg := fmt.Sprintf("error deleting directory: %s", err.Error())
					a.WriteOutput(msg)
				} else {
					msg := fmt.Sprintf("%s deleted", path)
					a.WriteOutput(msg)
					//reload parent children, and set current node to parent
					parent := nodeReference.Parent
					fs.SetFSChildren(parent)
					a.FSTree.SetCurrentNode(parent)
				}
				//Reset
				a.ResetInput()
				a.tvApp.SetFocus(a.FSTree)
			} else if answer == "n" {
				a.ResetInput()
				a.tvApp.SetFocus(a.FSTree)
			} else {
				//invalid answer
				a.WriteOutput("invalid key")
				a.Input.SetText("")
			}
		}

	}
	a.Prompt(prompt, doneFunc)
	return nil
}

func (a *App) SavePrompt(node *tview.TreeNode) error {
	reference := node.GetReference()
	nodeReference, ok := reference.(fs.FileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	path := nodeReference.FilePath
	totalDownloadSizeBytes := a.GDFS.GetTotalDownloadSizeBytes()
	totalDownloadSize := utils.ByteCountIEC(totalDownloadSizeBytes)
	diskName := filepath.VolumeName(path)
	freeDiskSpace, err := utils.GetDiskUsage(diskName)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Downloading %d files. Total size: %s. Free space on disk %s %s", len(a.GDFS.FilesToDownload), totalDownloadSize, diskName, freeDiskSpace)
	a.WriteOutput(msg)

	for _, gf := range a.GDFS.FilesToDownload {

		fullPath := filepath.Join(path, gf.Name)
		go a.DownloadFile(gf, fullPath)
		myFile := myfile.NewFile(gf, fullPath)
		a.FilesToDownload = append(a.FilesToDownload, myFile)
	}

	return nil
}

func (a *App) DownloadFile(file *drive.File, path string) {

	err := a.GDFS.DownloadFile(file, path)
	if err != nil {
		a.WriteOutput(err.Error())
	}
}
