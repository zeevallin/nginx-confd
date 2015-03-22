package main

import (
	"log"
)

func logFatalIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
