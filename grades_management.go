package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Options struct {
	Grades GradesOptions
	Email  EmailOptions
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
		getGrades(opt.Grades)
	} else if action == "email" {
		sendEmail(opt.Email, []string{"dalexis.da@gmail.com"})
	}
}
