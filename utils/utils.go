package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type GithubResponse struct {
	Name       string `json:"name"`
	ProfileUrl string `json:"html_url"`
	Company    string `json:"company"`
	Location   string `json:"location"`
}

func GetUserInput() string{
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to github go program")
	fmt.Println("---------------------")

	for {
		fmt.Println("Please provide a username of a user you would like to lookup")
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		fmt.Println(text)
		return text
	}
}

