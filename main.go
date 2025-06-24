package main

import (
	"log"
	"kubefix-cli/cmd"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd.Execute()
}


