package gcli

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func (cli GDriveCLI) GetFiles(parentID string) (*drive.FileList, error) {

	var r1 *drive.FileList
	var err error
	if parentID == "root" {
		r1, err = cli.Service.Files.List().PageSize(1000).Q(fmt.Sprintf("'%s' in parents", parentID)).Q(fmt.Sprintf("sharedWithMe")).OrderBy("folder,name").
			Fields("nextPageToken, files(id, name,mimeType,size)").Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve files: %v", err)
		}
	} else {
		r1, err = cli.Service.Files.List().PageSize(1000).Q(fmt.Sprintf("'%s' in parents", parentID)).OrderBy("folder,name").
			Fields("nextPageToken, files(id, name,mimeType,size)").Do()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve files: %v", err)
		}
	}

	if len(r1.Files) == 0 {
		return nil, fmt.Errorf("no files found")
	}
	return r1, nil
}
