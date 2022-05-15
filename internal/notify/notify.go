package notify

import (
	"fmt"

	"github.com/fatih/color"
)

func PrettyError(msg string) {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", red("Ivy Error: "), msg)
}

func PrettyWarning(msg string) {
	red := color.New(color.FgYellow, color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", red("Ivy Warning: "), msg)
}

func PrettyInfo(msg string) {
	red := color.New(color.FgCyan, color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", red("Ivy Info: "), msg)
}
