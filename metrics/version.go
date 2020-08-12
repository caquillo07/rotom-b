package metrics

import "time"

var (

	// Version is the version of the bot, injected at build time
	Version string

	// Commit is the commit used to build this version, injected at build time
	Commit string

	// Branch is the branch used to build the this version, injected at build
	// time
	Branch string

	// Date is when this binary was built, injected at build time
	Date time.Time
)
