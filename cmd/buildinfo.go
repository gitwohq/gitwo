package cmd

import "fmt"

// Public vars are already referenced in your root.go.
// We set them here and also wire Cobra's built-in version output.
func SetBuildInfo(v, c, d string) {
	// Make these visible anywhere in package cmd
	Version, Commit, Date = v, c, d

	// Let `gitwo --version` / `gitwo version` show a nice string
	rootCmd.Version = fmt.Sprintf("%s (commit %s, built %s)", v, c, d)

	// Optional: control format (safe to re-set)
	rootCmd.SetVersionTemplate("{{.Use}} version {{.Version}}\n")
}
