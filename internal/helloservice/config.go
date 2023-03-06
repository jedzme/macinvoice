package helloservice

import "fmt"

type Config struct {
	A string
	B bool
}

func (c Config) String() string {
	return fmt.Sprintf(`
=== Hello Service Config ===
A: %s
B: %t
`, c.A, c.B)
}
