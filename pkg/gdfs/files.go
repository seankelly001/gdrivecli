package gdfs

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func (gdfs *GDFileSystem) GetFiles(parentID string, shared bool) (*drive.FileList, error) {

	var r1 *drive.FileList
	var err error
	//if parentID == "root" {
	if shared {
		r1, err = gdfs.Service.Files.List().PageSize(1000).Q(fmt.Sprintf("'%s' in parents", parentID)).Q(fmt.Sprintf("sharedWithMe")).OrderBy("name").
			Fields("nextPageToken, files(id,name,mimeType,size)").Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve files: %v", err)
		}
	} else {
		r1, err = gdfs.Service.Files.List().PageSize(1000).Q(fmt.Sprintf("'%s' in parents", parentID)).OrderBy("name").
			Fields("nextPageToken, files(id,name,mimeType,size)").Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve files: %v", err)
		}
	}

	if len(r1.Files) == 0 {
		return nil, fmt.Errorf("no files found")
	}
	return r1, nil
}

func (gdfs *GDFileSystem) GetTotalDownloadSizeBytes() int64 {

	var totalSize int64 = 0
	for _, f := range gdfs.FilesToDownload {
		totalSize += f.Size
	}
	return totalSize
}
