package main

import (
    "io"
    "log"
)

var (
    Error *log.Logger // Logger for logging errors
    Info *log.Logger // Logger for logging important information
    Warn *log.Logger // Logger for logging anomolies that could lead to undesired behavior down the line
)

// InitLog: Initialize loggers
// Parameters:
//  errorWriter (io.Writer) - Writer to log errros to
//  infoWriter (io.Writer) - Writer to log info to
//  warnWriter (io.Writer) - Writer to log warnings to
func InitLog(errorWriter, infoWriter, warnWriter io.Writer) {
    Error = log.New(errorWriter, "[X] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
    Info = log.New(infoWriter, "(i) INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    Warn = log.New(warnWriter, "/!\\ WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}
