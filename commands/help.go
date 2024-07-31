package commands

import "fmt"

func Help() {
	fmt.Println(`
Usage: rawd <command> [arguments]

Output a friendly greeting

commands:
	run		   			Run the program
	watch				Watch the program for changes and re-run
	init				Initialize a new project

Use "rawd help <command>" for more information about a command.`)
}

func HelpRun() {
	fmt.Println(
		`Usage: rawd run [path]

Run the program

optional arguments:
	path		Relative path to project directory. Default is the current directory.`)
}

func HelpWatch() {
	fmt.Println(
		`Usage: rawd watch [path]
		
Watch the program for changes and re-run

optional arguments:
	path		Relative path to project directory. Default is the current directory.`)
}

func HelpInit() {
	fmt.Println(
		`Usage: rawd init

Initialize a new project using CLI menu`)
}
