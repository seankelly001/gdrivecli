package gdfs

// import (
// 	"gdrivecli/pkg/utils"

// 	"github.com/gdamore/tcell/v2"
// 	"github.com/rivo/tview"
// 	"google.golang.org/api/drive/v3"
// )

// type FileReference struct {
// 	File     *drive.File
// 	Download bool
// 	Parent   *tview.TreeNode
// 	Virtual  bool
// 	Shared   bool
// }

// func (gdfs *GDFileSystem) GenerateTree() *tview.TreeView {

// 	root := tview.NewTreeNode("root")
// 	rootFile := &drive.File{
// 		MimeType: "application/vnd.google-apps.folder",
// 	}
// 	rootReference := FileReference{
// 		File:     rootFile,
// 		Download: false,
// 		Parent:   nil,
// 		Virtual:  true,
// 		Shared:   false,
// 	}
// 	root.SetReference(rootReference)
// 	root.SetExpanded(true)
// 	root.SetColor(tcell.ColorRed)

// 	mine := tview.NewTreeNode("mine")
// 	mineFile := &drive.File{
// 		Id:       "root",
// 		MimeType: "application/vnd.google-apps.folder",
// 	}
// 	mineReference := FileReference{
// 		File:     mineFile,
// 		Download: false,
// 		Parent:   root,
// 		Virtual:  false,
// 		Shared:   false,
// 	}
// 	mine.SetReference(mineReference)
// 	mine.SetExpanded(false)

// 	shared := tview.NewTreeNode("shared")
// 	sharedFile := &drive.File{
// 		Id:       "root",
// 		MimeType: "application/vnd.google-apps.folder",
// 	}
// 	sharedReference := FileReference{
// 		File:     sharedFile,
// 		Download: false,
// 		Parent:   root,
// 		Virtual:  false,
// 		Shared:   true,
// 	}
// 	shared.SetReference(sharedReference)
// 	shared.SetExpanded(false)

// 	root.AddChild(mine)
// 	root.AddChild(shared)

// 	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
// 	tree.SetSelectedFunc(gdfs.NodeSelectedFunc)
// 	tree.SetInputCapture(gdfs.TreeInputFunc)
// 	tree.SetTitle("select files to download")
// 	tree.SetTitleColor(tcell.ColorOrange)
// 	tree.SetGraphicsColor(tcell.ColorSlateBlue)
// 	tree.SetBorderColor(tcell.ColorPink)
// 	tree.SetBorder(true)

// 	tree.SetBorderPadding(1, 1, 1, 1)
// 	return tree
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
