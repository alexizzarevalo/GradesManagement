package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func getNameWithExt(fileName, ext string) string {
	return fileName + ext
}

func extractGrades(xlsx *excelize.File, carneCell, gradeCell string) {
	// Se extrae el carnet y la nota de las celdas espeficiadas
	fmt.Println("Carne,Nota")
	for _, sheetName := range xlsx.GetSheetMap() {
		carne := xlsx.GetCellValue(sheetName, carneCell)
		grade := xlsx.GetCellValue(sheetName, gradeCell)
		fmt.Printf("%v,%v\n", carne, grade)
	}
}

var gradeCell string
var carneCell string
var fileName string

func main() {

	flag.StringVar(&gradeCell, "nota", "", "La celda del excel donde esta la nota del alumno ej: A10")
	flag.StringVar(&carneCell, "carnet", "", "La celda del excel donde esta el carnet del alumno ej: D100")
	flag.StringVar(&fileName, "archivo", "", "La ruta del archivo xlsx")
	flag.Parse()

	if gradeCell == "" || carneCell == "" {
		flag.PrintDefaults()
		return
	}

	ext := filepath.Ext(fileName)
	fileName = strings.Replace(fileName, ext, "", 1)

	original := getNameWithExt(fileName, ext)

	xlsx, err := excelize.OpenFile(original)
	if err != nil {
		log.Fatal(errors.New("no se pudo abrir el archivo " + original))
	}

	extractGrades(xlsx, carneCell, gradeCell)
}
