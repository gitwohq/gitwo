# Fish completion for gitwo
complete -c gitwo -f

# Main commands
complete -c gitwo -n __fish_use_subcommand -a new -d "Create new worktree"
complete -c gitwo -n __fish_use_subcommand -a list -d "List worktrees"
complete -c gitwo -n __fish_use_subcommand -a remove -d "Remove worktree"
complete -c gitwo -n __fish_use_subcommand -a config -d "Manage configuration"
complete -c gitwo -n __fish_use_subcommand -a init -d "Initialize gitwo"
complete -c gitwo -n __fish_use_subcommand -a shell-init -d "Print shell wrapper"
complete -c gitwo -n __fish_use_subcommand -a shell-install -d "Install shell wrapper"
complete -c gitwo -n __fish_use_subcommand -a code -d "Open worktree in editor"
complete -c gitwo -n __fish_use_subcommand -a completion -d "Generate shell completion"

# Worktree name completion for new/remove
complete -c gitwo -n "__fish_seen_subcommand_from new remove" -a "(git worktree list --porcelain 2>/dev/null | grep -E '^worktree' | cut -d' ' -f2 | xargs -n1 basename 2>/dev/null)"

# Shell completion for shell commands
complete -c gitwo -n "__fish_seen_subcommand_from shell-init shell-install completion" -a "bash zsh fish powershell"
