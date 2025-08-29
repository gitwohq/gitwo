package wt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AddResult contains information about the created worktree
type AddResult struct {
	Path   string
	Branch string
	Head   string
}

func Add(path, branch, startPoint string) (*AddResult, error) {
	return AddWithConfig(path, branch, startPoint, nil)
}

func AddWithConfig(path, branch, startPoint string, config *Config) (*AddResult, error) {
	return AddWithConfigSilent(path, branch, startPoint, config, false)
}

func AddWithConfigSilent(path, branch, startPoint string, config *Config, silent bool) (*AddResult, error) {
	// Create progress display
	var progress *ProgressDisplay
	if silent {
		progress = NewSilentProgressDisplay()
	} else {
		progress = NewProgressDisplay()
	}

	// Step 1: Check repository
	repoStep := progress.AddStep("Checking repository")
	repoRootPath, err := repoRoot()
	if err != nil {
		if !silent {
			progress.UpdateStep(repoStep, "ERROR", err.Error())
			progress.RenderStep(repoStep)
		}
		return nil, err
	}
	if !silent {
		progress.UpdateStep(repoStep, "OK", "")
		progress.RenderStep(repoStep)
	}

	// Get repository info
	repoBranch := ""
	if out, err := gitOut("branch", "--show-current"); err == nil {
		repoBranch = string(bytesTrimNL(out))
	}

	repoHead := ""
	if out, err := gitOut("rev-parse", "HEAD"); err == nil {
		repoHead = string(bytesTrimNL(out))
	}

	// Print header
	if !silent {
		PrintHeader(repoRootPath, repoBranch, repoHead, startPoint, branch, path)

		// Show branch normalization if applicable
		if strings.HasPrefix(branch, "feature/") {
			originalName := strings.TrimPrefix(branch, "feature/")
			fmt.Printf("note: normalized \"%s\" → \"%s\" (policy: feature/*)\n\n", originalName, branch)
		}
	}
	// Validate inputs
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if branch == "" {
		return nil, fmt.Errorf("branch cannot be empty")
	}
	if startPoint == "" {
		return nil, fmt.Errorf("start point cannot be empty")
	}
	// Step 2: Resolve base and fetch
	if !silent {
		fetchStep := progress.AddStep("Resolving base and fetching")
		if err := gitSilent("fetch", "--all", "--prune"); err != nil {
			// Ignore fetch errors - they don't affect worktree creation
			// This happens when no remote is configured, remote is unreachable, or repo doesn't exist yet
			progress.UpdateStep(fetchStep, "OK", "no remote configured")
		} else {
			// Check if start point is behind remote
			if refIsRemote(startPoint) {
				remoteRef := startPoint
				if strings.HasPrefix(startPoint, "origin/") {
					remoteRef = strings.TrimPrefix(startPoint, "origin/")
				}

				// Get local and remote commit hashes
				localHash, _ := gitOut("rev-parse", remoteRef)
				remoteHash, _ := gitOut("rev-parse", startPoint)

				if len(localHash) > 0 && len(remoteHash) > 0 && !bytes.Equal(localHash, remoteHash) {
					// Count commits ahead/behind
					out, _ := gitOut("rev-list", "--count", fmt.Sprintf("%s..%s", string(bytesTrimNL(localHash)), string(bytesTrimNL(remoteHash))))
					if len(out) > 0 {
						count := strings.TrimSpace(string(out))
						if count != "0" {
							fmt.Printf("WARN  start point %s is %s commits ahead of local %s\n", startPoint, count, remoteRef)
							fmt.Printf("      consider: git fetch --prune\n\n")
						}
					}
				}
			}

			progress.UpdateStep(fetchStep, "OK", fmt.Sprintf("%s up to date", startPoint))
		}
		progress.RenderStep(fetchStep)
	} else {
		// Silent mode: just do the fetch without output
		gitSilent("fetch", "--all", "--prune")
	}

	// Step 3: Create branch
	branchStep := progress.AddStep("Creating branch")

	// Check if target path already exists
	if _, err := os.Stat(path); err == nil {
		progress.UpdateStep(branchStep, "ERROR", fmt.Sprintf("target path exists: %s", path))
		progress.RenderStep(branchStep)
		return nil, fmt.Errorf("target path exists: %s\n      use --force to reuse or --path PATH to override", path)
	}

	// Ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		progress.UpdateStep(branchStep, "ERROR", fmt.Sprintf("mkdir: %v", err))
		progress.RenderStep(branchStep)
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	progress.UpdateStep(branchStep, "OK", branch)
	progress.RenderStep(branchStep)

	// Step 4: Create worktree
	worktreeStep := progress.AddStep("Creating worktree")
	if err := git("worktree", "add", "-B", branch, path, startPoint); err != nil {
		progress.UpdateStep(worktreeStep, "ERROR", err.Error())
		progress.RenderStep(worktreeStep)
		return nil, err
	}
	progress.UpdateStep(worktreeStep, "OK", path)
	progress.RenderStep(worktreeStep)

	// Step 5: Set upstream
	upstreamStep := progress.AddStep("Setting upstream")
	if refIsRemote(startPoint) {
		if err := run("bash", "-lc", fmt.Sprintf("cd %q && git branch --set-upstream-to %s %s", path, startPoint, branch)); err != nil {
			progress.UpdateStep(upstreamStep, "ERROR", err.Error())
		} else {
			progress.UpdateStep(upstreamStep, "OK", fmt.Sprintf("origin/%s", branch))
		}
	} else {
		progress.UpdateStep(upstreamStep, "OK", "no upstream set")
	}
	progress.RenderStep(upstreamStep)

	// Get the actual branch name that was checked out
	actualBranch := branch
	if out, err := gitOut("-C", path, "branch", "--show-current"); err == nil {
		actualBranch = string(bytesTrimNL(out))
	}

	// Get the HEAD commit
	head := ""
	if out, err := gitOut("-C", path, "rev-parse", "HEAD"); err == nil {
		head = string(bytesTrimNL(out))
	}

	result := &AddResult{
		Path:   path,
		Branch: actualBranch,
		Head:   head,
	}

	// Auto-switch if configured
	if config != nil && config.AutoSwitch {
		switchStep := progress.AddStep("Preparing auto-switch")
		progress.UpdateStep(switchStep, "OK", "ready to switch")
		progress.RenderStep(switchStep)
	}

	return result, nil
}
