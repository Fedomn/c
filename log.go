package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

const enableDebug = false // enable for debugging

var logger *log.Logger

func init() {
	if enableDebug {
		logPath := filepath.Dir(os.Args[0]) + "/debug.log"
		fd, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			color.Red("Failed to init debug.log.")
			os.Exit(1)
		}

		logger = log.New(fd, "", log.LstdFlags)
		logger.Println("Init debug logger successfully.")
	}
}

func debug(format string, v ...interface{}) {
	if enableDebug {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		if len(v) == 0 {
			logger.Printf(format)
		} else {
			logger.Printf(format, v...)
		}
	}
}
