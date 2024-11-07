package dbg

import (
	"fmt"
	"log"
	"os"
)

type DebugTopic int

const (
	CONFIG DebugTopic = iota
	READY
	SESSION
	MANAGER
	ICE
	ELSE
)

var DebugTopicToStr = map[DebugTopic]string{
	CONFIG:  "CONFIG",
	READY:   "READY",
	SESSION: "SESSION",
	MANAGER: "MANAGER",
	ICE:     "ICE",
	ELSE:    "ELSE",
}

type DebugMode int

const (
	SILENT DebugMode = iota
	STDOUT
	SINGLEFILE
)

var DebugModeToStr = map[DebugMode]string{
	SILENT:     "SILENT",
	STDOUT:     "STDOUT",
	SINGLEFILE: "SINGLEFILE",
}

var Mode DebugMode = SILENT
var FileLogger *log.Logger
var File *os.File

func Init(mode DebugMode) error {
	switch mode {
	case SILENT:
		Mode = SILENT
	case STDOUT:
		Mode = STDOUT
	case SINGLEFILE:
		Mode = SINGLEFILE
		var err error
		err = os.MkdirAll("./log/", 0755)
		if err != nil {
			return err
		}
		File, err = os.CreateTemp("./log/", "log")
		if err != nil {
			return err
		}
		FileLogger = log.New(File, "", log.LstdFlags)
	}
	return nil
}

func Close() {
	if File != nil {
		File.Close()
	}
}

func Println(topic DebugTopic, a ...interface{}) {
	switch Mode {
	case SILENT:
	case STDOUT:
		fmt.Printf("[%v] %v\n", DebugTopicToStr[topic], fmt.Sprint(a...))
	case SINGLEFILE:
		FileLogger.Printf("[%v] %v\n", DebugTopicToStr[topic], fmt.Sprint(a...))
	}
}

func Fatal(topic DebugTopic, a ...interface{}) {
	Println(topic, a...)
	os.Exit(1)
}
