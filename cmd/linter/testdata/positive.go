package main

import (
	"log"
	"os"
)

func main() {
	os.Exit(0)
	log.Fatal()
}

func positive() {
	log.Print("Positive func")
}

func mainNegative() {
	panic("Panic") // want "panic call is used"
}
