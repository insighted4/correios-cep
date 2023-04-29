// Package version contains version information for this app.
package version

import (
	"time"
)

// Version is set by the build scripts.
var (
	BuildTime  = time.Now().In(time.UTC).Format(time.Stamp + " 2006 UTC")
	CommitHash = ""
	Version    = ""
)
