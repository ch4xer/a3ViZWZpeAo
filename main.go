package main

import (
	"kubefix-cli/cmd"
	"log"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd.Execute()
}
