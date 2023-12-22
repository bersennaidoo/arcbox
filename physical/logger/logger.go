package logger

import (
	"log"
	"os"

	glog "github.com/kataras/golog"
)

func New() *glog.Logger {
	logger := glog.New()
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger.SetOutput(os.Stdout)
	logger.SetLevel("debug")
	logger.SetLevelOutput("info", file)

	return logger
}
