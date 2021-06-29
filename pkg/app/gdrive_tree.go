package app

import (
	"gdrivecli/pkg/config"
	"gdrivecli/pkg/utils"
	"log"

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
	originalNodeText := node.GetText()
	node.SetText(originalNodeText + "...")
	a.tvApp.ForceDraw()

	reference := node.GetReference()
	fileReference, ok := reference.(gdfs.GDFileReference)
	if !ok {
		a.WriteOutput("error casting")
		return
	}

	if fileReference.File.MimeType == "application/vnd.google-apps.folder" {
		if !fileReference.Virtual {
			if len(node.GetChildren()) == 0 {
				err := a.GDFS.SetChildren(node)
				if err != nil {
					log.Fatalf("err: %w", err)
				}
			}
		}
		node.SetExpanded(!node.IsExpanded())
	} else {
		fileReference.Download = !fileReference.Download
		if fileReference.Download {
			a.GDFS.FilesToDownload[fileReference.File.Id] = fileReference.File
			node.SetColor(config.TREE_DOWNLOAD_COLOUR)
		} else {
			delete(a.GDFS.FilesToDownload, fileReference.File.Id)
			node.SetColor(config.TREE_FILE_COLOUR)
		}

	}

	node.SetReference(fileReference)
	node.SetText(originalNodeText)
}

//For jumping to char
func (a *App) GDFSTreeInputFunc(event *tcell.EventKey) *tcell.EventKey {

	keyChar := string(event.Rune())
	if event.Rune() < 65 || event.Rune() > 122 {
		return event
	}
	currentNode := a.GDFSTree.GetCurrentNode()
	if currentNode != nil {

		reference := currentNode.GetReference()
		fileRererence, ok := reference.(gdfs.GDFileReference)
		if !ok {
			// TODO
		}
		parentNode := fileRererence.Parent
		if parentNode != nil {
			siblings := parentNode.GetChildren()
			if len(siblings) > 0 {
				jumpToNode := utils.FilterNodeChildren(siblings, keyChar)
				a.GDFSTree.SetCurrentNode(jumpToNode)
			}
		}
	}
	return event
}

// func (gdfs *GDFileSystem) SetChildren(node *tview.TreeNode) error {

// 	reference := node.GetReference()
// 	nodeReference, ok := reference.(FileReference)
// 	if !ok {
// 		return fmt.Errorf("error casting")
// 	}
// 	files, err := gdfs.GetFiles(nodeReference.File.Id, nodeReference.Shared)
// 	if err != nil {
// 		return err
// 	}
// 	for _, file := range files.Files {
// 		childReference := FileReference{
// 			File:     file,
// 			Download: false,
// 			Parent:   node,
// 			Virtual:  false,
// 			Shared:   false,
// 		}
// 		child := tview.NewTreeNode(file.Name)
// 		child.SetReference(childReference)
// 		child.SetExpanded(false)

// 		if file.MimeType == "application/vnd.google-apps.folder" {
// 			child.SetColor(tcell.ColorRed)
// 		}
// 		node.AddChild(child)
// 	}
// 	return nil
// }

// func (gdfs *GDFileSystem) NodeSelectedFunc(node *tview.TreeNode) {

// 	originalNodeText := node.GetText()
// 	node.SetText(originalNodeText + "...")
// 	//gdfs.App.ForceDraw()

// 	reference := node.GetReference()
// 	fileReference, ok := reference.(FileReference)
// 	if !ok {
// 		fmt.Println("no ok...")
// 		return
// 	}

// 	if fileReference.File.MimeType == "application/vnd.google-apps.folder" {

// 		if !fileReference.Virtual {
// 			if len(node.GetChildren()) == 0 {
// 				err := gdfs.SetChildren(node)
// 				if err != nil {
// 					log.Fatalf("err: %w", err)
// 				}
// 			}
// 		}
// 		node.SetExpanded(!node.IsExpanded())
// 	} else {
// 		fileReference.Download = !fileReference.Download
// 		if fileReference.Download {
// 			gdfs.FilesToDownload[fileReference.File.Id] = fileReference.File
// 			node.SetColor(tcell.ColorTeal)
// 		} else {
// 			delete(gdfs.FilesToDownload, fileReference.File.Id)
// 			node.SetColor(tcell.ColorWhite)
// 		}

// 	}

// 	node.SetReference(fileReference)
// 	node.SetText(originalNodeText)
// }

// func (gdfs *GDFileSystem) TreeInputFunc(event *tcell.EventKey) *tcell.EventKey {

// 	//event.Key()
// 	if event.Key() == tcell.KeyBackspace {
// 		gdfs.App.Stop()
// 	}
// 	// keyChar := ""
// 	// switch event.Key() {
// 	// case tcell.KeyCtrlA:
// 	// 	keyChar = "a"
// 	// case tcell.KeyCtrlB:
// 	// 	keyChar = "b"
// 	// }
// 	keyChar := string(event.Rune())
// 	//fmt.Printf("key: %s, rune: %d\n", keyString, event.Rune())
// 	if event.Rune() < 65 || event.Rune() > 122 {
// 		return event
// 	}
// 	currentNode := gdfs.Tree.GetCurrentNode()
// 	if currentNode != nil {
// 		//fmt.Printf("current node: %s\n", currentNode.GetText())

// 		reference := currentNode.GetReference()
// 		fileRererence, ok := reference.(FileReference)
// 		if !ok {
// 			// TODO
// 		}
// 		parentNode := fileRererence.Parent
// 		if parentNode != nil {
// 			//fmt.Printf("parent node: %s\n", parentNode.GetText())

// 			siblings := parentNode.GetChildren()
// 			if len(siblings) > 0 {
// 				jumpToNode := utils.FilterNodeChildren(siblings, keyChar)
// 				gdfs.Tree.SetCurrentNode(jumpToNode)
// 			}
// 		}

// 	}
// 	return event
// }

// func (cli *GDriveCLI) TreeMouseFunc (action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {

// 	if action != tview.MouseLeftClick {
// 		return
// 	}

// 	action.
// 	currentNode := cli.Tree.GetCurrentNode()
// 	if currentNode != nil {
// 		//fmt.Printf("current node: %s\n", currentNode.GetText())

// 		reference := currentNode.GetReference()
// 		fileRererence, ok := reference.(FileReference)
// 		if !ok {
// 			// TODO
// 		}
// 		parentNode := fileRererence.Parent
// 		if parentNode != nil {
// 			//fmt.Printf("parent node: %s\n", parentNode.GetText())

// 			siblings := parentNode.GetChildren()
// 			if len(siblings) > 0 {
// 				jumpToNode := utils.FilterNodeChildren(siblings, keyString)
// 				cli.Tree.SetCurrentNode(jumpToNode)
// 			}
// 		}

// 	}
// 	return event
// }
