package main

import (
	"errors"
	"fmt"
	"os"

	"rawdog-md/commands"

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
		if args[1] == "run" {
			commands.HelpRun()
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
		fmt.Println("Unknown command \"" + args[1] + "\"")
		return
	}

	if args[0] == "run" {
		if len(args) > 1 {
			os.Chdir(args[1])
		}

		err := commands.Run()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if args[0] == "watch" {
		if len(args) == 1 {
			err := commands.Watch(".")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err := commands.Watch(args[1])
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if args[0] == "init" {
		var (
			projectName string
			preset      string
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

		fmt.Printf("Project initialized at './%s' with preset '%s'\n", projectName, preset)
		fmt.Println("\nRun the following commands to get started:")
		fmt.Println("  cd", projectName)
		fmt.Println("\n(Optional, if you want to use git):")
		fmt.Println("  git init")
		fmt.Println("\nBegin development:")
		fmt.Println("  rawd watch")

		return
	}
}
