package main

import (
	"github.com/Komilov31/wget/internal/flags"
	"github.com/Komilov31/wget/internal/parser"
	"github.com/Komilov31/wget/internal/wget"
	"github.com/Komilov31/wget/internal/worker"
)

func main() {
	flags := flags.Parse()
	parser := parser.New(flags.FlagD)
	worker := worker.New(flags.Url, flags.FlagN, parser)
	wget := wget.New(worker)

	wget.Run()
}
