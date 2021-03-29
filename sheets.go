package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type SheetsOptions struct {
	Id          string
	Credentials string
	Cells       Cells
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tokFile := getFilePath("token.json")
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
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
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

func getFolderPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	appFolder := path.Join(homeDir, ".grades_management")

	return appFolder
}

func getFilePath(filename string) string {
	appFolder := getFolderPath()
	filePath := path.Join(appFolder, filename)
	return filePath
}

func getCredentialBytes(credentials string) []byte {
	var path string
	if strings.Compare(credentials, "") != 0 {
		path = credentials
	} else {
		path = getFilePath("credentials.json")
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	return b
}

func getSheetService(credentials string) *sheets.Service {
	b := getCredentialBytes(credentials)
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return srv
}

func checkAppFolder() {
	appFolder := getFolderPath()
	_, err := os.ReadDir(appFolder)
	if err != nil {
		os.Mkdir(appFolder, 0777)
	}
}

func getGradesFromSpreadSheet(opt SheetsOptions) {
	checkAppFolder() // Si no existe el folder de la app, se crea

	// Se extrae el carnet y la nota de las celdas espeficiadas
	srv := getSheetService(opt.Credentials)
	spreadsheetId := opt.Id

	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Fatalf("Error al intentar obtener la hoja de calculo. %v", err)
	}

	ranges := []string{}

	for _, sheet := range spreadsheet.Sheets {
		ranges = append(ranges, sheet.Properties.Title+"!"+opt.Cells.Carne)
		ranges = append(ranges, sheet.Properties.Title+"!"+opt.Cells.Grade)
	}

	values, err := srv.Spreadsheets.Values.BatchGet(spreadsheetId).Ranges(ranges...).Do()
	if err != nil {
		log.Fatalf("Error al intentar obtener los valores. %v", err)
	}

	fmt.Println("Carne,Nota")
	for i := 0; i < len(values.ValueRanges); i += 2 {
		carne := values.ValueRanges[i].Values[0][0]
		grade := values.ValueRanges[i+1].Values[0][0]
		fmt.Printf("%v,%v\n", carne, grade)
	}
}
