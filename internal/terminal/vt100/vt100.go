package vt100

import "fmt"

func Red(msg string) string {
	return fmt.Sprintf("\033[31m%s\033[39m", msg)
}

func Green(msg string) string {
	return fmt.Sprintf("\033[32m%s\033[39m", msg)
}
