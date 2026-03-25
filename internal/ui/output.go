package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
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

func Warn(msg string) {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("11")).
		Padding(0, 1)

	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Padding(0, 1)

	label := "WARN"
	fmt.Fprint(os.Stderr, labelStyle.Render(label))
	fmt.Fprintln(os.Stderr, msgStyle.Render(msg))
}

func Info(msg string) {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("12")).
		Padding(0, 1)

	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Padding(0, 1)

	label := "INFO"
	fmt.Fprint(os.Stderr, labelStyle.Render(label))
	fmt.Fprintln(os.Stderr, msgStyle.Render(msg))
}

func Select[T comparable](title string, options []huh.Option[T]) (T, error) {
	var result T

	err := huh.NewSelect[T]().
		Title(title).
		Options(options...).
		Value(&result).
		WithTheme(huh.ThemeBase16()).
		Run()

	return result, err
}

func Confirm(title string) (bool, error) {
	// TODO: Use Select yes/no instead?
	var result bool

	err := huh.NewConfirm().
		Title(title).
		Value(&result).
		WithTheme(huh.ThemeBase16()).
		Run()

	return result, err
}
