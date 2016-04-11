package golagraphite

import (
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {

	logf := &lumberjack.Logger{
		Filename:   "golagraphite.log",
		MaxSize:    128, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	}

	log.SetOutput(io.MultiWriter(logf, os.Stdout))
}
