# Zsh completion for gitwo
_gitwo() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    _arguments -C \
        '1: :->cmds' \
        '*:: :->args'

    case $state in
        cmds)
            _values 'gitwo commands' \
                'new[Create new worktree]' \
                'list[List worktrees]' \
                'remove[Remove worktree]' \
                'config[Manage configuration]' \
                'init[Initialize gitwo]' \
                'shell-init[Print shell wrapper]' \
                'shell-install[Install shell wrapper]' \
                'code[Open worktree in editor]' \
                'completion[Generate shell completion]'
            ;;
        args)
            case $line[1] in
                new|remove)
                    # Complete with existing worktree names
                    local worktrees=$(git worktree list --porcelain 2>/dev/null | grep -E '^worktree' | cut -d' ' -f2 | xargs -n1 basename 2>/dev/null)
                    if [ -n "$worktrees" ]; then
                        _values 'worktrees' ${=worktrees}
                    fi
                    ;;
                shell-init|shell-install|completion)
                    _values 'shells' 'bash' 'zsh' 'fish' 'powershell'
                    ;;
            esac
            ;;
    esac
}

compdef _gitwo gitwo
