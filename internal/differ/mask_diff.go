package differ

// MaskSecrets returns a copy of diffs with secret values replaced by "***".
// It uses IsSecret to determine which keys should be masked.
func MaskSecrets(diffs []DiffEntry) []DiffEntry {
	masked := make([]DiffEntry, len(diffs))
	for i, d := range diffs {
		entry := d
		if IsSecret(d.Key) {
			if entry.OldValue != "" {
				entry.OldValue = "***"
			}
			if entry.NewValue != "" {
				entry.NewValue = "***"
			}
		}
		masked[i] = entry
	}
	return masked
}
