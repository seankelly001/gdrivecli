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
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func main() {
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

	cli := gcli.GDriveCLI{
		Service:         srv,
		FilesToDownload: make(map[string]*drive.File, 0),
	}

	cli.Display("root")

	var totalSize int64
	fmt.Printf("List of files to download:\n")
	for _, file := range cli.FilesToDownload {
		fmt.Printf("|%-30v|%-20v\n", file.Name, utils.ByteCountIEC(file.Size))
		totalSize += file.Size
	}
	fmt.Printf("Total size of files to download: %s\n", utils.ByteCountIEC(totalSize))
	continueDL := false
	promptUse := &survey.Confirm{
		Message: "Do you want to continue?",
	}
	if err := survey.AskOne(promptUse, &continueDL); err != nil {
		log.Fatalf("Error prompting: %v", err)
	}
	if !continueDL {
		fmt.Println("Exiting...")
		os.Exit(0)
	}

	// TODO - choose folder
	rootPath := "E:/Users/Sean/Stuff/Media/TV Shows/Taskmaster/"
	fmt.Printf("Saving files to %s\n", rootPath)

	for _, file := range cli.FilesToDownload {
		cli.DownloadFile(rootPath, file.Name)
	}
	// 	fmt.Printf("Downloading %s...\n", file.Name)
	// 	path := rootPath + file.Name
	// 	f, err := os.Create(path)
	// 	defer f.Close()
	// 	if err != nil {
	// 		log.Fatalf("Unable to create file: %v", err)
	// 	}

	// 	waitChan := make(chan struct{})
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	go utils.DisplayProgress(ctx, waitChan, f, file.Size)

	// 	resp, err := srv.Files.Get(file.Id).Download()
	// 	if err != nil {
	// 		log.Fatalf("Unable to download file: %v", err)
	// 	}
	// 	defer resp.Body.Close()
	// 	_, err = io.Copy(f, resp.Body)
	// 	if err != nil {
	// 		log.Fatalf("Unable copy file: %v", err)
	// 	}

	// 	cancel()
	// 	<-waitChan
	// 	// need to wait here
	// }
}
