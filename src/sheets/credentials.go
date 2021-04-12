package sheets

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexizzarevalo/grades_management/src/msg"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

var scopes = []string{
	sheets.SpreadsheetsScope,
	drive.DriveFileScope,
}

func getFolderPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		msg.Error(err)
	}

	appFolder := filepath.Join(homeDir, ".grades_management")

	return appFolder
}

func getFilePath(filename string) string {
	appFolder := getFolderPath()
	filePath := filepath.Join(appFolder, filename)
	return filePath
}

func checkAppFolder() {
	appFolder := getFolderPath()
	_, err := os.ReadDir(appFolder)
	if err != nil {
		os.Mkdir(appFolder, 0777)
	}
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
		msg.ErrorWithoutExit(errors.New("No se pudo leer las credenciales: " + err.Error()))
		msg.Info("Solicite al desarrollador las credenciales o genere su propia credencial siguiendo esta guia:")
		msg.Info("https://github.com/alexizzarevalo/GradesManagement#credenciales-de-google-cloud-project")
		os.Exit(1)
	}
	return b
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
	msg.Info("Vaya al siguiente enlace en su navegador y luego pegue el código de autorización: ")
	fmt.Printf("\n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		msg.Error(errors.New("No se pudo leer el codigo de autorizacion: " + err.Error()))
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		msg.Error(errors.New("No se pudo obtener el token desde la web: " + err.Error()))
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
	msg.Info("Guardando el token oauth en: " + path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		msg.Error(errors.New("No se pudo guardar el token oauth: " + err.Error()))
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getHttpClient(credentials string) *http.Client {
	b := getCredentialBytes(credentials)
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		msg.Error(errors.New("No se pudo analizar el archivo de credeciales. " + err.Error()))
	}
	client := getClient(config)

	return client
}
