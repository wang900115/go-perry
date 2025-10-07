package func04

import "fmt"

func Greet(name string) string {
	if name == "" {
		return "Hello, Guest"
	}
	return fmt.Sprintf("Hello, %s", name)
}
