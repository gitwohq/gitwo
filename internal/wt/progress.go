package wt

import (
	"fmt"
	"time"
)

// ProgressStep represents a step in the worktree creation process
type ProgressStep struct {
	Message   string
	Status    string
	Details   string
	StartTime time.Time
}

// ProgressDisplay handles the display of progress steps
type ProgressDisplay struct {
	steps  []*ProgressStep
	silent bool
}

// NewProgressDisplay creates a new progress display
func NewProgressDisplay() *ProgressDisplay {
	return &ProgressDisplay{
		steps:  make([]*ProgressStep, 0),
		silent: false,
	}
}

// NewSilentProgressDisplay creates a silent progress display
func NewSilentProgressDisplay() *ProgressDisplay {
	return &ProgressDisplay{
		steps:  make([]*ProgressStep, 0),
		silent: true,
	}
}

// AddStep adds a new step to the progress display
func (pd *ProgressDisplay) AddStep(message string) *ProgressStep {
	step := &ProgressStep{
		Message:   message,
		Status:    "PENDING",
		StartTime: time.Now(),
	}
	pd.steps = append(pd.steps, step)
	return step
}

// UpdateStep updates a step with status and details
func (pd *ProgressDisplay) UpdateStep(step *ProgressStep, status, details string) {
	step.Status = status
	step.Details = details
}

// Render displays the current progress
func (pd *ProgressDisplay) Render() {
	if pd.silent {
		return
	}
	for _, step := range pd.steps {
		statusIcon := "⏳"
		if step.Status == "OK" {
			statusIcon = "✅"
		} else if step.Status == "ERROR" {
			statusIcon = "❌"
		}

		details := ""
		if step.Details != "" {
			details = fmt.Sprintf(" (%s)", step.Details)
		}

		fmt.Printf("• %s%s%s\n", step.Message, statusIcon, details)
	}
}

// RenderStep renders a single step immediately
func (pd *ProgressDisplay) RenderStep(step *ProgressStep) {
	if pd.silent {
		return
	}
	statusIcon := "⏳"
	if step.Status == "OK" {
		statusIcon = "✅"
	} else if step.Status == "ERROR" {
		statusIcon = "❌"
	}

	details := ""
	if step.Details != "" {
		details = fmt.Sprintf(" (%s)", step.Details)
	}

	fmt.Printf("• %s%s%s\n", step.Message, statusIcon, details)
}

// PrintHeader prints the header information
func PrintHeader(repoPath, repoBranch, repoHead, startPoint, branch, worktreePath string) {
	fmt.Printf("gitwo v0.1 • new\n")
	
	// Safely handle empty repoHead
	headDisplay := ""
	if len(repoHead) >= 8 {
		headDisplay = repoHead[:8]
	} else if len(repoHead) > 0 {
		headDisplay = repoHead
	} else {
		headDisplay = "unknown"
	}
	
	fmt.Printf("Repo        : %s  (%s@%s)\n", repoPath, repoBranch, headDisplay)
	fmt.Printf("Start point : %s\n", startPoint)
	fmt.Printf("Branch      : %s   (normalized from \"%s\")\n", branch, branch)
	fmt.Printf("Worktree    : %s\n\n", worktreePath)
}

// PrintFooter prints the footer with next steps
func PrintFooter(worktreePath, branch string) {
	fmt.Printf("\nDONE  Worktree ready 🚀\n")
	fmt.Printf("Next:\n")
	fmt.Printf("  cd %s\n", worktreePath)
	fmt.Printf("  git status\n")
	fmt.Printf("  git push -u origin %s\n", branch)
}

// PrintFooterWithSwitch prints the footer with auto-switch instructions
func PrintFooterWithSwitch(worktreePath, branch string) {
	fmt.Printf("\nDONE  Worktree ready 🚀\n")
	fmt.Printf("Next:\n")
	fmt.Printf("  cd %s  # Run this to switch to the new worktree\n", worktreePath)
	fmt.Printf("  git status\n")
	fmt.Printf("  git push -u origin %s\n", branch)
	fmt.Printf("\n⚡ For instant switching, add this to your shell:\n")
	fmt.Printf("   gitwo-switch() { cd %s; }\n", worktreePath)
	fmt.Printf("\n🔧 Pro tip: Use 'eval \"$(gitwo new test --shell | tail -1)\"' to auto-switch!\n")
}
