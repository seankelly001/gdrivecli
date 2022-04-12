package gdfs

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func (gdfs *GDFileSystem) GetFiles(parentID, orderBy string, shared bool) (*drive.FileList, error) {

	var r1 *drive.FileList
	var err error
	if shared {
		r1, err = gdfs.Service.Files.List().PageSize(1000).Q(fmt.Sprintf("'%s' in parents", parentID)).Q(fmt.Sprintf("sharedWithMe")).OrderBy(orderBy).
			Fields("nextPageToken, files(id,name,mimeType,size,modifiedTime)").Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve files: %v", err)
		}
	} else {
		r1, err = gdfs.Service.Files.List().PageSize(1000).Q(fmt.Sprintf("'%s' in parents", parentID)).OrderBy(orderBy).
			Fields("nextPageToken, files(id,name,mimeType,size,modifiedTime)").Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve files: %v", err)
		}
	}

	if len(r1.Files) == 0 {
		return nil, fmt.Errorf("no files found")
	}
	return r1, nil
}

func (gdfs *GDFileSystem) CreateFolder(name, parentID string) error {

	folder := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}
	_, err := gdfs.Service.Files.Create(folder).Do()
	return err
}

func (gdfs *GDFileSystem) DeleteFile(file *drive.File) error {

	return gdfs.Service.Files.Delete(file.Id).Do()
}
