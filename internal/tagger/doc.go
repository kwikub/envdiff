// Package tagger annotates .env entries with user-defined and automatic tags.
//
// Tags can be applied automatically (e.g. "secret" for sensitive keys,
// "empty" for blank values) or explicitly via the Options.Extra map.
// Once tagged, entries can be filtered with FilterByTag to isolate
// subsets of interest for reporting, masking, or export pipelines.
package tagger
