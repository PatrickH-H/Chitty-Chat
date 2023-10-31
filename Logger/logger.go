package Logger

import (
	"flag"
	"log"
	"os"
)

var (
	ErrorLogger *log.Logger
	FileLogger  *log.Logger
)

func init() {
	var logpath = "file.log"

	flag.Parse()
	var file, _ = os.Create(logpath)

	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Lshortfile)
	FileLogger = log.New(file, "", log.LstdFlags)
	FileLogger.SetFlags(0)

}
