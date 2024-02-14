package commands

import "fmt"

func Help() {
	fmt.Println(`
Usage: rawd <command> [arguments]

Output a friendly greeting

commands:
	run		   			Run the program

Use "rawd help <command>" for more information about a command.`)
}

func HelpRun() {
	fmt.Println(
`Usage: rawd run [path]

Run the program

optional arguments:
	path		Relative path to project directory. Default is the current directory.`)
}
