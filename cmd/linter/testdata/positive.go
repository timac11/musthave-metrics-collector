package main

import (
	"log"
	"os"
)

func main() {
	os.Exit(0)
	log.Fatal()
}

func mainNegative() {
	panic("Panic") // want "panic call is used"
}
