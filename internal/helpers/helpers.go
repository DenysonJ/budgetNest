package helpers

import (
	"fmt"
	"log"
)

// CheckError prints the provided message and error if the error is not nil.
func CheckError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
	}
}

// CheckFatal prints the provided message and error if the error is not nil, then exits the program.
func CheckFatal(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v\n", msg, err)
	}
}
