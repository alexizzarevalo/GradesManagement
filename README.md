# Grades Management

Programa de consola para extraer las notas de alumnos de un archivo Excel y para envío de correos para optimizar tareas que realizo como auxiliar

## Install dependencies

	go mod tidy

## Build

	go build

## Install

	go install

## Run

```bash
# Para extraer las notas del Excel (Muestra los datos en consola)
./grades_management grades options.json

# Para extraer las notas del Excel (Escribe las notas en el archivo notas.csv)
./grades_management grades options.json > notas.csv


# Para enviar correos
./grades_management email options.json
```