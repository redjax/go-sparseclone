package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/redjax/go-sparseclone/git"
)

type CLI struct {
	Provider string   `required:"" long:"provider" help:"Git provider: github, gitlab, codeberg." default:"github"`
	User     string   `required:"" short:"u" long:"username" help:"Git username or org."`
	Repo     string   `required:"" short:"r" long:"repository" help:"Repository name."`
	Output   string   `short:"o" long:"output-dir" help:"Output directory. Defaults to repo name."`
	Branch   string   `short:"b" long:"checkout-branch" default:"main" help:"Branch name to checkout."`
	Paths    []string `required:"" short:"p" long:"checkout-path" name:"path" help:"Paths to sparse-checkout (repeatable)."`
	Protocol string   `long:"protocol" help:"Clone protocol: ssh or https." default:"ssh"`
}

func main() {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("sparseclone"),
		kong.Description("Clone a git repo with sparse checkout in one step."),
	)

	// Check git is installed
	if !git.CheckGitInstalled() {
		log.Fatal("git is not installed or not in PATH")
	}

	// Validate provider input
	if !git.ValidateGitProvider(cli.Provider) {
		log.Fatalf("Unknown provider: %s", cli.Provider)
	}

	// Get output directory from repository's name
	outputDir := cli.Output
	if outputDir == "" || outputDir == "." {
		repoName := cli.Repo
		if strings.HasSuffix(repoName, ".git") {
			repoName = strings.TrimSuffix(repoName, ".git")
		}
		outputDir = repoName
	}

	// Build repo path from provider, protocol, user, and repo
	host := git.GetHostByProvider(cli.Provider)
	repoUrl := git.BuildRepoURL(cli.Protocol, host, cli.User, cli.Repo)

	// Clone repo without checking out
	if !git.GitClone(repoUrl, outputDir) {
		log.Fatalf("git clone failed")
	}

	// cd <outputDir>
	absOutputDir, err := filepath.Abs(outputDir)
	fmt.Printf("Output directory: %v\n", absOutputDir)
	if err != nil {
		log.Fatalf("Could not get absolute path: %v", err)
	}
	if _, err := os.Stat(absOutputDir); os.IsNotExist(err) {
		log.Fatalf("Output directory does not exist: %v", absOutputDir)
	}
	if err := os.Chdir(absOutputDir); err != nil {
		log.Fatalf("Could not enter output dir: %v", err)
	}

	// git sparse-checkout init --cone
	// initCmd := exec.Command("git", "sparse-checkout", "init", "--cone")
	// initCmd.Stdout, initCmd.Stderr = os.Stdout, os.Stderr
	// if err := initCmd.Run(); err != nil {
	// 	log.Fatalf("git sparse-checkout init failed: %v", err)
	// }
	if !git.GitSparseCheckoutInit() {
		log.Fatalf("git sparse-checkout init failed")
	}

	// git sparse-checkout set <paths...>
	// setCmd := exec.Command("git", append([]string{"sparse-checkout", "set"}, cli.Paths...)...)
	// setCmd.Stdout, setCmd.Stderr = os.Stdout, os.Stderr
	// if err := setCmd.Run(); err != nil {
	// 	log.Fatalf("git sparse-checkout set failed: %v", err)
	// }
	if !git.GitSparseCheckoutPaths(cli.Paths) {
		log.Fatalf("git sparse-checkout set failed")
	}

	// git checkout <branch>
	// coCmd := exec.Command("git", "checkout", cli.Branch)
	// coCmd.Stdout, coCmd.Stderr = os.Stdout, os.Stderr
	// if err := coCmd.Run(); err != nil {
	// 	log.Fatalf("git checkout failed: %v", err)
	// }

	if !git.GitCheckoutBranch(cli.Branch) {
		log.Fatalf("git checkout failed")
	}

	fmt.Println("Sparse clone complete!")
}
