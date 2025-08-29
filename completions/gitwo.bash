# Bash completion for gitwo
_gitwo() {
    local cur prev opts cmds
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    cmds="new list remove config init shell-init shell-install code completion"
    
    case "${prev}" in
        gitwo)
            COMPREPLY=( $(compgen -W "${cmds}" -- "${cur}") )
            return 0
            ;;
        new|remove)
            # Complete with existing worktree names
            local worktrees=$(git worktree list --porcelain 2>/dev/null | grep -E '^worktree' | cut -d' ' -f2 | xargs -n1 basename 2>/dev/null)
            if [ -n "$worktrees" ]; then
                COMPREPLY=( $(compgen -W "${worktrees}" -- "${cur}") )
            fi
            return 0
            ;;
        shell-init|shell-install)
            local shells="bash zsh fish powershell"
            COMPREPLY=( $(compgen -W "${shells}" -- "${cur}") )
            return 0
            ;;
        completion)
            local shells="bash zsh fish powershell"
            COMPREPLY=( $(compgen -W "${shells}" -- "${cur}") )
            return 0
            ;;
    esac
}

complete -F _gitwo gitwo
