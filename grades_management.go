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
	action := os.Args[1]
	var optionsFile string
	if len(os.Args) == 3 {
		optionsFile = os.Args[2]
	} else {
		log.Fatal(errors.New("debe especificar la ruta del options.json"))
	}

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
	} else if action == "email" || action == "email-only" {
		fmt.Println("NOTE: You need to turn on 'less secure apps' options: https://myaccount.google.com/lesssecureapps")
		fmt.Println("NOTE: If you have SecondFactorAuthentication: Application-specific password required: https://myaccount.google.com/apppasswords\n\r")

		if action == "email" {
			sheets.ExportSheetsInPDFAndSendEmail(opt.Sheets, opt.Email)
		} else {
			email.EmailOnly(opt.Email)
		}
	}
}
