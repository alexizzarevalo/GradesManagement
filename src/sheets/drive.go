package sheets

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alexizzarevalo/grades_management/src/email"
	"google.golang.org/api/drive/v3"
)

func getDriveService(credentials string) *drive.Service {
	srv, err := drive.New(getHttpClient(credentials))

	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return srv
}

func ListFiles(srv *drive.Service) {
	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
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

func Export(srv *drive.Service, spreadsheetId, name string) {
	resp, err := srv.Files.Export(spreadsheetId, "application/pdf").Download()
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func ExportSheetsInPDFAndSendEmail(opt SheetsOptions, emailOpt email.EmailOptions) {
	srv := getDriveService(opt.Credentials)
	srvSheets := getSheetService(opt.Credentials)
	students := email.ReadStudentsCsv()

	newSpreadsheets := CopySheetsIntoSeparateSpreadSheets(srvSheets, opt.Id)
	for _, newSpreadsheet := range newSpreadsheets {
		carnet := newSpreadsheet.Name
		pdfName := carnet + ".pdf"
		DeleteSheet(srvSheets, newSpreadsheet.SpreadsheetId, 0)
		Export(srv, newSpreadsheet.SpreadsheetId, pdfName)
		fmt.Println("Pdf generado: ", pdfName)

		to, err := email.GetEmailByCarne(carnet, students, emailOpt.StudentsCsv)
		if err != nil {
			fmt.Println(err)
			continue
		}
		email.SendEmailWithAttachment(emailOpt, []string{to}, pdfName, pdfName)
	}
}
