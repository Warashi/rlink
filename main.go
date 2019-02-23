package main

import (
	"flag"
	"log"
	"regexp"
	"strconv"
)

func main() {
	var relative, force, dryrun bool
	var ignore string
	flag.BoolVar(&relative, "r", false, "make symlink with relative path")
	flag.BoolVar(&force, "f", false, "create symlink if already exists file")
	flag.BoolVar(&dryrun, "d", false, "dry-run")
	flag.StringVar(&ignore, "i", "", "ignore pattern")
	flag.Parse()

	ign, err := regexp.Compile(ignore)
	if err != nil {
		log.Fatalf("regexp: Compile(%s): %+v", strconv.Quote(ignore), err)
	}

	if err := New(relative, force, dryrun, ign).MkLinks(flag.Arg(0), flag.Arg(1)); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
