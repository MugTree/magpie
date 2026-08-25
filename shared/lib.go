package shared

import (
	"fmt"
	"log"
	"runtime"
)

func LogError(err error) {
	_, file, line, ok := runtime.Caller(1) // 1 = caller of this function
	msg := err.Error()
	if ok {
		msg = fmt.Sprintf("%s (at %s:%d)", msg, file, line)
	}
	log.Println("ERROR: " + msg)
}
