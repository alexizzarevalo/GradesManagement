package excel

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/alexizzarevalo/grades_management/src/msg"
)

type Cells struct {
	Grade string
	Carne string
}

type ExcelOptions struct {
	File  string
	Cells Cells
}

func getNameWithExt(fileName, ext string) string {
	return fileName + ext
}

func extractGrades(xlsx *excelize.File, carneCell, gradeCell string) {
	// Se extrae el carnet y la nota de las celdas espeficiadas
	fmt.Println("Carne,Nota")
	var omited = "Se omitio: "
	for _, sheetName := range xlsx.GetSheetMap() {
		carne := xlsx.GetCellValue(sheetName, carneCell)
		grade := xlsx.GetCellValue(sheetName, gradeCell)
		if strings.Compare(carne, "") != 0 && strings.Compare(grade, "") != 0 {
			fmt.Printf("%v,%v\n", carne, grade)
		} else {
			omited += sheetName + ", "
		}
	}
	if strings.Compare(omited, "Se omitio: ") != 0 {
		msg.Warning(omited + "porque no se encontro carnet o nota.")
	}
}

func GetGrades(opt ExcelOptions) {
	ext := filepath.Ext(opt.File)
	fileName := strings.Replace(opt.File, ext, "", 1)

	original := getNameWithExt(fileName, ext)

	xlsx, err := excelize.OpenFile(original)
	if err != nil {
		msg.Error(errors.New("no se pudo abrir el archivo " + original))
	}

	extractGrades(xlsx, opt.Cells.Carne, opt.Cells.Grade)
}
