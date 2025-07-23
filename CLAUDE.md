# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gh-open is a GitHub CLI extension that quickly opens GitHub Pull Request pages in your browser. It automatically detects repository information and generates the appropriate GitHub URLs for new PR creation or PR list viewing.

## Build and Development Commands

```bash
# Build the binary
go build -o gh-open .

# Run go mod tidy to update dependencies
go mod tidy

# Test locally (from the gh-open directory)
./gh-open

# Test with verbose output
./gh-open --verbose

# Test with color options
./gh-open --color=always
./gh-open --color=never
./gh-open --color=auto

# Test list functionality
./gh-open --list
```

## Command Line Options

- `--verbose`: Show commands being executed
- `--color[=WHEN]`: Control color output
  - `always`: Always colorize output
  - `never`: Never colorize output  
  - `auto`: Colorize output only when stdout is a terminal (default)
- `--list`: Open pull requests list instead of new PR page

## Architecture

The entire application is contained in a single `main.go` file with the following key functions:

- `main()`: Entry point that orchestrates the workflow and parses command line flags
- `getMainRemote()`: Finds the main remote (priority: upstream, github, origin)
- `extractGitHubURL()`: Converts git remote URLs to GitHub repository format
- `getCurrentBranch()`: Gets the current git branch name
- `openBrowser()`: Opens URLs in the default browser across platforms
- `verboseLog()`: Logs commands when --verbose flag is enabled
- `colorizeOutput()`: Determines if color output should be enabled based on --color flag

## Key Implementation Details

1. **Command Line Parsing**: Uses Go's `flag` package for --verbose, --color, and --list options
2. **Verbose Logging**: Shows commands in `$ command args` format with magenta color
3. **Color Output**: Uses ANSI escape codes with isatty detection for terminal-aware coloring
4. **URL Generation**: Supports both new PR and PR list URLs based on --list flag
5. **Cross-platform Browser Opening**: Uses platform-specific commands (open, xdg-open, rundll32)
6. **Error Handling**: All errors exit with status code 1 and colored error messages to stderr

## GitHub URL Formats

- **New PR**: `https://github.com/owner/repo/pull/new/branch`
- **PR List**: `https://github.com/owner/repo/pulls`

## Release Process

GitHub Actions workflow (`.github/workflows/release.yml`) builds cross-platform binaries when a version tag is pushed:
- Darwin (amd64, arm64)
- Linux (amd64, arm64)  
- Windows (amd64)

To release a new version:
```bash
git tag v1.0.0
git push origin v1.0.0
```