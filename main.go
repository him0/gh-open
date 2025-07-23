package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/mattn/go-isatty"
)

type Remote struct {
	Name string
	URL  string
}

var (
	green      = "\033[32m"
	lightGreen = "\033[32;1m"
	red        = "\033[31m"
	lightRed   = "\033[31;1m"
	magenta    = "\033[35m"
	resetColor = "\033[0m"
	verbose    bool
	colorFlag  string
	listFlag   bool
)

func main() {
	// Parse command line flags
	flag.BoolVar(&verbose, "verbose", false, "verbose output")
	flag.StringVar(&colorFlag, "color", "auto", "colorize output (always, never, auto)")
	flag.BoolVar(&listFlag, "list", false, "open pull requests list instead of new PR")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fmt.Fprintf(os.Stderr, "  --verbose\n        verbose output\n")
		fmt.Fprintf(os.Stderr, "  --color string\n        colorize output (always, never, auto) (default \"auto\")\n")
		fmt.Fprintf(os.Stderr, "  --list\n        open pull requests list instead of new PR\n")
	}
	flag.Parse()

	// Set up color output
	colorize := colorizeOutput()
	if !colorize {
		green = ""
		lightGreen = ""
		red = ""
		lightRed = ""
		magenta = ""
		resetColor = ""
	}

	// Check if current directory is a git repository
	if err := checkGitRepo(); err != nil {
		exitWithError("fatal: Not a git repository")
	}

	// Get main remote (first available in priority order: upstream, github, origin)
	remote, err := getMainRemote()
	if err != nil {
		exitWithError(err.Error())
	}

	// Convert remote URL to GitHub URL
	repoURL, err := extractGitHubURL(remote.URL)
	if err != nil {
		exitWithError("Failed to extract GitHub repository URL: %s", err.Error())
	}

	var url string
	if listFlag {
		// Generate pull requests list URL
		url = fmt.Sprintf("https://github.com/%s/pulls", repoURL)
	} else {
		// Get current branch
		currentBranch, err := getCurrentBranch()
		if err != nil {
			exitWithError("Failed to get current branch: %s", err.Error())
		}

		// Generate new pull request URL
		url = fmt.Sprintf("https://github.com/%s/pull/new/%s", repoURL, currentBranch)
	}

	// Display URL
	fmt.Printf("%s%s%s\n", lightGreen, url, resetColor)

	// Open browser if available
	if err := openBrowser(url); err != nil {
		verboseLog("Failed to open browser", []string{err.Error()})
	}
}

func checkGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func getMainRemote() (*Remote, error) {
	// Priority order: upstream, github, origin, others
	priorityOrder := []string{"upstream", "github", "origin"}

	remotes, err := getRemotes()
	if err != nil {
		return nil, err
	}

	if len(remotes) == 0 {
		return nil, fmt.Errorf("no git remotes found")
	}

	// Check priority remotes first
	for _, priority := range priorityOrder {
		for _, remote := range remotes {
			if remote.Name == priority {
				return &remote, nil
			}
		}
	}

	// Return first remote if no priority match
	return &remotes[0], nil
}

func getRemotes() ([]Remote, error) {
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	remoteMap := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 && strings.HasSuffix(line, "(fetch)") {
			remoteMap[parts[0]] = parts[1]
		}
	}

	var remotes []Remote
	for name, url := range remoteMap {
		remotes = append(remotes, Remote{Name: name, URL: url})
	}

	return remotes, nil
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func extractGitHubURL(remoteURL string) (string, error) {
	// Remove .git suffix if present
	url := strings.TrimSuffix(remoteURL, ".git")

	// Handle SSH URLs (git@github.com:owner/repo)
	sshRegex := regexp.MustCompile(`^git@github\.com:(.+/.+)$`)
	if matches := sshRegex.FindStringSubmatch(url); len(matches) > 1 {
		return matches[1], nil
	}

	// Handle HTTPS URLs (https://github.com/owner/repo)
	httpsRegex := regexp.MustCompile(`^https://github\.com/(.+/.+)$`)
	if matches := httpsRegex.FindStringSubmatch(url); len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("not a GitHub repository URL: %s", remoteURL)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	verboseLog(cmd.Path, cmd.Args[1:])
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Don't wait for the browser to exit
	go cmd.Wait()
	return nil
}

func verboseLog(cmd string, args []string) {
	if verbose {
		msg := fmt.Sprintf("$ %s %s", cmd, strings.Join(args, " "))
		if isatty.IsTerminal(os.Stderr.Fd()) && magenta != "" {
			msg = fmt.Sprintf("%s%s%s", magenta, msg, resetColor)
		}
		fmt.Fprintln(os.Stderr, msg)
	}
}

func colorizeOutput() bool {
	switch colorFlag {
	case "always":
		return true
	case "never":
		return false
	case "auto":
		return isatty.IsTerminal(os.Stdout.Fd())
	default:
		return isatty.IsTerminal(os.Stdout.Fd())
	}
}

func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", lightRed, fmt.Sprintf(format, args...), resetColor)
	os.Exit(1)
}