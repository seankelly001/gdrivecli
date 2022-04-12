package file_system

import (
	"fmt"
	"gdrivecli/pkg/config"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/disk"
)

type FileSystem struct {
	FilesToUpload map[string]os.FileInfo
}

type FileReference struct {
	FileName string
	FilePath string
	Upload   bool
	Parent   *tview.TreeNode
}

func NewFS() *FileSystem {
	fs := &FileSystem{
		FilesToUpload: make(map[string]os.FileInfo, 0),
	}
	return fs
}

func (fs *FileSystem) GetPartitionNodes() ([]*tview.TreeNode, error) {

	var nodes []*tview.TreeNode
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}
	for _, partition := range partitions {
		partitionNode := tview.NewTreeNode(partition.Device)
		if err != nil {
			return nil, err
		}
		partitionRef := FileReference{
			FileName: partition.Device,
			FilePath: partition.Device + "\\",
		}
		partitionNode.SetExpanded(false)
		partitionNode.SetReference(partitionRef)
		partitionNode.SetColor(config.TREE_DIR_COLOUR)
		nodes = append(nodes, partitionNode)
	}
	return nodes, nil

}

func (fs *FileSystem) SetFSChildren(node *tview.TreeNode) error {
	node.ClearChildren()
	ref := node.GetReference()
	nodeRef, ok := ref.(FileReference)
	if !ok {
		return fmt.Errorf("error casting")
	}
	f, err := os.Open(nodeRef.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	files, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range files {
		childReference := FileReference{
			FileName: file.Name(),
			FilePath: filepath.Join(nodeRef.FilePath, file.Name()),
			Upload:   false,
			Parent:   node,
		}
		child := tview.NewTreeNode(file.Name())
		child.SetReference(childReference)
		child.SetExpanded(false)

		if file.IsDir() {
			child.SetColor(config.TREE_DIR_COLOUR)
		} else {
			child.SetColor(config.TREE_FILE_COLOUR)
		}
		node.AddChild(child)
	}
	return nil
}

func (fs *FileSystem) CreateDir(path string) error {

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileSystem) Delete(path string) error {
	return os.RemoveAll(path)
}

func (fs *FileSystem) IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func (fs *FileSystem) GetTotalUploadSizeBytes() int64 {

	var totalSize int64 = 0
	for _, fileInfo := range fs.FilesToUpload {
		totalSize += fileInfo.Size()
	}
	return totalSize
}
