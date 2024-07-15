package main

import (
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	if err := NewKvsServer().Serve(); err != nil {
		log.Fatal(err)
	}
}
