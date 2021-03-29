package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type SheetsOptions struct {
	Id    string
	Cells Cells
}

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

func getSheetService() *sheets.Service {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

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

func getGradesFromSpreadSheet(opt SheetsOptions) {
	// Se extrae el carnet y la nota de las celdas espeficiadas
	srv := getSheetService()
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
