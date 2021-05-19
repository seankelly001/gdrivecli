package gcli

import (
	"fmt"
	"gdrivecli/pkg/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"google.golang.org/api/drive/v3"
)

type GDriveCLI struct {
	Service         *drive.Service
	FilesToDownload map[string]*drive.File
	SizeBytes       int
	App             *tview.Application
	Tree            *tview.TreeView
}

type FileReference struct {
	File     *drive.File
	Download bool
	Parent   *tview.TreeNode
	Virtual  bool
	Shared   bool
}

func (cli GDriveCLI) Display() error {

	if err := cli.App.SetRoot(cli.Tree, true).Run(); err != nil {
		return err
	}
	return nil
}

func (cli *GDriveCLI) GenerateTree() *tview.TreeView {

	root := tview.NewTreeNode("root")
	rootFile := &drive.File{
		MimeType: "application/vnd.google-apps.folder",
	}
	rootReference := FileReference{
		File:     rootFile,
		Download: false,
		Parent:   nil,
		Virtual:  true,
		Shared:   false,
	}
	root.SetReference(rootReference)
	root.SetExpanded(false)
	root.SetColor(tcell.ColorRed)

	mine := tview.NewTreeNode("mine")
	mineFile := &drive.File{
		Id:       "root",
		MimeType: "application/vnd.google-apps.folder",
	}
	mineReference := FileReference{
		File:     mineFile,
		Download: false,
		Parent:   root,
		Virtual:  false,
		Shared:   false,
	}
	mine.SetReference(mineReference)
	mine.SetExpanded(false)

	shared := tview.NewTreeNode("shared")
	sharedFile := &drive.File{
		Id:       "root",
		MimeType: "application/vnd.google-apps.folder",
	}
	sharedReference := FileReference{
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
	tree.SetSelectedFunc(cli.NodeSelectedFunc)
	tree.SetInputCapture(cli.TreeInputFunc)
	tree.SetTitle("select files to download")
	tree.SetTitleColor(tcell.ColorOrange)
	tree.SetGraphicsColor(tcell.ColorSlateBlue)
	tree.SetBorderColor(tcell.ColorPink)
	tree.SetBorder(true)

	tree.SetBorderPadding(1, 1, 1, 1)
	return tree
}

func (cli *GDriveCLI) SetChildren(node *tview.TreeNode) error {

	reference := node.GetReference()
	nodeReference, ok := reference.(FileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	files, err := cli.GetFiles(nodeReference.File.Id, nodeReference.Shared)
	if err != nil {
		return err
	}
	for _, file := range files.Files {
		childReference := FileReference{
			File:     file,
			Download: false,
			Parent:   node,
			Virtual:  false,
			Shared:   false,
		}
		child := tview.NewTreeNode(file.Name)
		child.SetReference(childReference)
		child.SetExpanded(false)

		if file.MimeType == "application/vnd.google-apps.folder" {
			child.SetColor(tcell.ColorRed)
		}
		node.AddChild(child)
	}
	return nil
}

func (cli *GDriveCLI) NodeSelectedFunc(node *tview.TreeNode) {

	originalNodeText := node.GetText()
	node.SetText(originalNodeText + "...")
	//cli.App.
	cli.App.ForceDraw()

	//fmt.Println("node selected!")
	reference := node.GetReference()
	fileReference, ok := reference.(FileReference)
	if !ok {
		fmt.Println("no ok...")
		return
	}

	if fileReference.File.MimeType == "application/vnd.google-apps.folder" {

		if !fileReference.Virtual {
			if len(node.GetChildren()) == 0 {
				err := cli.SetChildren(node)
				if err != nil {
					// TODO
				}
			}
		}
		node.SetExpanded(!node.IsExpanded())
	} else {
		fileReference.Download = !fileReference.Download
		if fileReference.Download {
			cli.FilesToDownload[fileReference.File.Id] = fileReference.File
			node.SetColor(tcell.ColorTeal)
		} else {
			delete(cli.FilesToDownload, fileReference.File.Id)
			node.SetColor(tcell.ColorWhite)
		}

	}

	node.SetReference(fileReference)
	node.SetText(originalNodeText)
}

func (cli *GDriveCLI) TreeInputFunc(event *tcell.EventKey) *tcell.EventKey {

	//event.Key()
	if event.Key() == tcell.KeyBackspace {
		cli.App.Stop()
	}
	// keyChar := ""
	// switch event.Key() {
	// case tcell.KeyCtrlA:
	// 	keyChar = "a"
	// case tcell.KeyCtrlB:
	// 	keyChar = "b"
	// }
	keyChar := string(event.Rune())
	//fmt.Printf("key: %s, rune: %d\n", keyString, event.Rune())
	if event.Rune() < 65 || event.Rune() > 122 {
		return event
	}
	currentNode := cli.Tree.GetCurrentNode()
	if currentNode != nil {
		//fmt.Printf("current node: %s\n", currentNode.GetText())

		reference := currentNode.GetReference()
		fileRererence, ok := reference.(FileReference)
		if !ok {
			// TODO
		}
		parentNode := fileRererence.Parent
		if parentNode != nil {
			//fmt.Printf("parent node: %s\n", parentNode.GetText())

			siblings := parentNode.GetChildren()
			if len(siblings) > 0 {
				jumpToNode := utils.FilterNodeChildren(siblings, keyChar)
				cli.Tree.SetCurrentNode(jumpToNode)
			}
		}

	}
	return event
}

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
