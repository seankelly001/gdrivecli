package app

import (
	"fmt"
	"gdrivecli/pkg/config"
	"gdrivecli/pkg/errors"
	fs "gdrivecli/pkg/file_system"
	"gdrivecli/pkg/myfile"
	"gdrivecli/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

	partitionNodes, err := a.FS.GetPartitionNodes()
	if err != nil {
		return nil, err
	}
	for _, node := range partitionNodes {
		root.AddChild(node)
	}

	return tree, nil
}

func (a *App) FSNodeSelectedFunc(node *tview.TreeNode) {

	nodeRef, ok := node.GetReference().(fs.FileReference)
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
			err := a.FS.SetFSChildren(node)
			if err != nil {
				a.WriteOutput(err.Error())
			}
		}
		node.SetExpanded(!node.IsExpanded())
	} else {
		err = a.ToggleUpload(node, false, false)
		if err != nil {
			a.WriteOutput(err.Error())
		}
	}
}

// Toggle a file(node) to upload
// If force is true, it will use the upload param
// Otherwise, download is ignored
func (a *App) ToggleUpload(node *tview.TreeNode, force, upload bool) error {
	nodeRef, ok := node.GetReference().(fs.FileReference)
	if !ok {
		return errors.ErrCasting
	}

	if force {
		nodeRef.Upload = upload
	} else {
		nodeRef.Upload = !nodeRef.Upload
	}

	mf, ok := a.FilesToUpload[nodeRef.FilePath]
	if ok {
		if mf.InProgress || mf.Done {
			return nil
		}
	}

	if nodeRef.Upload {
		myFile, err := myfile.NewUploadFile(nodeRef.FilePath, node)
		if err != nil {
			return err
		}
		a.FilesToUpload[nodeRef.FilePath] = myFile

		node.SetColor(config.TREE_DOWNLOAD_COLOUR)
	} else {
		delete(a.FilesToUpload, nodeRef.FilePath)
		node.SetColor(config.TREE_FILE_COLOUR)
	}
	node.SetReference(nodeRef)
	return nil
}

//Determine what to do when a node has input from keyboard
func (a *App) FSInputFunc(event *tcell.EventKey) *tcell.EventKey {

	//TODO error handling
	//TODO select all prompt
	cur := a.FSTree.GetCurrentNode()
	switch event.Key() {
	case tcell.KeyInsert:
		a.CreateDirPrompt(cur)
	case tcell.KeyDelete:
		a.DeletePrompt(cur)
	case tcell.KeyCtrlS:
		a.DownloadPrompt(cur)
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
			err := a.FS.CreateDir(newDirPath)
			msg := ""
			if err != nil {
				msg = fmt.Sprintf("Error creating directory: %s", err.Error())
			} else {
				msg = "Created directory: " + path
			}
			a.WriteOutput(msg)
			cur := a.FSTree.GetCurrentNode()
			a.FS.SetFSChildren(cur)
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
	isDir, err := a.FS.IsDir(path)
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
				err := a.FS.Delete(path)
				if err != nil {
					msg := fmt.Sprintf("error deleting directory: %s", err.Error())
					a.WriteOutput(msg)
				} else {
					msg := fmt.Sprintf("Deleted directory %s", path)
					a.WriteOutput(msg)
					//reload parent children, and set current node to parent
					parent := nodeReference.Parent
					a.FS.SetFSChildren(parent)
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

func (a *App) DownloadPrompt(node *tview.TreeNode) error {
	reference := node.GetReference()
	nodeReference, ok := reference.(fs.FileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	path := nodeReference.FilePath
	totalDownloadSize := a.GetTotalDownloadSize()
	diskName := filepath.VolumeName(path)
	freeDiskSpace, err := utils.GetDiskUsage(diskName)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Downloading %d files. Total size: %s. Free space on disk %s %s", len(a.FilesToDownload), totalDownloadSize, diskName, freeDiskSpace)
	a.WriteOutput(msg)

	for _, mf := range a.FilesToDownload {

		fullPath := filepath.Join(path, mf.Name)
		mf.Path = fullPath
		go a.DownloadFile(mf)
	}

	return nil
}

func (a *App) DownloadFile(mf *myfile.MyFile) {

	//delete(a.FilesToDownload, mf.Name)
	mf.InProgress = true
	mf.Node.SetColor(config.TREE_IN_PROGRESS_COLOUR)
	err := a.GDFS.DownloadFile(mf.GFile, mf.Path)
	if err != nil {
		a.WriteOutput(err.Error())
	}
	// mf.InProgress = false
	// mf.Done = true
}
