package github

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v71/github"
	"github.com/neilfarmer/go-git/internal/config"
	"golang.org/x/oauth2"
)

func GetRepos(config config.Config) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repos, _, err := client.Repositories.ListByAuthenticatedUser(ctx, nil)
	if err != nil {
		fmt.Println("Error listing repositories:", err)
		return
	}
	os.Mkdir("github", 0755)
	os.Chdir("github")

	for _, repo := range repos {
		fmt.Printf("Cloning Repository: %s\n", *repo.Name)
		fmt.Printf("Clone URL: %s\n", *repo.CloneURL)
		cloneURL := strings.Replace(*repo.CloneURL, "https://", fmt.Sprintf("https://%s@", config.Token), 1)

		cmd := exec.Command("git", "clone", cloneURL)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error cloning repository %s: %v\n", *repo.Name, err)
		} else {
			fmt.Printf("Successfully cloned repository %s\n", *repo.Name)
		}
		time.Sleep(5 * time.Second) // Sleep for 5 seconds between clones
	}
}