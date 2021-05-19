package main

import (
	"fmt"
	"gdrivecli/pkg/gcli"
	"gdrivecli/pkg/token"
	"gdrivecli/pkg/utils"
	"io/ioutil"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rivo/tview"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"

	"path/filepath"
)

func main() {

	var err error
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := token.GetClient(config)

	srv, err := drive.New(client)

	fmt.Printf("%+v", srv)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	app := tview.NewApplication()

	cli := gcli.GDriveCLI{
		Service:         srv,
		FilesToDownload: make(map[string]*drive.File, 0),
		App:             app,
	}

	tree := cli.GenerateTree()
	cli.Tree = tree

	err = cli.Display()
	if err != nil {
		log.Fatalf("Unable to start application: %v", err)
	}

	if len(cli.FilesToDownload) == 0 {
		fmt.Println("no files to download")
		os.Exit(0)
	}

	for _, file := range cli.FilesToDownload {
		if file.Size == 0 {
			fmt.Printf("file is empty, will not download: %s\n", file.Name)
			delete(cli.FilesToDownload, file.Id)
		}
	}

	var totalSize int64
	fmt.Printf("List of files to download:\n")
	for _, file := range cli.FilesToDownload {
		fmt.Printf("|%-60v|%-20v\n", file.Name, utils.ByteCountIEC(file.Size))
		totalSize += file.Size
	}

	totalDownloadSize := utils.ByteCountIEC(totalSize)
	continueDL := false
	rootFolder := ""
	for !continueDL {
		rootFolder, err = utils.ChooseFolder()
		if err != nil {
			log.Fatalf("Error choosing folder: %w", err)
		}
		diskName := filepath.VolumeName(rootFolder)
		freeDiskSpace, err := utils.GetDiskUsage(diskName)
		if err != nil {
			log.Fatalf("Error getting disk space: %w", err)
		}

		msg := fmt.Sprintf(`Are you sure you want to download files to %s?
Total Download Size: %s
Free Space On Drive: %s`, rootFolder, totalDownloadSize, freeDiskSpace)
		promptUse := &survey.Confirm{
			Message: msg,
		}
		if err := survey.AskOne(promptUse, &continueDL); err != nil {
			log.Fatalf("Error prompting: %v", err)
		}
	}

	fmt.Printf("Saving files to %s\n", rootFolder)

	for fileName := range cli.FilesToDownload {

		fmt.Printf("Lets download: %s\n", fileName)
		cli.DownloadFile(rootFolder, fileName)
	}
}
