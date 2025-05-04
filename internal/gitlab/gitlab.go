package gitlab

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/neilfarmer/go-git/internal/config"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)


func GetRepos(config config.Config) {
	client, err := SetupClient(config)
	if err != nil {
		log.Fatalf("Failed to set up GitLab client: %v", err)
	}

	os.Mkdir("gitlab", 0755)
	os.Chdir("gitlab")

	groups, err := GetGroups(client)
	if err != nil {
		log.Fatalf("Failed to get groups: %v", err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // limit concurrent clones

	for _, group := range groups {
		groupName := group.FullPath
		slog.Debug("Group Path", "groupName", groupName)
		if err := os.MkdirAll(groupName, 0755); err != nil {
			log.Fatalf("Failed to create directory for group %s: %v", groupName, err)
		}

		projects, _, err := client.Groups.ListGroupProjects(group.ID, &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		})
		if err != nil {
			log.Fatalf("Failed to list projects for group %s: %v", groupName, err)
		}

		for _, project := range projects {
			targetDir := filepath.Join(groupName, project.Name)
			if _, err := os.Stat(targetDir); os.IsNotExist(err) {
				wg.Add(1)
				sem <- struct{}{}
				go func(project *gitlab.Project, targetDir string) {
					defer wg.Done()
					defer func() { <-sem }()
					cloneRepo(config, project, targetDir)
				}(project, targetDir)
			}
		}
	}

	wg.Wait()
}

func cloneRepo(config config.Config, project *gitlab.Project, targetDir string) {
	repoUrl := strings.Replace(project.HTTPURLToRepo, "https://", fmt.Sprintf("https://oauth2:%s@", config.Token), 1)

	slog.Debug("Project Name", "name", project.Name)
	slog.Debug("Clone URL", "cloneUrl", repoUrl)
	slog.Debug("Target Directory", "targetDir", targetDir)

	cmd := exec.Command("git", "clone", repoUrl, targetDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error cloning repository %s: %v\n", project.Name, err)
	} else {
		log.Printf("Successfully cloned repository %s\n", project.Name)
	}
}

func GetGroups(client *gitlab.Client) ([]gitlab.Group, error) {
	groups, _, err := client.Groups.ListGroups(&gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	})
	if err != nil {
		log.Fatalf("Failed to list groups: %v", err)
	}

	var groupList []gitlab.Group
	for _, group := range groups {
		groupList = append(groupList, *group)
	}
	return groupList, nil
}

func SetupClient(config config.Config) (*gitlab.Client, error) {
	if config.Url == "" {
		slog.Debug("GitLab URL not set, using default")
		config.Url = "https://gitlab.com"
	}
	client, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.Url+"/api/v4"))
	slog.Debug("GitLab Url", "url", config.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}
	return client, nil
}

func GraphRepos(config config.Config) {
	client, err := SetupClient(config)
	if err != nil {
		log.Fatalf("Failed to set up client: %v", err)
	}

	groups, err := GetGroups(client)
	if err != nil {
		log.Fatalf("Failed to get groups: %v", err)
	}

	// Map by full path for parent tracking
	groupMap := make(map[string]*gitlab.Group)
	for _, group := range groups {
		groupMap[group.FullPath] = &group
	}

	printGroupTree(client, groupMap, "", "")
}

func printGroupTree(client *gitlab.Client, groupMap map[string]*gitlab.Group, parentPath string, prefix string) {
	children := getChildPaths(groupMap, parentPath)

	for i, fullPath := range children {
		group := groupMap[fullPath]
		isLast := i == len(children)-1

		branch := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			branch = "└── "
			nextPrefix = prefix + "    "
		}

		fmt.Printf("%s%s%s\n", prefix, branch, group.Name)
		printProjects(client, group.ID, nextPrefix)
		printGroupTree(client, groupMap, fullPath, nextPrefix)
	}
}

func getChildPaths(groupMap map[string]*gitlab.Group, parentPath string) []string {
	var children []string
	for fullPath := range groupMap {
		if getParentPath(fullPath) == parentPath {
			children = append(children, fullPath)
		}
	}
	return children
}

func getParentPath(fullPath string) string {
	if idx := strings.LastIndex(fullPath, "/"); idx != -1 {
		return fullPath[:idx]
	}
	return ""
}

func printProjects(client *gitlab.Client, groupID int, prefix string) {
	projects, _, err := client.Groups.ListGroupProjects(groupID, &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100, Page: 1},
	})
	if err != nil {
		log.Printf("Failed to list projects for group ID %d: %v", groupID, err)
		return
	}

	for i, project := range projects {
		isLast := i == len(projects)-1
		branch := "├── "
		if isLast {
			branch = "└── "
		}
		fmt.Printf("%s%s%s\n", prefix, branch, project.Name)
	}
}