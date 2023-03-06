package csv

import "fmt"

type Config struct {
	A int32
	B bool
}

func (c Config) String() string {
	return fmt.Sprintf(`
=== CSV Service Config ===
A: %s
B: %t
`, c.A, c.B)
}
