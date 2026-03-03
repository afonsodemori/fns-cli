package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func HandleError(err error) {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("9")).
		Padding(0, 1)

	fmt.Fprintln(os.Stderr, style.Render(err.Error()))
	os.Exit(1)
}
