# PowerShell completion for gitwo
Register-ArgumentCompleter -Native -CommandName gitwo -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    
    $completions = @(
        [System.Management.Automation.CompletionResult]::new('new', 'new', 'ParameterValue', 'Create new worktree'),
        [System.Management.Automation.CompletionResult]::new('list', 'list', 'ParameterValue', 'List worktrees'),
        [System.Management.Automation.CompletionResult]::new('remove', 'remove', 'ParameterValue', 'Remove worktree'),
        [System.Management.Automation.CompletionResult]::new('config', 'config', 'ParameterValue', 'Manage configuration'),
        [System.Management.Automation.CompletionResult]::new('init', 'init', 'ParameterValue', 'Initialize gitwo'),
        [System.Management.Automation.CompletionResult]::new('shell-init', 'shell-init', 'ParameterValue', 'Print shell wrapper'),
        [System.Management.Automation.CompletionResult]::new('shell-install', 'shell-install', 'ParameterValue', 'Install shell wrapper'),
        [System.Management.Automation.CompletionResult]::new('code', 'code', 'ParameterValue', 'Open worktree in editor'),
        [System.Management.Automation.CompletionResult]::new('completion', 'completion', 'ParameterValue', 'Generate shell completion')
    )
    
    $completions.Where{ $_.CompletionText -like "$wordToComplete*" }
}
