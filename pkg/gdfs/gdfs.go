package gdfs

import (
	"fmt"
	"gdrivecli/pkg/config"

	"github.com/rivo/tview"
	drive "google.golang.org/api/drive/v3"
)

type GDFileSystem struct {
	Service         *drive.Service
	FilesToDownload map[string]*drive.File
	//SizeBytes       int
}

type GDFileReference struct {
	File     *drive.File
	Download bool
	Parent   *tview.TreeNode
	Virtual  bool
	Shared   bool
}

func NewGDFS() (*GDFileSystem, error) {
	srv, err := GetGDriveService()
	if err != nil {
		return nil, fmt.Errorf("error getting gdrive service: %w", err)
	}
	gdfs := &GDFileSystem{
		Service:         srv,
		FilesToDownload: make(map[string]*drive.File, 0),
	}
	return gdfs, nil
}

func (gdfs *GDFileSystem) SetChildren(node *tview.TreeNode) error {

	reference := node.GetReference()
	nodeReference, ok := reference.(GDFileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	files, err := gdfs.GetFiles(nodeReference.File.Id, nodeReference.Shared)
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
		}
		child := tview.NewTreeNode(file.Name)
		child.SetReference(childReference)
		child.SetExpanded(false)

		if file.MimeType == "application/vnd.google-apps.folder" {
			child.SetColor(config.TREE_DIR_COLOUR)
		} else {
			child.SetColor(config.TREE_FILE_COLOUR)
		}
		node.AddChild(child)
	}
	return nil
}
