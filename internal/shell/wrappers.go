package shell

// GenerateBashZshWrapper generates a bash/zsh wrapper function
func GenerateBashZshWrapper() string {
	return `# gitwo wrapper: auto-cd after "gitwo new ..."
gitwo() {
    if [[ "$1" == "new" ]]; then
        export GITWO_WRAPPER=1
        local line dest
        line=$(command gitwo "$@" --shell | tail -n1) || return $?
        dest=${line#cd }
        [[ -n "$dest" ]] && builtin cd "$dest"
    else
        export GITWO_WRAPPER=1
        command gitwo "$@"
    fi
}

# gitwo-dev wrapper: auto-cd after "gitwo-dev new ..."
# Unalias gitwo-dev if it exists, then define as function
unalias gitwo-dev 2>/dev/null || true
gitwo-dev() {
    if [[ "$1" == "new" ]]; then
        export GITWO_WRAPPER=1
        local line dest
        line=$(command gitwo-dev "$@" --shell | tail -n1) || return $?
        dest=${line#cd }
        [[ -n "$dest" ]] && builtin cd "$dest"
    else
        export GITWO_WRAPPER=1
        command gitwo-dev "$@"
    fi
}
`
}

// GenerateFishWrapper generates a fish shell wrapper function
func GenerateFishWrapper() string {
	return `# gitwo wrapper: auto-cd after "gitwo new ..."
function gitwo --wraps gitwo --description "gitwo with auto-cd for new"
    set -x GITWO_WRAPPER 1
    if test (count $argv) -ge 1; and test $argv[1] = new
        set line (command gitwo $argv --shell | tail -n 1)
        set dest (string replace -r '^cd ' '' -- $line)
        if test -n "$dest"
            cd "$dest"
        end
    else
        command gitwo $argv
    end
end
`
}

// GeneratePowerShellWrapper generates a PowerShell wrapper function
func GeneratePowerShellWrapper() string {
	return `# gitwo wrapper: auto-cd after "gitwo new ..."
function gitwo {
    param([Parameter(ValueFromRemainingArguments=$true)] $Args)
    $env:GITWO_WRAPPER = "1"
    if ($Args.Count -gt 0 -and $Args[0] -eq 'new') {
        $line = & gitwo @Args --shell | Select-Object -Last 1
        $dest = $line -replace '^cd\s+',''
        if ($dest) {
            Set-Location -Path $dest
        }
    } else {
        & gitwo @Args
    }
}
`
}

// GenerateWrapper generates a wrapper for the specified shell
func GenerateWrapper(shell string) string {
	switch shell {
	case "bash", "zsh":
		return GenerateBashZshWrapper()
	case "fish":
		return GenerateFishWrapper()
	case "powershell":
		return GeneratePowerShellWrapper()
	default:
		return GenerateBashZshWrapper()
	}
}
