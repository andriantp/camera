package common

import (
	"time"
)

const (
	Timeout = 30 * time.Second
	// Wait    = 5 * time.Second
	Wait = 500 * time.Millisecond

	// Defaults used when Profile / Channel are not set in .env.
	DefaultProfile = "profile_1"
	DefaultChannel = 1
)
