// Package interpolator resolves variable references embedded in .env file
// values using ${VAR} or $VAR syntax.
//
// References are looked up first in the provided entry slice, then in an
// optional overrides map, and optionally in the OS environment when
// Options.FallbackToOS is enabled.
//
// In strict mode an unresolved reference causes an error; otherwise the
// original token is left in place.
package interpolator
