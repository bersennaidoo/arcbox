package logger

import (
	"log"
	"os"

	glog "github.com/kataras/golog"
)

func New() *glog.Logger {
	logger := glog.New()
	infolog, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	errorlog, err := os.OpenFile("errorlog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger.SetOutput(os.Stdout)
	logger.SetLevel("debug")
	logger.SetLevelOutput("info", infolog)
	logger.SetLevelOutput("error", errorlog)

	return logger
}
