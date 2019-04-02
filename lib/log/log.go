package log

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func Init(logLevel string) {
	Error = log.New(os.Stderr,
		"ERROR:\t",
		log.Ldate|log.Ltime)
	Warn = log.New(os.Stdout,
		"WARN:\t",
		log.Ldate|log.Ltime)
	Info = log.New(os.Stdout,
		"INFO:\t",
		log.Ldate|log.Ltime)
	Debug = log.New(os.Stdout,
		"DEBUG:\t",
		log.Ldate|log.Ltime)

	if logLevel == "debug" {
		return
	} else if logLevel == "info" {
		Debug.SetOutput(ioutil.Discard)
	} else if logLevel == "warn" {
		Debug.SetOutput(ioutil.Discard)
		Info.SetOutput(ioutil.Discard)
	} else if logLevel == "error" {
		Debug.SetOutput(ioutil.Discard)
		Info.SetOutput(ioutil.Discard)
		Warn.SetOutput(ioutil.Discard)
	} else {
		log.Fatal("Invalid log level " + logLevel)
	}
}
