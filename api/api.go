package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type GithubResponse struct {
	Name       string `json:"name"`
	ProfileUrl string `json:"html_url"`
	Company    string `json:"company"`
	Location   string `json:"location"`
}

type GithubRepoRequirements struct {
	Name        string `json:"name"`
	Homepage    string `json:"homepage"`
	Description string `json:"description"`
	Url         string `json:"html_url"`
}

func askUserForInput(input string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		req := fmt.Sprintf("Please enter the %s: ", input)
		fmt.Print(req)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		return text
	}
}

func goDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}

func GetUserFromGithub() {

	username := askUserForInput("username")
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

func CreateRepo() {
	repoName := askUserForInput("repo name")
	repoDescription := askUserForInput("description")
	templateOwner := askUserForInput("template owner")
	templateRepo := askUserForInput("template repo name")
	authToken := goDotEnvVar("AUTH_KEY")

	values := map[string]string{"name": repoName, "description": repoDescription}
	json_data, err := json.Marshal(values)
	authroization := "token " + authToken
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	fmt.Println("Preparing to create repo")

	post_url := fmt.Sprintf("https://api.github.com/repos/%s/%s", templateOwner, templateRepo)
	req, err := http.NewRequest("POST", post_url, bytes.NewBuffer(json_data))
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", authroization)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Something went wrong")
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	var githubRepo GithubRepoRequirements
	json.Unmarshal(body, &githubRepo)
	fmt.Println(githubRepo.Name)
	fmt.Println("Repo Created! Happy Hacking.")
	fmt.Println("")
	fmt.Printf("Name: %s\nDescription: %s\nUrl: %s\nHomepage: %s\n", githubRepo.Name, githubRepo.Description, githubRepo.Url, githubRepo.Homepage)
	fmt.Println("")

}
