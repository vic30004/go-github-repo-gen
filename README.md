# Github Repo Generator

This is a small application that provides you with the option to look up users on github, generate new repos from templates, or to create a new repo from scratch.

## How to use

1. Fork the repo
2. Clone the repo locally
3. Go to repo and add a .env file
4. Generate a personal access token from github. Follow the instruction in this [link](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token).
4. Add an env variable inside the .env called AUTH_KEY and assign its value to the access token. 
5. Add another env variable called LOCALPATH and set it to the path where you would like folders to be created. Example `/Users/<name>/Desktop/Repositories/`
6. Run `go run main.go`