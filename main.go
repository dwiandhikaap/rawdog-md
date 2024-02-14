package main

import (
	"fmt"
	"os"

	"rawdog-md/commands"
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
		fmt.Println("Unknown command \"" + args[1] + "\"")
		return
	}

	if args[0] == "run" {
		if len(args) == 1 {
			err := commands.Run(".")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err := commands.Run(args[1])
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
}
