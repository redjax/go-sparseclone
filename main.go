package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
)

var providerMap = map[string]string{
	"github":   "github.com",
	"gitlab":   "gitlab.com",
	"codeberg": "codeberg.org",
}

type CLI struct {
	Provider string   `required:"" long:"provider" help:"Git provider: github, gitlab, codeberg." default:"github"`
	User     string   `required:"" short:"u" long:"username" help:"Git username or org."`
	Repo     string   `required:"" short:"r" long:"repository" help:"Repository name."`
	Output   string   `short:"o" long:"output-dir" help:"Output directory. Defaults to repo name."`
	Branch   string   `short:"b" long:"checkout-branch" default:"main" help:"Branch name to checkout."`
	Paths    []string `required:"" short:"p" long:"checkout-path" name:"path" help:"Paths to sparse-checkout (repeatable)."`
	Protocol string   `long:"protocol" help:"Clone protocol: ssh or https." default:"ssh"`
}

func buildRepoURL(protocol, host, user, repo string) string {
	switch protocol {
	case "ssh":
		if !strings.HasSuffix(repo, ".git") {
			repo += ".git"
		}
		return fmt.Sprintf("git@%s:%s/%s", host, user, repo)
	case "https":
		return fmt.Sprintf("https://%s/%s/%s", host, user, repo)
	default:
		log.Fatalf("Unknown protocol: %s", protocol)
		return ""
	}
}

func main() {
	var cli CLI
	kong.Parse(&cli,
		kong.Name("sparseclone"),
		kong.Description("Clone a git repo with sparse checkout in one step."),
	)

	// Check for git
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal("git is not installed or not in PATH")
	}

	// Validate provider
	host, ok := providerMap[strings.ToLower(cli.Provider)]
	if !ok {
		log.Fatalf("Unknown provider: %s", cli.Provider)
	}

	// Determine output directory
	outputDir := cli.Output
	if outputDir == "" || outputDir == "." {
		repoName := cli.Repo
		if strings.HasSuffix(repoName, ".git") {
			repoName = strings.TrimSuffix(repoName, ".git")
		}
		outputDir = repoName
	}

	// Compose repo URL
	repoURL := buildRepoURL(cli.Protocol, host, cli.User, cli.Repo)
	fmt.Printf("Cloning from %s...\n", repoURL)

	// Step 1: git clone --no-checkout ...
	cloneCmd := exec.Command("git", "clone", "--no-checkout", repoURL, outputDir)
	cloneCmd.Stdout, cloneCmd.Stderr = os.Stdout, os.Stderr
	if err := cloneCmd.Run(); err != nil {
		log.Fatalf("git clone failed: %v", err)
	}

	// Step 2: cd <outputDir>
	absOutputDir, err := filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("Could not get absolute path: %v", err)
	}
	if _, err := os.Stat(absOutputDir); os.IsNotExist(err) {
		log.Fatalf("Output directory does not exist: %v", absOutputDir)
	}
	if err := os.Chdir(absOutputDir); err != nil {
		log.Fatalf("Could not enter output dir: %v", err)
	}

	// Step 3: git sparse-checkout init --cone
	initCmd := exec.Command("git", "sparse-checkout", "init", "--cone")
	initCmd.Stdout, initCmd.Stderr = os.Stdout, os.Stderr
	if err := initCmd.Run(); err != nil {
		log.Fatalf("git sparse-checkout init failed: %v", err)
	}

	// Step 4: git sparse-checkout set <paths...>
	setCmd := exec.Command("git", append([]string{"sparse-checkout", "set"}, cli.Paths...)...)
	setCmd.Stdout, setCmd.Stderr = os.Stdout, os.Stderr
	if err := setCmd.Run(); err != nil {
		log.Fatalf("git sparse-checkout set failed: %v", err)
	}

	// Step 5: git checkout <branch>
	coCmd := exec.Command("git", "checkout", cli.Branch)
	coCmd.Stdout, coCmd.Stderr = os.Stdout, os.Stderr
	if err := coCmd.Run(); err != nil {
		log.Fatalf("git checkout failed: %v", err)
	}

	fmt.Println("Sparse clone complete!")
}
