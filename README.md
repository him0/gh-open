# gh-open

GitHub CLI Extension to quickly open GitHub Pull Request pages in your browser.

## Features

- Automatically detects the main remote (upstream, github, origin priority)
- Opens new pull request page with current branch
- Can open pull requests list page with `--list` flag
- Automatically opens browser on supported platforms (macOS, Linux, Windows)
- Colored output for better visibility
- Verbose logging support to see commands being executed

## Installation

```bash
gh extension install him0/gh-open
```

### Upgrade

```bash
gh extension upgrade gh-open
```

## Usage

```bash
gh open [options]
```

### Options

- `--verbose`: Show commands being executed
- `--color[=WHEN]`: Control color output
  - `always`: Always colorize output
  - `never`: Never colorize output  
  - `auto`: Colorize output only when stdout is a terminal (default)
- `--list`: Open pull requests list instead of new PR page

### Examples

```bash
# Open new pull request page for current branch
gh open

# Open pull requests list
gh open --list

# Open with verbose output
gh open --verbose

# Open with color control
gh open --color=always
gh open --color=never

# Combine options
gh open --list --verbose --color=never
```

This will:
1. Check if the current directory is a git repository
2. Find the main remote (priority: upstream, github, origin)
3. Extract GitHub repository URL from remote
4. Generate appropriate GitHub URL:
   - New PR: `https://github.com/owner/repo/pull/new/branch`
   - PR list: `https://github.com/owner/repo/pulls`
5. Display the URL and open it in your default browser

### Example Output

```bash
$ gh open --verbose
https://github.com/owner/repo/pull/new/feature-branch
$ open https://github.com/owner/repo/pull/new/feature-branch
```

```bash
$ gh open --list
https://github.com/owner/repo/pulls
```

## Requirements

- [GitHub CLI](https://cli.github.com/) (gh) must be installed and authenticated
- Git must be installed
- The current directory must be a git repository with GitHub remotes
- Browser must be available for automatic opening (optional)

## Notes

- Supports SSH and HTTPS GitHub remote URLs
- Automatically opens browser on macOS (open), Linux (xdg-open), and Windows (rundll32)
- Only works with GitHub repositories (github.com)

## License

MIT