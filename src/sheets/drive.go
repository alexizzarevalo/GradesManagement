package sheets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/alexizzarevalo/grades_management/src/email"
	"github.com/alexizzarevalo/grades_management/src/msg"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getDriveService(credentials string) *drive.Service {
	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(getHttpClient(credentials)))

	if err != nil {
		msg.Error(errors.New("No se pudo recuperar el cliente de Drive. " + err.Error()))
	}

	return srv
}

func ListFiles(srv *drive.Service) {
	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		msg.Error(errors.New("No se pudo recuperar los archivos. " + err.Error()))
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Println(len(i.ExportLinks))
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}

func Export(srv *drive.Service, spreadsheetId, name string) error {
	resp, err := srv.Files.Export(spreadsheetId, "application/pdf").Download()
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func Pdf(srv *drive.Service, srvSheets *sheets.Service, newSpreadsheet NewSpreadsheet, wg *sync.WaitGroup) {
	carnet := newSpreadsheet.Name
	pdfName := carnet + ".pdf"
	err := DeleteSheet(srvSheets, newSpreadsheet.SpreadsheetId, 0)
	if err != nil {
		msg.ErrorWithoutExit(errors.New("No se pudo eliminar la hoja del spreadsheet. " + err.Error()))
		wg.Done()
		return
	}
	err = Export(srv, newSpreadsheet.SpreadsheetId, pdfName)
	if err != nil {
		msg.ErrorWithoutExit(errors.New("No se pudo exportar el archivo de Google Sheet " + pdfName + " " + err.Error()))
		wg.Done()
		return
	}
	msg.Success("Pdf generado: " + pdfName)
	wg.Done()
}

func DeleteSpreadsheet(srv *drive.Service, newSpreadsheet NewSpreadsheet, wg *sync.WaitGroup) {
	err := srv.Files.Delete(newSpreadsheet.SpreadsheetId).Do()
	if err != nil {
		msg.ErrorWithoutExit(errors.New("No se pudo eliminar el archivo de Google Drive " + newSpreadsheet.Name + ".pdf " + err.Error()))
	}
	wg.Done()
}

func ExportSheetsInPDF(opt SheetsOptions) {
	srv := getDriveService(opt.Credentials)
	srvSheets := getSheetService(opt.Credentials)

	newSpreadsheets, err := CopySheetsIntoSeparateSpreadSheets(srvSheets, opt.Id)
	if err != nil {
		msg.Error(errors.New("No se pudo copiar las hojas a un nuevo spreadsheet. " + err.Error()))
	}

	var wg sync.WaitGroup

	// Se exporta a PDF cada Spreadsheet
	wg.Add(len(newSpreadsheets))
	for _, newSpreadsheet := range newSpreadsheets {
		go Pdf(srv, srvSheets, newSpreadsheet, &wg)
	}
	wg.Wait()

	// Se eliminan las spreadsheets creadas
	wg.Add(len(newSpreadsheets))
	for _, newSpreadsheet := range newSpreadsheets {
		go DeleteSpreadsheet(srv, newSpreadsheet, &wg)
	}
	wg.Wait()
}

func ExportSheetsInPDFAndSendEmail(opt SheetsOptions, emailOpt email.EmailOptions) {
	ExportSheetsInPDF(opt)
	email.EmailOnly(emailOpt)
}
