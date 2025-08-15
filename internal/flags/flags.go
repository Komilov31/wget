package flags

import (
	"log"

	"github.com/pborman/getopt"
)

type Flags struct {
	FlagN int
	FlagD int
	Url   string
}

func Parse() *Flags {
	flagN := getopt.Int('n', 3, "number of workers to download resourses in parallel")
	flagD := getopt.Int('d', 1, "max recursion depth(levels of downloads)")
	getopt.Parse()

	flags := Flags{
		FlagD: *flagD,
		FlagN: *flagN,
	}

	args := getopt.Args()
	if len(args) == 0 {
		log.SetFlags(0)
		log.Fatal("wget: missing url")
	}

	if len(args) > 0 {
		flags.Url = args[0]
	}

	return &flags
}
