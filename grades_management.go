package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alexizzarevalo/grades_management/src/email"
	"github.com/alexizzarevalo/grades_management/src/excel"
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
		fmt.Println("grades_management necesita la accion a realizar seguido de la ruta del archivo de opciones")
		fmt.Println("grades_management [grades|grades-excel|email|email-only] options.json")
		os.Exit(1)
	}
	action := os.Args[1]

	content, err := os.ReadFile(optionsFile)
	if err != nil {
		log.Fatal(errors.New("No se pudo abrir el archivo " + optionsFile))
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
		fmt.Println("NOTE: You need to turn on 'less secure apps' options: https://myaccount.google.com/lesssecureapps")
		fmt.Println("NOTE: If you have SecondFactorAuthentication: Application-specific password required: https://myaccount.google.com/apppasswords\n\r")

		if action == "email" {
			email.EmailOnly(opt.Email)
		} else {
			sheets.ExportSheetsInPDFAndSendEmail(opt.Sheets, opt.Email)
		}
	}
}
