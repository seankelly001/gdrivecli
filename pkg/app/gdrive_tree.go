package app

import (
	"fmt"
	"gdrivecli/pkg/config"
	"gdrivecli/pkg/errors"
	"gdrivecli/pkg/myfile"
	"gdrivecli/pkg/utils"
	"strings"

	"gdrivecli/pkg/gdfs"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/drive/v3"
)

func NewGDFSTree(a *App) *tview.TreeView {

	root := tview.NewTreeNode("root")
	rootFile := &drive.File{
		MimeType: "application/vnd.google-apps.folder",
	}
	rootReference := gdfs.GDFileReference{
		File:     rootFile,
		Download: false,
		Parent:   nil,
		Virtual:  true,
		Shared:   false,
		OrderBy:  "name",
	}
	root.SetReference(rootReference)
	root.SetExpanded(true)
	root.SetColor(config.TREE_DIR_COLOUR)

	mine := tview.NewTreeNode("My Drive").SetColor(config.TREE_DIR_COLOUR)
	mineFile := &drive.File{
		Id:       "root",
		MimeType: "application/vnd.google-apps.folder",
	}
	mineReference := gdfs.GDFileReference{
		File:     mineFile,
		Download: false,
		Parent:   root,
		Virtual:  false,
		Shared:   false,
		OrderBy:  "name",
	}
	mine.SetReference(mineReference)
	mine.SetExpanded(false)

	shared := tview.NewTreeNode("Shared With Me").SetColor(config.TREE_DIR_COLOUR)
	sharedFile := &drive.File{
		Id:       "root",
		MimeType: "application/vnd.google-apps.folder",
	}
	sharedReference := gdfs.GDFileReference{
		File:     sharedFile,
		Download: false,
		Parent:   root,
		Virtual:  false,
		Shared:   true,
		OrderBy:  "name",
	}
	shared.SetReference(sharedReference)
	shared.SetExpanded(false)

	root.AddChild(mine)
	root.AddChild(shared)

	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetSelectedFunc(a.GDFSNodeSelectedFunc)
	tree.SetInputCapture(a.GDFSTreeInputFunc)
	tree.SetTitle("Google Drive")
	tree.SetTitleColor(config.TITLE_COLOUR)
	tree.SetGraphicsColor(config.TREE_GRAPHICS_COLOUR)
	tree.SetBorderColor(config.BORDER_COLOUR)
	tree.SetBorder(true)
	tree.SetBorderPadding(1, 1, 1, 1)

	return tree
}

func (a *App) GDFSNodeSelectedFunc(node *tview.TreeNode) {

	//Do this to signal that work is being done
	originalNodeText := node.GetText()
	node.SetText(originalNodeText + "...")
	a.tvApp.ForceDraw()

	nodeRef, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		a.WriteOutput("error casting")
		return
	}
	if utils.IsGDFolder(nodeRef.File) {
		if !nodeRef.Virtual {
			if len(node.GetChildren()) == 0 {
				err := a.GDFS.SetChildren(node)
				if err != nil {
					a.WriteOutput(fmt.Sprintf("error setting children: %s", err.Error()))
				}
			}
		}
		node.SetExpanded(!node.IsExpanded())
	} else {
		a.ToggleDownload(node, false, false)
	}
	node.SetText(originalNodeText)
}

//Toggle a file(node) to download
//If force is true, it will use the download param
//Otherwise, download is ignored
func (a *App) ToggleDownload(node *tview.TreeNode, force, download bool) {
	nodeRef, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		//TODO
	}

	if force {
		nodeRef.Download = download
	} else {
		nodeRef.Download = !nodeRef.Download
	}

	mf, ok := a.FilesToDownload[nodeRef.File.Id]
	if ok {
		if mf.InProgress || mf.Done {
			return
		}
	}

	if nodeRef.Download {
		myFile := myfile.NewDownloadFile(nodeRef.File, node)
		a.FilesToDownload[nodeRef.File.Id] = myFile

		node.SetColor(config.TREE_DOWNLOAD_COLOUR)
	} else {
		delete(a.FilesToDownload, nodeRef.File.Id)
		node.SetColor(config.TREE_FILE_COLOUR)
	}
	node.SetReference(nodeRef)
}

//For jumping to char
func (a *App) GDFSTreeInputFunc(event *tcell.EventKey) *tcell.EventKey {

	//Ctrl + S: Upload (tbd??)
	//Ctrl + N: Sort by name (reversible)
	//Ctrl + D: Sort by date (reversible)
	//Ctrl + A: Select All
	//Ctrl + U: Deselct All
	var err error
	cur := a.GDFSTree.GetCurrentNode()

	if event.Rune() > 65 && event.Rune() < 122 {
		//Jump to char
		keyChar := string(event.Rune())
		err = a.JumpToNodeKeyPrompt(cur, keyChar)
		if err != nil {
			a.WriteOutput(fmt.Sprintf("error switching: %s", err.Error()))
		}
		return nil
	}

	switch event.Key() {
	case tcell.KeyCtrlS:
		err = a.UploadPrompt(cur)
	case tcell.KeyCtrlN:
		err = a.SwitchOutputPrompt(cur, "name")
	case tcell.KeyCtrlD:
		err = a.SwitchOutputPrompt(cur, "modifiedTime")
	case tcell.KeyCtrlA:
		err = a.SelectAllPrompt(cur, true)
	case tcell.KeyCtrlU:
		err = a.SelectAllPrompt(cur, false)
	case tcell.KeyInsert:
		err = a.CreateGDFolderPrompt(cur)
	case tcell.KeyDelete:
		err = a.DeleteGDPrompt(cur)
	case tcell.KeyPgUp:
		err = a.JumpToNodePos(cur, -15)
	case tcell.KeyPgDn:
		err = a.JumpToNodePos(cur, 15)
	case tcell.KeyBackspace:
		err = a.JumpParentPrompt(cur)
	default:
		return event
	}
	if err != nil {
		a.WriteOutput(fmt.Sprintf("error switching: %s", err.Error()))
	}

	return nil
}

func (a *App) UploadPrompt(node *tview.TreeNode) error {
	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	totalUploadSize, err := a.GetTotalUploadSize()
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("Uploading %d files. Total size: %s", len(a.FilesToUpload), totalUploadSize)
	a.WriteOutput(msg)

	for _, mf := range a.FilesToUpload {
		go a.UploadFile(ref.File.Id, mf)
	}

	return nil
}

func (a *App) SwitchOutputPrompt(node *tview.TreeNode, orderBy string) error {

	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		// TODO
		return nil
	}
	parent := ref.Parent
	parentRef := parent.GetReference().(gdfs.GDFileReference)
	currentOrderBy := parentRef.OrderBy
	var newOrderBy string

	a.WriteOutput(fmt.Sprintf("ob: %s, cob: %s", orderBy, currentOrderBy))

	switch {
	case orderBy == "name" && currentOrderBy == "name":
		newOrderBy = "name desc"
	case orderBy == "modifiedTime" && currentOrderBy == "modifiedTime":
		newOrderBy = "modifiedTime desc"
	default:
		newOrderBy = orderBy
	}

	parentRef.OrderBy = newOrderBy
	parent.SetReference(parentRef)

	err := a.GDFS.SetChildren(parent)
	if err != nil {
		return err
	}
	parent.SetExpanded(true)
	a.GDFSTree.SetCurrentNode(parent.GetChildren()[0])

	return nil
}

func (a *App) CreateGDFolderPrompt(node *tview.TreeNode) error {

	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		return errors.ErrCasting
	}
	file := ref.File
	if !utils.IsGDFolder(file) {
		//TODO - get parent??
		return nil
	}
	prompt := "Enter dir name: "
	doneFunc := func(key tcell.Key) {
		if key == tcell.KeyEnter {
			//Create folder with name from input, write message, reset nodes children and expand node
			newFolderName := a.GetInputText()
			var msg string = ""
			err := a.GDFS.CreateFolder(newFolderName, file.Id)
			if err != nil {
				msg = fmt.Sprintf("Error creating directory: %s", err.Error())
			} else {
				msg = "Created directory: " + newFolderName
			}
			a.WriteOutput(msg)
			a.GDFS.SetChildren(node)
			node.Expand()
		}
		//Reset
		a.ResetInput()
		a.tvApp.SetFocus(a.GDFSTree)
	}
	a.Prompt(prompt, doneFunc)
	return nil
}

func (a *App) DeleteGDPrompt(node *tview.TreeNode) error {
	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		return errors.ErrCasting
	}
	isDir := utils.IsGDFolder(ref.File)
	name := ref.File.Name
	var prompt string
	if isDir {
		prompt = fmt.Sprintf("Delete gd directory %s ? (y/n): ", name)
	} else {
		prompt = fmt.Sprintf("Delete gd file %s ? (y/n): ", name)
	}

	doneFunc := func(key tcell.Key) {
		if key == tcell.KeyEnter {
			//Create directory with name from input, write message, reset nodes children and expand node
			answer := strings.ToLower(a.GetInputText())
			if answer == "y" {
				//delete dir
				err := a.GDFS.DeleteFile(ref.File)
				if err != nil {
					msg := fmt.Sprintf("error deleting gd file: %s", err.Error())
					a.WriteOutput(msg)
				} else {
					msg := fmt.Sprintf("Deleted gd file %s", ref.File.Name)
					a.WriteOutput(msg)
					//reload parent children, and set current node to parent
					parent := ref.Parent
					a.GDFS.SetChildren(parent)
					a.GDFSTree.SetCurrentNode(parent)
				}
				//Reset
				a.ResetInput()
				a.tvApp.SetFocus(a.GDFSTree)
			} else if answer == "n" {
				a.ResetInput()
				a.tvApp.SetFocus(a.GDFSTree)
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

func (a *App) SelectAllPrompt(node *tview.TreeNode, download bool) error {

	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		// TODO
		return nil
	}
	parent := ref.Parent

	children := parent.GetChildren()
	for _, c := range children {
		childRef, ok := c.GetReference().(gdfs.GDFileReference)
		if !ok {
			// TODO
			return nil
		}
		if !utils.IsGDFolder(childRef.File) {
			a.ToggleDownload(c, true, download)
		}
	}

	return nil
}

func (a *App) JumpToNodeKeyPrompt(node *tview.TreeNode, keyChar string) error {

	var jumpToNode *tview.TreeNode
	var siblings []*tview.TreeNode

	if node == nil {
		return nil
	}

	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		return errors.ErrCasting
	}
	parentNode := ref.Parent
	if parentNode != nil {
		siblings = parentNode.GetChildren()
	} else {
		return nil
	}
	jumpToNode = utils.FilterNodeChildren(siblings, keyChar)
	a.GDFSTree.SetCurrentNode(jumpToNode)
	return nil
}

func (a *App) JumpToNodePos(node *tview.TreeNode, pos int) error {

	var jumpToNode *tview.TreeNode
	var siblings []*tview.TreeNode

	if node == nil {
		return nil
	}

	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		return errors.ErrCasting
	}
	parentNode := ref.Parent
	if parentNode != nil {
		siblings = parentNode.GetChildren()
	} else {
		return nil
	}

	jumpToNode = utils.JumpNodePosition(node, siblings, pos)
	a.GDFSTree.SetCurrentNode(jumpToNode)
	return nil
}

func (a *App) JumpParentPrompt(node *tview.TreeNode) error {

	if node == nil {
		return nil
	}

	ref, ok := node.GetReference().(gdfs.GDFileReference)
	if !ok {
		return errors.ErrCasting
	}
	parentNode := ref.Parent
	if parentNode != nil {
		a.GDFSTree.SetCurrentNode(parentNode)
	}
	return nil
}

func (a *App) UploadFile(parentID string, mf *myfile.MyFile) {

	mf.InProgress = true
	mf.Node.SetColor(config.TREE_IN_PROGRESS_COLOUR)

	gFile := &drive.File{
		Name:    mf.Name,
		Parents: []string{parentID},
	}
	mf.GFile = gFile
	err := a.GDFS.UploadFile(mf)
	if err != nil {
		a.WriteOutput(err.Error())
	}
}
