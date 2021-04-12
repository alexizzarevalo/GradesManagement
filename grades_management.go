package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/alexizzarevalo/grades_management/src/email"
	"github.com/alexizzarevalo/grades_management/src/excel"
	"github.com/alexizzarevalo/grades_management/src/msg"
	"github.com/alexizzarevalo/grades_management/src/sheets"
)

type Options struct {
	Excel  excel.ExcelOptions
	Sheets sheets.SheetsOptions
	Email  email.EmailOptions
}

func main() {
	var optionsFile string
	if len(os.Args) == 3 {
		optionsFile = os.Args[2]
	} else {
		fmt.Println("Usage:")
		fmt.Println("grades_management [comando] options.json\n\r")
		fmt.Println("comandos:")
		fmt.Println("\tgrades         Extrae la notas de un Spreadsheet de Google Sheet")
		fmt.Println("\tgrades-excel   Extrae las notas de un Excel de Microsoft Office")
		fmt.Println("\texport         Exporta cada hoja de un Spreadsheet a PDF (Lo guarda en el directorio actual)")
		fmt.Println("\temail          Envia correo electronico a cada alumno con los PDF (Debe existir un csv de alumnos)")
		fmt.Println("\texport-email   Ejecuta el comando export y luego email")
		fmt.Println()
		fmt.Println("\toptions.json   Ruta del archivo de opciones en formato json (Puede tener cualquier nombre)")
		fmt.Println()
		fmt.Println("Para mas informacion de los comandos, visite: https://github.com/alexizzarevalo/GradesManagement#grades-management")
		os.Exit(1)
	}
	action := os.Args[1]

	content, err := os.ReadFile(optionsFile)
	if err != nil {
		msg.Error(errors.New("No se pudo abrir el archivo " + optionsFile))
	}

	opt := Options{}
	json.Unmarshal(content, &opt)

	if action == "grades" {
		sheets.GetGradesFromSpreadSheet(opt.Sheets)
	} else if action == "grades-excel" {
		excel.GetGrades(opt.Excel)
	} else if action == "export" {
		sheets.ExportSheetsInPDF(opt.Sheets)
	} else if action == "email" || action == "export-email" {
		msg.Info("Debes habilitar la opcion de aplicaciones inseguras en tu cuenta de Gmail: https://myaccount.google.com/lesssecureapps")
		msg.Info("Si tienes un segundo factor de autenticacion, en su lugar debes crear una password de aplicacion: https://myaccount.google.com/apppasswords\n\r")

		if action == "email" {
			email.EmailOnly(opt.Email)
		} else {
			sheets.ExportSheetsInPDFAndSendEmail(opt.Sheets, opt.Email)
		}
	}
}
