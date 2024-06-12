package internal

import (
	"github.com/pkg/errors"
	"log"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func CheckError(err error) {
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			log.Printf("%+s:%d\n", f, f)
			log.Println(err)
		}
	}

	if err != nil {
		log.Println(err.Error())
	}
}
