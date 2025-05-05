package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func BuildRepoURL(protocol, host, user, repo string) string {
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

func GetHostByProvider(provider string) string {
	// Get host by provider
	host, ok := providerMap[strings.ToLower(provider)]
	if !ok {
		log.Fatalf("Unknown provider: %s", provider)
	}

	return host
}

func CheckGitInstalled() bool {
	// Check for git
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal("git is not installed or not in PATH")

		return false
	}

	return true
}

func ValidateGitProvider(provider string) bool {
	// Validate provider
	_, ok := providerMap[strings.ToLower(provider)]
	if !ok {
		log.Fatalf("Unknown provider: %s", provider)

		return false
	}

	return true
}

func GitClone(repoURL string, outputDir string) bool {
	// git clone --no-checkout ...
	cloneCmd := exec.Command("git", "clone", "--no-checkout", repoURL, outputDir)
	cloneCmd.Stdout, cloneCmd.Stderr = os.Stdout, os.Stderr
	if err := cloneCmd.Run(); err != nil {
		log.Fatalf("git clone failed: %v", err)

		return false
	}

	return true
}

func GitSparseCheckoutInit() bool {
	// git sparse-checkout init --cone
	initCmd := exec.Command("git", "sparse-checkout", "init", "--cone")

	initCmd.Stdout, initCmd.Stderr = os.Stdout, os.Stderr

	if err := initCmd.Run(); err != nil {
		log.Fatalf("git sparse-checkout init failed: %v", err)

		return false
	}

	return true
}

func GitSparseCheckoutPaths(paths []string) bool {
	// git sparse-checkout set <paths...>
	setCmd := exec.Command("git", append([]string{"sparse-checkout", "set"}, paths...)...)

	setCmd.Stdout, setCmd.Stderr = os.Stdout, os.Stderr

	if err := setCmd.Run(); err != nil {
		log.Fatalf("git sparse-checkout set failed: %v", err)

		return false
	}

	return true
}

func GitCheckoutBranch(branch string) bool {
	coCmd := exec.Command("git", "checkout", branch)
	coCmd.Stdout, coCmd.Stderr = os.Stdout, os.Stderr
	if err := coCmd.Run(); err != nil {
		log.Fatalf("git checkout failed: %v", err)

		return false
	}

	return true
}
