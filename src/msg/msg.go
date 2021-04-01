package msg

import (
	"fmt"
	"log"
	"os"

	"github.com/TwinProduction/go-color"
)

func Error(err error) {
	fmt.Println(color.Ize(color.Red, err.Error()))
	os.Exit(1)
}

func ErrorWithoutExit(err error) {
	fmt.Println(color.Ize(color.Red, err.Error()))
}

func Warning(msg string) {
	fmt.Println(color.Ize(color.Yellow, "WARNING: "+msg))
}

func Success(msg string) {
	log.SetPrefix("0")
	log.Println(color.Ize(color.Green, msg))
}

func Info(msg string) {
	fmt.Println(color.Ize(color.Blue, "INFO: "+msg))
}
