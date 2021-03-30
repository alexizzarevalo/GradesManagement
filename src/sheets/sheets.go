package sheets

import (
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/sheets/v4"
)

type Cells struct {
	Grade string
	Carne string
}

type SheetsOptions struct {
	Id          string
	Credentials string
	Cells       Cells
}

func getSheetService(credentials string) *sheets.Service {
	srv, err := sheets.New(getHttpClient(credentials))

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return srv
}

func GetGradesFromSpreadSheet(opt SheetsOptions) {
	checkAppFolder() // Si no existe el folder de la app, se crea

	// Se extrae el carnet y la nota de las celdas espeficiadas
	srv := getSheetService(opt.Credentials)
	spreadsheetId := opt.Id

	spreadsheet := GetSpreadsheetById(srv, spreadsheetId)

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

func GetSheetByName(spreadsheet *sheets.Spreadsheet, name string) *sheets.Sheet {
	for _, sheet := range spreadsheet.Sheets {
		if strings.Compare(sheet.Properties.Title, name) == 0 {
			return sheet
		}
	}

	return nil
}

func createSpreadsheet(srv *sheets.Service, title string) string {
	spreadsheet, err := srv.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{Title: title},
	}).Do()

	if err != nil {
		log.Fatal(err)
	}

	return spreadsheet.SpreadsheetId
}

func GetSpreadsheetById(srv *sheets.Service, spreadsheetId string) *sheets.Spreadsheet {
	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		log.Fatalf("Error al intentar obtener la hoja de calculo. %v", err)
	}
	return spreadsheet
}

type NewSpreadsheet struct {
	SpreadsheetId string
	SheetId       int64
	Name          string
}

func CopySheetsIntoSeparateSpreadSheets(srv *sheets.Service, fromSpreadsheetId string) []NewSpreadsheet {
	// Se obtiene la spreadsheet deseada
	originalSpreadsheet := GetSpreadsheetById(srv, fromSpreadsheetId)
	var newSpreadsheetIds []NewSpreadsheet
	// Se recorren sus hojas
	for _, sheet := range originalSpreadsheet.Sheets {
		// Se crea un nuevo spreadsheet con el mismo nombre de la hoja actual
		newSpreadsheetId := createSpreadsheet(srv, sheet.Properties.Title)
		// Se copia la hoja del spreadsheet original al nuevo spreadsheet
		properties, err := srv.Spreadsheets.Sheets.CopyTo(
			originalSpreadsheet.SpreadsheetId,
			sheet.Properties.SheetId,
			&sheets.CopySheetToAnotherSpreadsheetRequest{
				DestinationSpreadsheetId: newSpreadsheetId,
			}).Do()

		if err != nil {
			log.Fatal(err)
		}

		newSpreadsheetIds = append(newSpreadsheetIds, NewSpreadsheet{
			SpreadsheetId: newSpreadsheetId,
			SheetId:       properties.SheetId,
			Name:          sheet.Properties.Title,
		})
	}

	return newSpreadsheetIds
}

func DeleteSheet(srv *sheets.Service, spreadsheetId string, sheetId int64) {
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteSheet: &sheets.DeleteSheetRequest{
					SheetId: sheetId,
				},
			},
		},
	}).Do()

	if err != nil {
		log.Fatal(err)
	}
}
