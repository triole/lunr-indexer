package main

import "fmt"

func potentialEmptyLine() {
	if CLI.Watch == false {
		fmt.Printf("\n")
	}
}
