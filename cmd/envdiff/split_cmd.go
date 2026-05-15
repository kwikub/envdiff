package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/splitter"
)

// runSplit is the entry-point for the `envdiff split` sub-command.
//
// Usage:
//
//	envdiff split --src .env --dest ./split [--strategy prefix|alpha] [--buckets 4]
//	envdiff split --src .env --dest ./split --map 'db.env=DB_,REDIS_;app.env=APP_'
func runSplit(args []string) error {
	var (
		src      = ".env"
		dest     = "./split"
		strategy = "prefix"
		buckets  = 2
		mapStr   = ""
	)

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--src":
			i++
			src = args[i]
		case "--dest":
			i++
			dest = args[i]
		case "--strategy":
			i++
			strategy = args[i]
		case "--buckets":
			i++
			fmt.Sscanf(args[i], "%d", &buckets)
		case "--map":
			i++
			mapStr = args[i]
		}
	}

	opts := splitter.Options{
		Buckets: buckets,
	}

	switch strategy {
	case "alpha":
		opts.Strategy = splitter.StrategyAlpha
	default:
		opts.Strategy = splitter.StrategyPrefix
	}

	if mapStr != "" {
		opts.Mapping = parseMapFlag(mapStr)
	}

	counts, err := splitter.Split(src, dest, opts)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	fmt.Fprintf(os.Stdout, "Split %s → %s\n", src, dest)
	return enc.Encode(counts)
}

// parseMapFlag parses "db.env=DB_,REDIS_;app.env=APP_" into a mapping.
func parseMapFlag(s string) map[string][]string {
	out := map[string][]string{}
	for _, part := range strings.Split(s, ";") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		file := strings.TrimSpace(kv[0])
		prefixes := strings.Split(kv[1], ",")
		for i, p := range prefixes {
			prefixes[i] = strings.TrimSpace(p)
		}
		out[file] = prefixes
	}
	return out
}
