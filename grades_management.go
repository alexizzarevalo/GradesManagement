package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func getNameWithExt(fileName, ext string) string {
	return fileName + "." + ext
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

func main() {
	fileName := "Notas"
	ext := "xlsx"

	original := getNameWithExt(fileName, ext)

	xlsx, err := excelize.OpenFile(original)
	if err != nil {
		log.Fatal(errors.New("no se pudo abrir el archivo" + original))
	}

	extractGrades(xlsx, "D3", "D103")
}
