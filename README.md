# Gitwo - Git Worktree Helper CLI

A powerful command-line tool that extends `git worktree` functionality with advanced features, auto-switching, shell integration, and language/framework detection.

## 🚀 Features

- **Smart Worktree Management**: Create, list, and remove worktrees with intelligent branch naming
- **Auto-Switch**: Automatically change directory to new worktrees after creation
- **Shell Integration**: Cross-shell support (bash, zsh, fish, PowerShell) with auto-cd functionality
- **Language/Framework Detection**: Automatic detection of project types (Rails, Node.js, Go, Java, etc.)
- **Hook System**: Pre and post-execution hooks for custom automation
- **Configuration Management**: Flexible YAML-based configuration with environment variable overrides
- **Progress Display**: Beautiful, real-time progress indicators
- **Shell Completions**: Full completion support for all major shells

## 📦 Installation

### From Source
```bash
git clone https://github.com/gitwohq/gitwo.git
cd gitwo
make build
# Add to PATH: export PATH="$PWD/bin:$PATH"
```

### Homebrew (Coming Soon)
```bash
brew install gitwohq/gitwo/gitwo
```

## 🎯 Quick Start

### Basic Usage
```bash
# Create a new worktree
gitwo new feature-add-openapi

# List all worktrees
gitwo list

# Remove a worktree
gitwo remove feature-add-openapi
# or
gitwo rm feature-add-openapi
```

### Auto-Switch Setup
Enable automatic directory switching after creating worktrees:

```bash
# Install shell wrapper (one-time setup)
gitwo shell-install

# Or for current session only
eval "$(gitwo shell-init)"  # bash/zsh
gitwo shell-init --shell fish | source  # fish
gitwo shell-init --shell pwsh | Out-String | Invoke-Expression  # PowerShell
```

After setup, `gitwo new feature-name` will automatically `cd` into the new worktree!

## Release workflow (TL;DR)

**Local dev**
```bash
make check
make snapshot-fast   # quick dry-run (no docker/sign/sbom)
make snapshot        # full dry-run (no publish)
```

## 📋 Commands

### Core Commands

#### `gitwo new <name>`
Create a new worktree with smart branch naming.

```bash
gitwo new feature-add-openapi
# Creates: ../feature-add-openapi with branch feature/feature-add-openapi
```

**Flags:**
- `--start-point <ref>`: Specify start point (default: HEAD)
- `--shell`: Output only the cd command for shell integration
- `--switch`: Auto-switch to new worktree (if not using shell wrapper)

#### `gitwo list`
List all worktrees with detailed information.

```bash
gitwo list
# Output:
# /path/to/repo (main) [bare]
# /path/to/repo/../feature-add-openapi (feature/feature-add-openapi)
```

#### `gitwo remove <worktree>`
Remove a worktree by name or path.

```bash
gitwo remove feature-add-openapi
gitwo remove ../feature-add-openapi
gitwo rm feature-add-openapi  # alias
```

### Shell Integration

#### `gitwo shell-init [--shell <shell>]`
Print shell wrapper for auto-cd functionality.

```bash
gitwo shell-init                    # Auto-detect shell
gitwo shell-init --shell bash       # Specify shell
gitwo shell-init --shell zsh
gitwo shell-init --shell fish
gitwo shell-init --shell powershell
```

#### `gitwo shell-install [--shell <shell>]`
Install shell wrapper into your shell profile.

```bash
gitwo shell-install                 # Auto-detect and install
gitwo shell-install --shell zsh     # Install for specific shell
```

### Configuration

#### `gitwo config`
Manage gitwo configuration.

```bash
gitwo config                        # Show current config
gitwo config auto-switch true       # Enable auto-switch
gitwo config auto-switch false      # Disable auto-switch
```

## ⚙️ Configuration

Gitwo uses a YAML configuration file at `.gitwo/config.yml`:

```yaml
# Core settings
worktrees_dir: ".."                 # Where to place worktrees
name_template: "${REPO}-${BRANCH}"  # Directory name template
main_branch: "origin/main"          # Default base branch
editor_cmd: "code -g"               # Editor command
post_add_open_editor: true          # Auto-open editor after creation

# Auto-switch configuration
auto_switch: true
default_branch_prefix: "feature/"

# Language/Framework detection
language: "rails"  # auto-detected
framework: "rails"  # auto-detected

# Shell configuration
shell:
  type: "zsh"  # auto-detected
  auto_cd: true
  wrapper_installed: true

# Hook configuration
hooks:
  enabled: true
  pre_add:
    - type: "command"
      command: "bundle install"
      description: "Install Ruby dependencies"
  post_add:
    - type: "command"
      command: "bin/setup"
      description: "Run Rails setup script"
```

## 🔧 Environment Variables

Gitwo supports environment variable overrides:

```bash
export GITWO_EDITOR="code -g"           # Override editor command
export GITWO_OPEN="true"                 # Auto-open editor after new
export GITWO_SYNC_ON_NEW="true"          # Run git fetch before creating worktree
export GITWO_DOCKER_UP="true"            # Hint for hooks to bring up Docker
export GITWO_VERBOSE="1"                 # Logging verbosity (0|1|2)
export GITWO_COLOR="auto"                # Color mode (auto|always|never)
export GITWO_TIMEOUT="30"                # Hook execution timeout (seconds)
```

## 🪝 Hook System

Gitwo supports pre and post-execution hooks for automation:

### Hook Files
Create `.gitwo/hooks/pre_add.yml` or `.gitwo/hooks/post_add.yml`:

```yaml
hooks:
  - type: "command"
    command: "bundle install"
    description: "Install Ruby dependencies"
  - type: "command"
    command: "bin/setup"
    description: "Run Rails setup script"
```

### Hook Environment Variables
Hooks have access to these environment variables:

```bash
GITWO_ACTION="add"                    # Current action
GITWO_REPO="my-project"               # Repository name
GITWO_BRANCH="feature/new-feature"    # Branch name
GITWO_PATH="../new-feature"           # Worktree path
GITWO_WORKTREES_DIR=".."              # Worktrees directory
GITWO_MAIN_BRANCH="origin/main"       # Main branch
GITWO_NAME_TEMPLATE="${REPO}-${BRANCH}" # Name template
GITWO_EDITOR="code -g"                # Editor command
```

## 🎨 Language/Framework Detection

Gitwo automatically detects project types and applies appropriate defaults:

### Supported Languages
- **Ruby/Rails**: Detects `Gemfile`, `config/application.rb`, `bin/rails`
- **Node.js/Next.js**: Detects `package.json`, `next.config.js`
- **Go**: Detects `go.mod`
- **Java**: Detects `pom.xml`, `build.gradle`
- **Python**: Detects `requirements.txt`, `setup.py`, `pyproject.toml`
- **PHP**: Detects `composer.json`
- **Rust**: Detects `Cargo.toml`

### Framework Detection
- **Rails**: `config/application.rb` + `bin/rails`
- **Next.js**: `next.config.js` + `package.json` with `next` dependency
- **React**: `package.json` with `react` dependency
- **Vue**: `package.json` with `vue` dependency
- **Express**: `package.json` with `express` dependency
- **Gin**: `go.mod` with `gin-gonic/gin`
- **Spring Boot**: `pom.xml` with `spring-boot` dependency

## 🐚 Shell Completions

Install shell completions for better UX:

### Bash
```bash
# Add to ~/.bashrc
source <(gitwo completion bash)
```

### Zsh
```bash
# Add to ~/.zshrc
source <(gitwo completion zsh)
```

### Fish
```bash
# Add to ~/.config/fish/config.fish
gitwo completion fish | source
```

### PowerShell
```bash
# Add to PowerShell profile
gitwo completion powershell | Out-String | Invoke-Expression
```

## 🔄 Auto-Switch Examples

### With Shell Wrapper (Recommended)
```bash
# Install wrapper once
gitwo shell-install

# Now every gitwo new automatically switches
gitwo new feature-add-openapi
# ✅ Automatically cd's into ../feature-add-openapi
```

### Without Shell Wrapper
```bash
# Use --switch flag
gitwo new feature-add-openapi --switch

# Or use eval
eval "$(gitwo new feature-add-openapi --shell | tail -1)"
```

## 🏗️ Project Structure

```
gitwo/
├── bin/gitwo                    # Compiled binary
├── cmd/gitwo/                   # CLI commands
├── internal/                    # Internal packages
│   ├── wt/                     # Worktree management
│   ├── shell/                  # Shell integration
│   ├── detection/              # Language/framework detection
│   ├── hooks/                  # Hook system
│   └── config/                 # Configuration management
├── completions/                # Shell completion files
├── templates/                  # Template files
└── README.md
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `gitwo new feature-name`
3. Make your changes
4. Add tests for new functionality
5. Run tests: `make test`
6. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [worktree_manager](https://github.com/nacyot/worktree_manager)
- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Uses [testify](https://github.com/stretchr/testify) for testing

---

**Happy coding with Gitwo! 🚀**
