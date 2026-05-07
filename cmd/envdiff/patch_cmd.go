package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/patcher"
)

// runPatch handles the "patch" sub-command.
// Usage: envdiff patch -f .env -set KEY=VALUE [-set ...] [-del KEY] [-del ...] [-write]
func runPatch(args []string) error {
	fs := flag.NewFlagSet("patch", flag.ContinueOnError)

	var (
		filePath  = fs.String("f", ".env", "path to the .env file to patch")
		writeBack = fs.Bool("write", false, "write patched result back to the file")
		setFlags  multiFlag
		delFlags  multiFlag
	)
	fs.Var(&setFlags, "set", "set KEY=VALUE (repeatable)")
	fs.Var(&delFlags, "del", "delete KEY (repeatable)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	var patches []patcher.Patch

	for _, kv := range setFlags {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("patch: -set value must be KEY=VALUE, got %q", kv)
		}
		patches = append(patches, patcher.Patch{Op: patcher.OpSet, Key: parts[0], Value: parts[1]})
	}

	for _, key := range delFlags {
		patches = append(patches, patcher.Patch{Op: patcher.OpDelete, Key: key})
	}

	res, err := patcher.Apply(*filePath, patches)
	if err != nil {
		return err
	}

	output := patcher.Format(res.Entries)

	if *writeBack {
		if err := os.WriteFile(*filePath, []byte(output), 0o600); err != nil {
			return fmt.Errorf("patch: write %q: %w", *filePath, err)
		}
		fmt.Fprintf(os.Stderr, "patch: applied %d change(s), skipped %d\n",
			len(res.Applied), len(res.Skipped))
	} else {
		fmt.Print(output)
	}

	return nil
}

// multiFlag is a flag.Value that collects repeated -flag values.
type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ", ") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }
