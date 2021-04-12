package sheets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alexizzarevalo/grades_management/src/email"
	"github.com/alexizzarevalo/grades_management/src/msg"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
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

func ExportSheetsInPDF(opt SheetsOptions) {
	srv := getDriveService(opt.Credentials)
	srvSheets := getSheetService(opt.Credentials)

	newSpreadsheets, err := CopySheetsIntoSeparateSpreadSheets(srvSheets, opt.Id)
	if err != nil {
		msg.Error(errors.New("No se pudo copiar las hojas a un nuevo spreadsheet. " + err.Error()))
	}
	for _, newSpreadsheet := range newSpreadsheets {
		carnet := newSpreadsheet.Name
		pdfName := carnet + ".pdf"
		err := DeleteSheet(srvSheets, newSpreadsheet.SpreadsheetId, 0)
		if err != nil {
			msg.ErrorWithoutExit(errors.New("No se pudo eliminar la hoja del spreadsheet. " + err.Error()))
			continue
		}
		err = Export(srv, newSpreadsheet.SpreadsheetId, pdfName)
		if err != nil {
			msg.ErrorWithoutExit(errors.New("No se pudo exportar el archivo de Google Sheet " + pdfName + " " + err.Error()))
			continue
		}
		msg.Success("Pdf generado: " + pdfName)
	}
}

func ExportSheetsInPDFAndSendEmail(opt SheetsOptions, emailOpt email.EmailOptions) {
	ExportSheetsInPDF(opt)
	email.EmailOnly(emailOpt)
}
