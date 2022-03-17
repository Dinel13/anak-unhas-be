package helper

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

var errorLog *log.Logger

func logErorr(err error) {
	errorLog = log.New(os.Stdin, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	errorLog.Println(trace)
}
