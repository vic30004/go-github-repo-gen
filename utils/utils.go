package utils

import (
	"bufio"
	"fmt"
	"go-github/api"
	"os"
	"strings"
)

func printOptions() {
	fmt.Println("What would you like to do?")
	fmt.Println("(A) Look up user on github")
	fmt.Println("(B) Create a repo with template")
	fmt.Println("(C) Create a new repo")
	fmt.Println("(E) To exit")
	fmt.Print("-> ")
}

func GetUserInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to github go program")
	fmt.Println("---------------------")

	for {
		printOptions()
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if strings.ToUpper(text) != "A" && strings.ToUpper(text) != "B" && strings.ToUpper(text) != "E" && strings.ToUpper(text) != "C" {
			fmt.Println("Please select a valid option")
		} else {
			if strings.ToUpper(text) == "A" {
				api.GetUserFromGithub()
			} else if strings.ToUpper(text) == "E" {
				fmt.Println("Hope to see you again!")
				return
			} else if strings.ToUpper(text) == "B" {
				api.CreateRepoFromTemplate()
			} else if strings.ToUpper(text) == "C" {
				api.CreateNewRepo()
			}
		}
	}
}
