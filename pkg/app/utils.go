package app

import (
	"gdrivecli/pkg/utils"
	"os"
)

func (a *App) GetTotalDownloadSize() string {

	var totalSize int64 = 0
	for _, f := range a.FilesToDownload {
		totalSize += f.GFile.Size
	}
	return utils.ByteCountIEC(totalSize)
}

func (a *App) GetTotalUploadSize() (string, error) {

	var totalSize int64 = 0
	for _, mf := range a.FilesToUpload {
		f, err := os.Open(mf.Path)
		if err != nil {
			return "", err
		}
		fStat, err := f.Stat()
		if err != nil {
			return "", err
		}
		totalSize += fStat.Size()
	}
	return utils.ByteCountIEC(totalSize), nil
}
