package sheets

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alexizzarevalo/grades_management/src/msg"
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
		msg.Error(errors.New("No se pudo recuperar el cliente de Sheets. " + err.Error()))
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
		msg.Error(errors.New("Error al intentar obtener los valores de carnet y nota. " + err.Error()))
	}

	fmt.Println("Carne,Nota")
	var omited = "Se omitio: "
	for i := 0; i < len(values.ValueRanges); i += 2 {
		if len(values.ValueRanges[i].Values) == 1 && len(values.ValueRanges[i+1].Values) == 1 {
			carne := values.ValueRanges[i].Values[0][0]
			grade := values.ValueRanges[i+1].Values[0][0]
			fmt.Printf("%v,%v\n", carne, grade)
		} else {
			omited += strings.Split(values.ValueRanges[i].Range, "!")[0] + ", "
		}
	}
	if strings.Compare(omited, "Se omitio: ") != 0 {
		msg.Warning(omited + "porque no se encontro carnet o nota.")
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
		msg.Error(errors.New("No se pudo crear la spreadsheet " + title + " " + err.Error()))
	}

	return spreadsheet.SpreadsheetId
}

func GetSpreadsheetById(srv *sheets.Service, spreadsheetId string) *sheets.Spreadsheet {
	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		msg.Error(errors.New("No se pudo obtener la spreadsheet. " + err.Error()))
	}
	return spreadsheet
}

type NewSpreadsheet struct {
	SpreadsheetId string
	SheetId       int64
	Name          string
}

func CopySheetsIntoSeparateSpreadSheets(srv *sheets.Service, fromSpreadsheetId string) ([]NewSpreadsheet, error) {
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
			return nil, err
		}

		newSpreadsheetIds = append(newSpreadsheetIds, NewSpreadsheet{
			SpreadsheetId: newSpreadsheetId,
			SheetId:       properties.SheetId,
			Name:          sheet.Properties.Title,
		})
	}

	return newSpreadsheetIds, nil
}

func DeleteSheet(srv *sheets.Service, spreadsheetId string, sheetId int64) error {
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteSheet: &sheets.DeleteSheetRequest{
					SheetId: sheetId,
				},
			},
		},
	}).Do()

	return err
}
