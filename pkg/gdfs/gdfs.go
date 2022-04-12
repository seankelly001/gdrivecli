package gdfs

import (
	"fmt"
	"gdrivecli/pkg/config"
	"gdrivecli/pkg/utils"

	"github.com/rivo/tview"
	drive "google.golang.org/api/drive/v3"
)

type GDFileSystem struct {
	Service *drive.Service
	//FilesToDownload map[string]*drive.File
	//FilesToDownload map[string]*myfile.MyFile

	//SizeBytes       int
}

type GDFileReference struct {
	File     *drive.File
	Download bool
	Parent   *tview.TreeNode
	Virtual  bool
	Shared   bool
	OrderBy  string
}

func NewGDFS() (*GDFileSystem, error) {
	srv, err := GetGDriveService()
	if err != nil {
		return nil, fmt.Errorf("error getting gdrive service: %w", err)
	}
	gdfs := &GDFileSystem{
		Service: srv,
		//ilesToDownload: make(map[string]*myfile.MyFile, 0),
	}
	return gdfs, nil
}

func (gdfs *GDFileSystem) SetChildren(node *tview.TreeNode) error {

	node.ClearChildren()
	reference := node.GetReference()
	nodeReference, ok := reference.(GDFileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	files, err := gdfs.GetFiles(nodeReference.File.Id, nodeReference.OrderBy, nodeReference.Shared)
	if err != nil {
		return err
	}
	for _, file := range files.Files {
		childReference := GDFileReference{
			File:     file,
			Download: false,
			Parent:   node,
			Virtual:  false,
			Shared:   false,
			OrderBy:  "name",
		}
		child := tview.NewTreeNode(file.Name)
		child.SetReference(childReference)
		child.SetExpanded(false)

		if utils.IsGDFolder(file) {
			child.SetColor(config.TREE_DIR_COLOUR)
		} else {
			child.SetColor(config.TREE_FILE_COLOUR)
		}
		node.AddChild(child)
	}
	return nil
}
