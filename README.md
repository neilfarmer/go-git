# Go Git

Go Git is a command-line tool that authenticates with a fine-grained GitHub token and clones all accessible repositories (public and private) for the user.

## 📋 Prerequisites

Before using this tool, ensure the following prerequisites are met:

| Requirement  | Description                                                    |
| ------------ | -------------------------------------------------------------- |
| Go           | Go 1.18 or higher installed                                    |
| Git          | Git must be installed and available in your system's PATH      |
| GitHub Token | Fine-grained personal access token with access to repositories |
| Config File  | A JSON config file located at `~/.config/go-git/config.json`   |

Example `config.json`:

```json
{
  "token": "ghp_xxx...xxx",
  "scm_name": "github"
}
```

## 🚀 Usage

Run the program from the command line:

```bash
go run main.go
```

It will:

- Load your GitHub token from the config file
- Authenticate with GitHub
- List all repositories accessible to your token
- Clone them into the current directory

## 🛡 Permissions Required

| Permission Type   | Permission Level | Purpose                                       |
| ----------------- | ---------------- | --------------------------------------------- |
| Repository Access | All Repositories | Grants access to all public and private repos |
| Contents          | Read-only        | Allows cloning and reading repository files   |
| Metadata          | Read-only        | Enables reading repository metadata           |

## 🔒 Security Warning

This application currently injects the GitHub token into the `git clone` URL. Do not share your terminal or logs where the token could be exposed. Future versions may add credential helper support for improved security.
