package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dwiandhikaap/rawdog-md/commands"

	"github.com/charmbracelet/huh"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		commands.Help()
		return
	}

	if args[0] == "help" {
		if len(args) == 1 {
			commands.Help()
			return
		}
		if args[1] == "build" {
			commands.HelpBuild()
			return
		}
		if args[1] == "watch" {
			commands.HelpWatch()
			return
		}
		if args[1] == "init" {
			commands.HelpInit()
			return
		}
		if args[1] == "version" {
			commands.HelpVersion()
			return
		}

		commands.Help()
		fmt.Println("\nUnknown command \"" + args[1] + "\"")
		return
	}

	if args[0] == "build" {
		firstArg := "."
		if len(args) > 1 {
			firstArg = args[1]
		}

		err := commands.Build(firstArg)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if args[0] == "watch" {
		// eat the first command
		args = args[1:]

		var port int = 3000
		for index, arg := range args {
			if arg == "--port" || arg == "-p" {
				// need at least 2 arguments (flag and value)
				if len(args) >= 2 {
					parsedPort, err := strconv.ParseInt(args[index+1], 10, 32)
					port = int(parsedPort)
					if err != nil {
						fmt.Println("Invalid port number")
						return
					}

					// eat the flag and value from the arguments
					args = append(args[:index], args[index+2:]...)
				} else {
					fmt.Println("Port number is unspecified")
					return
				}
			}
		}

		// get last argument
		path := "."
		if len(args) > 0 {
			path = args[len(args)-1]
		}

		err := commands.Watch(path, port)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if args[0] == "init" {
		var (
			projectName string = "my-blog"
			preset      string = "basic"
		)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Pick a project name").
					Value(&projectName).
					Validate(func(str string) error {
						if str == "" {
							return errors.New("project name cannot be empty")
						}

						// check if the directory already exists
						if _, err := os.Stat(str); !os.IsNotExist(err) {
							return fmt.Errorf("error: directory \"%s\" already exists", str)
						}
						return nil
					}),
			),
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Pick a preset").
					Options(
						huh.NewOption("skeleton - very minimal setup to get going", "skeleton"),
						huh.NewOption("basic - basic setup with a few pages and styling", "basic"),
						huh.NewOption("docs - documentation setup with side navigation and markdown", "docs"),
					).
					Value(&preset),
			),
		)

		err := form.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		err = commands.Init(projectName, preset)
		if err != nil {
			fmt.Println(err)
			return
		}

		return
	}

	if args[0] == "version" {
		commands.Version()
		return
	}

	commands.Help()
	fmt.Println("\nUnknown command \"" + args[0] + "\"")
}
