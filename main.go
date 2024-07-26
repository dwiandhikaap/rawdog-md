package main

import (
	"fmt"
	"os"
	"time"

	"rawdog-md/commands"
)

func repeatingNonsense(ch chan int, value int, wait time.Duration) {
	for {
		time.Sleep(wait * time.Second)
		ch <- value
	}
}

func main2() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go repeatingNonsense(ch1, 420, 1)
	go repeatingNonsense(ch1, 69, 3)

	go func() {
		for {
			fmt.Println("Waiting for value...")
			select {
			case v1 := <-ch1:
				fmt.Println(v1)
			case v2 := <-ch2:
				fmt.Println(v2)
			}
		}
	}()

	<-make(chan int)
}

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
}
