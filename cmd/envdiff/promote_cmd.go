package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/promoter"
)

func runPromote(args []string) error {
	fs := flag.NewFlagSet("promote", flag.ContinueOnError)
	skipExisting := fs.Bool("skip-existing", false, "skip keys already present in the target")
	dryRun := fs.Bool("dry-run", false, "report changes without writing the target file")
	keys := fs.String("keys", "", "comma-separated allow-list of keys to promote")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("usage: envdiff promote [flags] <src> <dst>")
	}

	src := fs.Arg(0)
	dst := fs.Arg(1)

	var keyList []string
	if *keys != "" {
		for _, k := range strings.Split(*keys, ",") {
			k = strings.TrimSpace(k)
			if k != "" {
				keyList = append(keyList, k)
			}
		}
	}

	opts := promoter.Options{
		SkipExisting: *skipExisting,
		DryRun:       *dryRun,
		Keys:         keyList,
	}

	results, err := promoter.Promote(src, dst, opts)
	if err != nil {
		return err
	}

	for _, r := range results {
		switch r.Action {
		case "promoted":
			if r.OldValue == "" {
				fmt.Printf("+ %s=%s\n", r.Key, r.NewValue)
			} else {
				fmt.Printf("~ %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
			}
		case "skipped":
			fmt.Printf("= %s (skipped)\n", r.Key)
		case "dry-run":
			if r.OldValue == "" {
				fmt.Printf("[dry-run] + %s=%s\n", r.Key, r.NewValue)
			} else {
				fmt.Printf("[dry-run] ~ %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
			}
		}
	}
	return nil
}
