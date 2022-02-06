package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type GithubResponse struct {
	Name       string `json:"name"`
	ProfileUrl string `json:"html_url"`
	Company    string `json:"company"`
	Location   string `json:"location"`
}

type GithubRepoRequirements struct {
	Name           string `json:"name"`
	AutoInit       bool   `json:"autoInit"`
	Private        bool   `json:"private"`
	GitignoreTempl string `json:"gitignore_template"`
}

func AskUserForGithubUsername() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Please enter the username: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		return text
	}
}

func GetUserFromGithub() {

	username := AskUserForGithubUsername()
	fmt.Println("Looking up user...")
	user_to_search := fmt.Sprintf("https://api.github.com/users/%s", username)
	response, err := http.Get(user_to_search)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Processing data...")
	fmt.Println("")
	var responseObject GithubResponse

	json.Unmarshal(responseData, &responseObject)

	fmt.Printf("Name: %s\nProfile url: %s\nCompany: %s\nLocation: %s\n", responseObject.Name, responseObject.ProfileUrl, responseObject.Company, responseObject.Location)
	fmt.Println("")
}
