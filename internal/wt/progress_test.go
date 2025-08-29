package wt

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgressDisplay(t *testing.T) {
	t.Run("NewProgressDisplay", func(t *testing.T) {
		pd := NewProgressDisplay()
		assert.NotNil(t, pd)
		assert.Empty(t, pd.steps, "New progress display should have no steps")
	})

	t.Run("AddStep", func(t *testing.T) {
		pd := NewProgressDisplay()
		
		step := pd.AddStep("Test step")
		assert.NotNil(t, step)
		assert.Equal(t, "Test step", step.Message)
		assert.Equal(t, "PENDING", step.Status)
		assert.Len(t, pd.steps, 1)
		assert.Equal(t, step, pd.steps[0])
	})

	t.Run("UpdateStep", func(t *testing.T) {
		pd := NewProgressDisplay()
		step := pd.AddStep("Test step")
		
		pd.UpdateStep(step, "OK", "success")
		assert.Equal(t, "OK", step.Status)
		assert.Equal(t, "success", step.Details)
	})

	t.Run("RenderStep", func(t *testing.T) {
		pd := NewProgressDisplay()
		step := pd.AddStep("Test step")
		pd.UpdateStep(step, "OK", "success")
		
		// Capture output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		pd.RenderStep(step)
		
		w.Close()
		os.Stdout = originalStdout
		
		buf.ReadFrom(r)
		output := buf.String()
		
		assert.Contains(t, output, "• Test step✅ (success)")
	})

	t.Run("Render", func(t *testing.T) {
		pd := NewProgressDisplay()
		step1 := pd.AddStep("Step 1")
		step2 := pd.AddStep("Step 2")
		
		pd.UpdateStep(step1, "OK", "success")
		pd.UpdateStep(step2, "ERROR", "failed")
		
		// Capture output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		pd.Render()
		
		w.Close()
		os.Stdout = originalStdout
		
		buf.ReadFrom(r)
		output := buf.String()
		
		assert.Contains(t, output, "• Step 1✅ (success)")
		assert.Contains(t, output, "• Step 2❌ (failed)")
	})
}

func TestPrintHeader(t *testing.T) {
	t.Run("PrintHeader_WithValidInputs", func(t *testing.T) {
		// Capture output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		PrintHeader("/path/to/repo", "main", "abc123def456", "HEAD", "feature/test", "../test")
		
		w.Close()
		os.Stdout = originalStdout
		
		buf.ReadFrom(r)
		output := buf.String()
		
		assert.Contains(t, output, "gitwo v0.1 • new")
		assert.Contains(t, output, "Repo        : /path/to/repo  (main@abc123de)")
		assert.Contains(t, output, "Start point : HEAD")
		assert.Contains(t, output, "Branch      : feature/test   (normalized from \"feature/test\")")
		assert.Contains(t, output, "Worktree    : ../test")
	})

	t.Run("PrintHeader_WithEmptyRepoHead", func(t *testing.T) {
		// Capture output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		PrintHeader("/path/to/repo", "main", "", "HEAD", "feature/test", "../test")
		
		w.Close()
		os.Stdout = originalStdout
		
		buf.ReadFrom(r)
		output := buf.String()
		
		assert.Contains(t, output, "Repo        : /path/to/repo  (main@unknown)")
	})

	t.Run("PrintHeader_WithShortRepoHead", func(t *testing.T) {
		// Capture output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		PrintHeader("/path/to/repo", "main", "abc", "HEAD", "feature/test", "../test")
		
		w.Close()
		os.Stdout = originalStdout
		
		buf.ReadFrom(r)
		output := buf.String()
		
		assert.Contains(t, output, "Repo        : /path/to/repo  (main@abc)")
	})
}

func TestPrintFooter(t *testing.T) {
	t.Run("PrintFooter", func(t *testing.T) {
		// Capture output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		PrintFooter("../test", "feature/test")
		
		w.Close()
		os.Stdout = originalStdout
		
		buf.ReadFrom(r)
		output := buf.String()
		
		assert.Contains(t, output, "DONE  Worktree ready 🚀")
		assert.Contains(t, output, "  cd ../test")
		assert.Contains(t, output, "  git status")
		assert.Contains(t, output, "  git push -u origin feature/test")
	})
}

func TestPrintFooterWithSwitch(t *testing.T) {
	t.Run("PrintFooterWithSwitch", func(t *testing.T) {
		// Capture output
		var buf bytes.Buffer
		originalStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		
		PrintFooterWithSwitch("../test", "feature/test")
		
		w.Close()
		os.Stdout = originalStdout
		
		buf.ReadFrom(r)
		output := buf.String()
		
		assert.Contains(t, output, "DONE  Worktree ready 🚀")
		assert.Contains(t, output, "  cd ../test  # Run this to switch to the new worktree")
		assert.Contains(t, output, "⚡ For instant switching, add this to your shell:")
		assert.Contains(t, output, "   gitwo-switch() { cd ../test; }")
		assert.Contains(t, output, "🔧 Pro tip: Use 'eval \"$(gitwo new test --shell | tail -1)\"' to auto-switch!")
	})
}
