package commands

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func Help() {
	style1 := lipgloss.NewStyle()
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8dd2"))

	style3 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4a55f2"))

	style4 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4a55f2")).
		Width(24)

	fmt.Println(style1.Render("Usage: ") +
		style2.Render("rawd ") +
		style3.Render("[command] [arguments]") + "\n\n" +
		style1.Render("Output a friendly greeting") + "\n\n" +
		style1.Render("commands:") + "\n" +
		style4.Render("\tbuild") +
		style1.Render("Build the project") + "\n" +
		style4.Render("\twatch") +
		style1.Render("Watch the program for changes and re-build") + "\n" +
		style4.Render("\tinit") +
		style1.Render("Initialize a new project") + "\n" +
		style4.Render("\tversion") +
		style1.Render("Output the version of the program") + "\n\n" +
		style1.Render("Use \"") +
		style2.Render("rawd help <command>") +
		style1.Render("\" for more information about a command."))
}

func HelpBuild() {
	style1 := lipgloss.NewStyle()
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8dd2"))
	style3 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4a55f2")).
		Width(24)

	fmt.Println(style1.Render("Usage:") +
		style2.Render(" rawd build") +
		style3.Render(" [path]") + "\n\n" +
		style1.Render("Build the project") + "\n\n" +
		style1.Render("optional arguments:") + "\n" +
		style3.Render("\tpath") +
		style1.Render("Relative path to project directory. Default is the current directory."))
}

func HelpWatch() {
	style1 := lipgloss.NewStyle()
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8dd2"))
	style3 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4a55f2")).
		Width(24)

	fmt.Println(style1.Render("Usage:") +
		style2.Render(" rawd watch") +
		style3.Render(" [options] [path]") + "\n\n" +
		style1.Render("Watch the program for changes and re-build") + "\n\n" +
		style1.Render("optional arguments:") + "\n" +

		style3.Render("\t-p, --port [port]") +
		style1.Render("Port number to run the server on. Default is 3000.") + "\n" +
		style3.Render("\tpath") +
		style1.Render("Relative path to project directory. Default is the current directory."))
}

func HelpInit() {
	style1 := lipgloss.NewStyle()
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8dd2"))

	fmt.Println(style1.Render("Usage:") +
		style2.Render(" rawd init") + "\n\n" +
		style1.Render("Initialize a new project using CLI menu"))
}

func HelpVersion() {
	style1 := lipgloss.NewStyle()
	style2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff8dd2"))

	fmt.Println(style1.Render("Usage:") +
		style2.Render(" rawd version") + "\n\n" +
		style1.Render("Output the version of the program"))
}
