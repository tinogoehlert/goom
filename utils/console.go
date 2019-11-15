package utils

import (
	"fmt"
	"os"

	"github.com/ttacon/chalk"
)

// Console is the GOOM Console, used to print GOOM messages
// to the system console.
type Console struct{}

// GoomConsole is the root console of GOOM.
var GoomConsole = &Console{}

// Greenf applies green console coloring to the given format and values.
func (c *Console) Greenf(format string, a ...interface{}) string {
	return fmt.Sprintf(chalk.Green.Color(format), a...)
}

// Green prints green console text using the given format and values.
func (c *Console) Green(format string, a ...interface{}) {
	fmt.Println(c.Greenf(format, a...))
}

// Redf applies red console coloring to the given format and values.
func (c *Console) Redf(format string, a ...interface{}) string {
	return fmt.Sprintf(chalk.Red.Color(format), a...)
}

// Red prints red console text using the given format and values.
func (c *Console) Red(format string, a ...interface{}) {
	fmt.Println(c.Redf(format, a...))
}

// Print prints console text in the default color using the given format and values.
func (c *Console) Print(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf(format, a...))
}

// Fatalf prints console text in the default color using the given format and values.
func (c *Console) Fatalf(format string, a ...interface{}) {
	fmt.Println(c.Redf(format, a...))
	os.Exit(2)
}
