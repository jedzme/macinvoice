package http

import (
	"fmt"
	"time"
)

type Config struct {
	MaxRetries    int
	ClientTimeout time.Duration
}

func (c Config) String() string {
	return fmt.Sprintf(`
=== HTTP Service Config ===
MaxRetries: %d
ClientTimeout: %s
`, c.MaxRetries, c.ClientTimeout)
}
