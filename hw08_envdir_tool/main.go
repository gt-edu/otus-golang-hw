package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Invalid amount of arguments")
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Error during reading variables: %v", err)
	}

	os.Exit(RunCmd(os.Args[2:], env))
}
