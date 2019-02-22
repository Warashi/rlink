package main

import (
	"flag"
	"log"
)

func main() {
	var relative, force bool
	flag.BoolVar(&relative, "r", false, "make symlink with relative path")
	flag.BoolVar(&force, "f", false, "create symlink if already exists file")
	flag.Parse()
	if err := MkLinks(flag.Arg(0), flag.Arg(1), relative, force); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
