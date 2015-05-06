package src

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

func init() {
	repoGroup, err := CLI.AddCommand("repo",
		"describe the current repository",
		"The repo command displays information the current repository. If there is a remote repository, its information is also displayed.",
		&repoCmd,
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = repoGroup
}

type RepoCmd struct{}

var repoCmd RepoCmd

func (c *RepoCmd) Execute(args []string) error {
	if localRepo != nil {
		fmt.Println("# Local repository:")
		fmt.Printf("Root directory:\t%s\n", localRepo.RootDir)
		fmt.Printf("VCS type:\t%s\n", localRepo.VCSType)
		fmt.Printf("Commit ID:\t%s\n", localRepo.CommitID)
		fmt.Printf("Clone URL:\t%s\n", localRepo.CloneURL)

		fmt.Println()

		if localRepoErr != nil {
			fmt.Printf("# Warning: %s\n", localRepoErr)
			fmt.Println("# Not trying to fetch and display remote repository information due to the above error.")
		} else {
			cl := NewAPIClientWithAuthIfPresent()
			remoteRepo, err := getRemoteRepo(cl)
			if err == nil {
				fmt.Println("# Remote repository:")
				printRemoteRepo(remoteRepo)
			} else if sourcegraph.IsHTTPErrorCode(err, http.StatusNotFound) {
				fmt.Println("# No remote repository found.")
				fmt.Println("# Use `src remote add` to create it.")
			} else {
				fmt.Printf("# Error getting remote repository: %s.\n", err)
			}
		}
	} else {
		fmt.Println("# No local git/hg repository found in or above the current directory.")
		if localRepoErr != nil {
			fmt.Printf("# Error was: %s.\n", localRepoErr)
		}
	}
	return nil
}

func printRemoteRepo(repo *sourcegraph.Repo) {
	fmt.Printf("URI:\t\t%s\n", repo.URI)
	if repo.URIAlias != "" {
		fmt.Printf("URI alias:\t%s\n", repo.URIAlias)
	}

	if GlobalOpt.Verbose {
		if repo.Description != "" {
			fmt.Printf("Description:\t%s\n", repo.Description)
		}
	}

	if repo.HomepageURL != "" {
		fmt.Printf("Homepage:\t%s\n", repo.HomepageURL)
	}
	if repo.HTTPCloneURL != "" {
		fmt.Printf("Clone (HTTP):\t%s (%s)\n", repo.HTTPCloneURL, repo.VCS)
	}
	if repo.SSHCloneURL != "" {
		fmt.Printf("Clone (SSH):\t%s (%s)\n", repo.SSHCloneURL, repo.VCS)
	}
	fmt.Printf("Default branch:\t%s\n", repo.DefaultBranch)
	if repo.Language != "" {
		fmt.Printf("Language:\t%s\n", repo.Language)
	}
	if repo.Private {
		fmt.Printf("Private:\tyes\n")
	}
	fmt.Printf("Created at:\t%s\n", repo.CreatedAt)
	fmt.Printf("Updated at:\t%s\n", repo.UpdatedAt)
	fmt.Printf("Pushed at:\t%s\n", repo.PushedAt)
}
