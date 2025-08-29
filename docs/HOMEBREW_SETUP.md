# Homebrew Package Setup for gitwo

This guide explains how to set up and maintain the Homebrew package for gitwo.

## Prerequisites

1. **GitHub Repository**: `gitwohq/gitwo`
2. **Homebrew Tap Repository**: `gitwohq/homebrew-tap`
3. **GitHub Token**: With permissions to create releases and push to tap

## Step 1: Create Homebrew Tap Repository

Create a new repository: `gitwohq/homebrew-tap`

```bash
# Clone the tap repository
git clone https://github.com/gitwohq/homebrew-tap.git
cd homebrew-tap

# Create Formula directory
mkdir Formula
```

## Step 2: Configure GitHub Actions for Releases

Create `.github/workflows/release.yml` in the main gitwo repository:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

## Step 3: Set Up GitHub Secrets

In your GitHub repository settings, add these secrets:

- `GITHUB_TOKEN`: Automatically provided by GitHub
- `HOMEBREW_TAP_GITHUB_TOKEN`: Personal access token with repo permissions

## Step 4: Create a Release

```bash
# Tag a new version
git tag v0.1.0
git push origin v0.1.0
```

This will trigger:
1. GitHub Actions builds the binaries
2. Creates a GitHub release
3. Updates the Homebrew tap automatically

## Step 5: Test the Homebrew Package

```bash
# Add the tap
brew tap gitwohq/tap

# Install gitwo
brew install gitwo

# Test installation
gitwo --version
```

## Step 6: Manual Installation (for testing)

If you want to test locally:

```bash
# Build locally
make build

# Create a test formula
brew install --build-from-source Formula/gitwo.rb
```

## File Structure

```
gitwo/
├── .goreleaser.yaml          # Release configuration
├── Formula/
│   └── gitwo.rb             # Homebrew formula template
├── completions/              # Shell completions
│   ├── gitwo.bash
│   ├── gitwo.zsh
│   ├── gitwo.fish
│   └── gitwo.ps1
└── docs/
    └── HOMEBREW_SETUP.md    # This file
```

## Homebrew Formula Features

The formula includes:

- ✅ **Multi-platform support**: macOS (Intel/ARM) and Linux (Intel/ARM)
- ✅ **Shell completions**: bash, zsh, fish, PowerShell
- ✅ **Helpful caveats**: Instructions for shell wrapper setup
- ✅ **Auto-updates**: Via goreleaser integration

## Installation Instructions for Users

Once published, users can install with:

```bash
# Add the tap
brew tap gitwohq/tap

# Install gitwo
brew install gitwo

# Enable auto-cd (optional)
gitwo shell-install
```

## Troubleshooting

### Common Issues

1. **Formula not found**: Ensure the tap is added correctly
2. **SHA256 mismatch**: Rebuild and update the formula
3. **Permission denied**: Check GitHub token permissions

### Updating the Formula

The formula is automatically updated by goreleaser, but you can manually update:

```bash
# Update the formula
brew update
brew upgrade gitwo
```

## Next Steps

1. Set up the `gitwohq/homebrew-tap` repository
2. Configure GitHub Actions
3. Create the first release with `v0.1.0` tag
4. Test the installation process
5. Update documentation with installation instructions
