package main

import (
	"fmt"
	"log"

	"github.com/aiq/go-argsfile"
)

func main() {

	args, err := argsfile.Args()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("resolved args:")
	fmt.Println("--------------")
	for i, v := range args {
		fmt.Printf("%  d. %s\n", i, v)
	}
}
