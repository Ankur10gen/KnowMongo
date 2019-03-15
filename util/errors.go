package util

import (
	"log"
	"os"
)

func BigError(msg string)  {
	log.Println(msg)
	os.Exit(1)
}

func SmallError(msg string)  {
	log.Println(msg)
}