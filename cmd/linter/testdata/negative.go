package main

import (
	"log"
	"os"
)

func negative() {
	os.Exit(0)     // want "os.Exit is used outside main package"
	log.Fatal()    // want "log.Fatal is used outside main package"
	panic("Panic") // want "panic call is used"
}
