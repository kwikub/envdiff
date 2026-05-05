package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	maskSecrets := flag.Bool("mask", true", "Mask secret values in output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <base.env> <target.env>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	baseEntries, err := parser.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing base file: %v\n", err)
		os.Exit(1)
	}

	targetEntries, err := parser.Parse(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing target file: %v\n", err)
		os.Exit(1)
	}

	diffs := differ.Diff(baseEntries, targetEntries)

	if *maskSecrets {
		diffs = differ.MaskSecrets(diffs)
	}

	fmt := reporter.FormatText
	if *format == "json" {
		fmt = reporter.FormatJSON
	}

	if err := reporter.Report(os.Stdout, diffs, fmt); err != nil {
		fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
		os.Exit(1)
	}
}
