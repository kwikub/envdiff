// Package grouper partitions a slice of env entries into named groups based on
// key prefix segments.
//
// Keys are split on a configurable separator (default "_"). Everything before
// the first separator becomes the group name. For example:
//
//	DB_HOST  → group "DB"
//	DB_PORT  → group "DB"
//	APP_NAME → group "APP"
//
// Keys that contain no separator are placed in a synthetic "(ungrouped)" group
// when Options.IncludeUngrouped is set to true.
//
// Groups are returned in alphabetical order; entries within each group preserve
// their original order.
package grouper
