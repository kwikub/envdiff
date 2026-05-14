package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/tagger"
)

func runTag(args []string) error {
	fs := flag.NewFlagSet("tag", flag.ContinueOnError)
	file := fs.String("file", "", "path to .env file (required)")
	tagSecrets := fs.Bool("secrets", false, "auto-tag secret-looking keys")
	tagEmpty := fs.Bool("empty", false, "auto-tag keys with empty values")
	extraRaw := fs.String("tag", "", "extra tag in name=KEY1,KEY2 format (repeatable via comma-separated)")
	filterTag := fs.String("filter", "", "only output entries carrying this tag")
	outputJSON := fs.Bool("json", false, "output results as JSON")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *file == "" {
		return fmt.Errorf("tag: --file is required")
	}

	entries, err := parser.Parse(*file)
	if err != nil {
		return fmt.Errorf("tag: parse %s: %w", *file, err)
	}

	opts := tagger.Options{
		TagSecrets: *tagSecrets,
		TagEmpty:   *tagEmpty,
		Extra:      make(map[string][]string),
	}

	if *extraRaw != "" {
		parts := strings.SplitN(*extraRaw, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("tag: --tag must be in name=KEY1,KEY2 format")
		}
		tagName := strings.TrimSpace(parts[0])
		keys := strings.Split(parts[1], ",")
		for i, k := range keys {
			keys[i] = strings.TrimSpace(k)
		}
		opts.Extra[tagName] = keys
	}

	results, err := tagger.Apply(entries, opts)
	if err != nil {
		return fmt.Errorf("tag: %w", err)
	}

	if *filterTag != "" {
		results = tagger.FilterByTag(results, *filterTag)
	}

	if *outputJSON {
		return json.NewEncoder(os.Stdout).Encode(results)
	}

	for _, r := range results {
		tagStr := ""
		if len(r.Tags) > 0 {
			tagStr = fmt.Sprintf(" [%s]", strings.Join(r.Tags, ", "))
		}
		fmt.Printf("%s=%s%s\n", r.Key, r.Value, tagStr)
	}
	return nil
}
