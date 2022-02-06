package main

import (
	"encoding/json"
	"fmt"
	"go-github/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	userInput := utils.GetUserInput()
	fmt.Printf("Searching github for %s\n", userInput)
	user_to_search := fmt.Sprintf("https://api.github.com/users/%s", userInput)

	response, err := http.Get(user_to_search)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject utils.GithubResponse

	json.Unmarshal(responseData, &responseObject)

	fmt.Printf("Name: %s\nProfile url: %s\nCompany: %s\nLocation: %s", responseObject.Name, responseObject.ProfileUrl,responseObject.Company,responseObject.Location)
}
	
