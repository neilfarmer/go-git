# Go Git

Go Git is a CLI tool that authenticates with GitHub or GitLab using a personal access token, then recursively fetches all accessible repositories (public and private) and clones them in a structured directory layout. It supports concurrent cloning for improved speed and structured logging with support for configurable verbosity.

## 📋 Prerequisites

Before using this tool, ensure the following are installed and configured:

| Requirement | Description                                                  |
| ----------- | ------------------------------------------------------------ |
| Go          | Go 1.18 or higher installed                                  |
| Git         | Git must be installed and available in your system's PATH    |
| Token       | GitHub or GitLab token with access to repositories           |
| Config File | A JSON config file located at `~/.config/go-git/config.json` |

## 🧩 Installation

Download the binary from the [Releases page](https://github.com/youruser/go-git/releases), or install via `curl`:

````bash
# For Linux AMD64
curl -L https://github.com/youruser/go-git/releases/latest/download/go-git-linux-amd64 -o /usr/local/bin/go-git
chmod +x /usr/local/bin/go-git

# For macOS ARM64 (Apple Silicon)
curl -L https://github.com/youruser/go-git/releases/latest/download/go-git-darwin-arm64 -o /usr/local/bin/go-git
chmod +x /usr/local/bin/go-git
```

### GitHub Token Setup

1. Visit https://github.com/settings/tokens?type=beta
2. Click **"Generate new token"**, select **"Fine-grained token"**
3. Give it access to:
   - Repositories (All Repositories)
   - Contents (Read-only)
   - Metadata (Read-only)
4. Save the token and place it in your config.

### GitLab Token Setup

1. Visit https://gitlab.com/-/profile/personal_access_tokens
2. Generate a token with:
   - `read_api`
   - `read_repository`
3. Save the token and place it in your config.

### Example `~/.config/go-git/config.json`

```json
{
  "token": "ghp_xxx...xxx",
  "scm_name": "github"
}
````

or for GitLab:

```json
{
  "token": "glpat_xxx...xxx",
  "scm_name": "gitlab"
}
```

## 🚀 Usage

Run from the command line:

```bash
go run main.go graph
```

### Example Commands

```bash
# Clone all accessible GitHub repositories with debug logging
go run main.go graph --verbose debug

# Synchronize repositories (fetch and prune)
go run main.go sync --verbose info

# Build the binary and run it
go build -o go-git
./go-git graph
```

## 🧭 Commands

| Command | Description                                                                |
| ------- | -------------------------------------------------------------------------- |
| `graph` | Display a tree of GitLab/GitHub groups and projects, with optional cloning |
| `sync`  | Clone or pull all repositories into local folders in parallel              |

## ⚙️ CLI Flags

| Flag        | Shorthand | Description                                     |
| ----------- | --------- | ----------------------------------------------- |
| `--verbose` | `-v`      | Logging level: `debug`, `info`, `warn`, `error` |

> 🔹 Note: The configuration file is expected at `~/.config/go-git/config.json` and is loaded automatically. It must contain a personal access token and a `scm_name` value of `github` or `gitlab`.

## 🛠 Features

- Authenticates with GitHub or GitLab via PAT
- Lists all repositories accessible to the token
- Clones all repositories concurrently with throttling
- Organized output using ASCII tree
- Verbose logging with `--verbose debug`

## 🔒 Security Warning

This tool uses the token directly in the clone URL. Do not expose logs or URLs containing the token. Future versions may include credential helper integration.

## Release

Releases are handled via tags

```bash
export APP_TAG="v1.x.x"

git tag ${APP_TAG}
git push origin ${APP_TAG}
```
