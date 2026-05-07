package inspector

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Print writes a human-readable summary to w.
func Print(w io.Writer, s *Summary) {
	fmt.Fprintf(w, "Total keys   : %d\n", s.TotalKeys)
	fmt.Fprintf(w, "Secret keys  : %d\n", s.SecretKeys)
	fmt.Fprintf(w, "Empty values : %d\n", s.EmptyValues)
	fmt.Fprintf(w, "Longest key  : %s\n", s.LongestKey)

	if len(s.DuplicateKeys) > 0 {
		sort.Strings(s.DuplicateKeys)
		fmt.Fprintf(w, "Duplicates   : %s\n", strings.Join(s.DuplicateKeys, ", "))
	} else {
		fmt.Fprintf(w, "Duplicates   : none\n")
	}

	if len(s.Groups) > 0 {
		sort.Strings(s.Groups)
		fmt.Fprintf(w, "Groups       : %s\n", strings.Join(s.Groups, ", "))
	} else {
		fmt.Fprintf(w, "Groups       : none\n")
	}
}
