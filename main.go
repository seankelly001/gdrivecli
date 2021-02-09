package main

import (
	"encoding/json"
	"fmt"
	"gdrivecli/pkg/gcli"
	"gdrivecli/pkg/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

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
	client := getClient(config)

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
		//fmt.Printf("%s", file.Name, utils.ByteCountIEC(file.Size))

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

	fmt.Println("Saving files to C:/Users/Sean/Documents/gdrivecli")

	for _, file := range cli.FilesToDownload {
		fmt.Printf("Downloading %s...\n", file.Name)
		path := "E:/Users/Sean/Stuff/Media/TV Shows/Taskmaster/" + file.Name
		f, err := os.Create(path)
		defer f.Close()
		if err != nil {
			log.Fatalf("Unable to create file: %v", err)
		}

		resp, err := srv.Files.Get(file.Id).Download()
		if err != nil {
			log.Fatalf("Unable to download file: %v", err)
		}
		defer resp.Body.Close()
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			log.Fatalf("Unable copy file: %v", err)
		}
	}
}
