package logger

import (
	"log"
	"os"
)

var (
	INFO       *log.Logger
	WARN       *log.Logger
	ERRO       *log.Logger
	initalized bool = false
)

func Init() {
	if initalized {
		return
	}
	INFO = log.New(os.Stdout, "\x1b[32mINFO:\x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	WARN = log.New(os.Stdout, "\x1b[33mWARN:\x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	ERRO = log.New(os.Stderr, "\x1b[31mERRO:\x1b[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	initalized = true
}
