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
	"os/exec"
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

func createReadme(name string) {
	readmeName:=fmt.Sprintf("# %s", name)
	cmd := exec.Command("echo", readmeName)

	// Make test file
	testFile, err := os.Create("README.md")
	if err != nil {
		panic(err)
	}
	defer testFile.Close()

	// Redirect the output here (this is the key part)
	cmd.Stdout = testFile

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.Wait()
}

func commitChanges() {
	commit := exec.Command("git", "commit", "-m", "'Initial commit'")
	var out bytes.Buffer
	var stderr bytes.Buffer
	commit.Stdout = &out
	commit.Stderr = &stderr
	err := commit.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Println("Result: " + out.String())
	fmt.Println("Initial Commit made")
}

func initRepo(name string, url string) {
	localPath := goDotEnvVar("LOCALPATH")
	dir := fmt.Sprintf("%s/%s", localPath, name)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Folder created")
	os.Chdir(dir)

	init := exec.Command("git", "init").Run()
	if init != nil {
		log.Fatal(init)
	}
	fmt.Println("Git init complete")

	// create readme
	createReadme(name)

	// add changes
	addChanges := exec.Command("git", "add", ".").Run()
	if addChanges != nil {
		log.Fatal(addChanges)
	}
	fmt.Println("Readme added")

	// commit changes
	commitChanges()

	// make sure we are on main branch
	branch := exec.Command("git", "branch", "-M", "main").Run()
	if branch != nil {
		log.Fatal(branch)
	}

	//Connect to repo
	repoAddress := fmt.Sprintf("%s.git", url)
	remote := exec.Command("git", "remote", "add", "origin", repoAddress).Run()
	if remote != nil {
		log.Fatal(remote.Error())
	}
	fmt.Println("Connection to repo complete")

	// push changes
	push := exec.Command("git", "push", "-u", "origin", "main").Run()
	if push != nil {
		log.Fatal(push.Error())
	}
	fmt.Println("Changes pushed to github!")
}

func makePostRequest(url string, json_data []byte) GithubRepoRequirements {
	authToken := goDotEnvVar("AUTH_KEY")
	authroization := "token " + authToken

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
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
	return githubRepo
}

func CreateRepoFromTemplate() {
	repoName := askUserForInput("repo name")
	repoDescription := askUserForInput("description")
	templateOwner := askUserForInput("template owner")
	templateRepo := askUserForInput("template repo name")

	values := map[string]string{"name": repoName, "description": repoDescription}
	json_data, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}
	fmt.Println("Preparing to create repo")

	post_url := fmt.Sprintf("https://api.github.com/repos/%s/%s", templateOwner, templateRepo)

	res := makePostRequest(post_url, json_data)

	fmt.Println(res.Name)
	fmt.Println("Repo Created! Happy Hacking.")
	fmt.Println("")
	fmt.Printf("Name: %s\nDescription: %s\nUrl: %s\nHomepage: %s\n", res.Name, res.Description, res.Url, res.Homepage)
	fmt.Println("")
}

func CreateNewRepo() {
	repoName := askUserForInput("repo name")
	repoDescription := askUserForInput("description")
	values := map[string]string{"name": repoName, "description": repoDescription}

	json_data, err := json.Marshal(values)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	fmt.Println("Preparing to create repo")

	post_url := fmt.Sprintf("https://api.github.com/user/repos")
	res := makePostRequest(post_url, json_data)

	fmt.Println("Repo Created!")
	fmt.Println("")
	fmt.Printf("Name: %s\nDescription: %s\nUrl: %s\nHomepage: %s\n", res.Name, res.Description, res.Url, res.Homepage)
	fmt.Println("")
	fmt.Println("Initializing repo")

	initRepo(res.Name, res.Url)
}
