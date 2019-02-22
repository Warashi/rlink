package main

import (
	"flag"
	"log"
)

func main() {
	var relative, force bool
	var ignore string
	flag.BoolVar(&relative, "r", false, "make symlink with relative path")
	flag.BoolVar(&force, "f", false, "create symlink if already exists file")
	flag.StringVar(&ignore, "i", "", "ignore pattern")
	flag.Parse()
	if err := MkLinks(flag.Arg(0), flag.Arg(1), relative, force, ignore); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
