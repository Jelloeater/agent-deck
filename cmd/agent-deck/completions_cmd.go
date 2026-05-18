package main

import (
	"fmt"
	"os"
	"strings"
)

// handleCompletions generates shell completion scripts for Bash, ZSH, and Fish
func handleCompletions(args []string) {
	if len(args) == 0 {
		printCompletionsHelp()
		return
	}

	shell := args[0]
	switch strings.ToLower(shell) {
	case "bash":
		generateBashCompletion()
	case "zsh":
		generateZshCompletion()
	case "fish":
		generateFishCompletion()
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported shell '%s'\n", shell)
		fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, fish")
		os.Exit(1)
	}
}

func printCompletionsHelp() {
	fmt.Println("Usage: agent-deck completions <shell>")
	fmt.Println()
	fmt.Println("Generate shell completion scripts for agent-deck")
	fmt.Println()
	fmt.Println("Supported shells:")
	fmt.Println("  bash    Generate Bash completion script")
	fmt.Println("  zsh     Generate ZSH completion script")
	fmt.Println("  fish    Generate Fish completion script")
	fmt.Println()
	fmt.Println("Installation:")
	fmt.Println()
	fmt.Println("Bash:")
	fmt.Println("  agent-deck completions bash > /etc/bash_completion.d/agent-deck")
	fmt.Println("  # or for user-level:")
	fmt.Println("  agent-deck completions bash > ~/.bash_completion.d/agent-deck")
	fmt.Println()
	fmt.Println("ZSH:")
	fmt.Println("  agent-deck completions zsh > ~/.zsh/completions/_agent-deck")
	fmt.Println("  # Add to .zshrc if not already present:")
	fmt.Println("  # fpath=(~/.zsh/completions $fpath)")
	fmt.Println("  # autoload -U compinit && compinit")
	fmt.Println()
	fmt.Println("Fish:")
	fmt.Println("  agent-deck completions fish > ~/.config/fish/completions/agent-deck.fish")
}

func generateBashCompletion() {
	// Bash completion script using complete builtin
	script := `# bash completion for agent-deck
_agent_deck_completions() {
    local cur prev words cword
    _init_completion || return

    # Top-level commands
    local commands="add launch try list ls remove rm rename mv status session mcp skill plugin codex-hooks gemini-hooks group worktree wt web remote conductor profile update costs inbox feedback watcher openclaw oc notify-daemon hook-handler codex-notify hooks uninstall debug-dump completions version help"

    # Subcommands for each command
    local session_cmds="start stop restart fork attach show send stream search color move info set-telegram-warnings"
    local mcp_cmds="list attached attach detach"
    local skill_cmds="list attached attach detach source"
    local plugin_cmds="list install uninstall update info"
    local group_cmds="list create delete move rename set-parent"
    local worktree_cmds="list info cleanup"
    local remote_cmds="add remove rm list ls sessions attach rename update"
    local conductor_cmds="setup teardown status list"
    local profile_cmds="list create delete default migrate"
    local codex_hooks_cmds="install uninstall status"
    local gemini_hooks_cmds="install uninstall status"
    local skill_source_cmds="list add remove update"

    # If we're at the first argument, complete top-level commands
    if [[ $cword -eq 1 ]]; then
        COMPREPLY=($(compgen -W "$commands" -- "$cur"))
        return
    fi

    # Get the main command
    local main_cmd="${words[1]}"

    # Complete subcommands based on main command
    case "$main_cmd" in
        session)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$session_cmds" -- "$cur"))
            fi
            ;;
        mcp)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$mcp_cmds" -- "$cur"))
            fi
            ;;
        skill)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$skill_cmds" -- "$cur"))
            elif [[ $cword -eq 3 && "${words[2]}" == "source" ]]; then
                COMPREPLY=($(compgen -W "$skill_source_cmds" -- "$cur"))
            fi
            ;;
        plugin)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$plugin_cmds" -- "$cur"))
            fi
            ;;
        group)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$group_cmds" -- "$cur"))
            fi
            ;;
        worktree|wt)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$worktree_cmds" -- "$cur"))
            fi
            ;;
        remote)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$remote_cmds" -- "$cur"))
            fi
            ;;
        conductor)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$conductor_cmds" -- "$cur"))
            fi
            ;;
        profile)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$profile_cmds" -- "$cur"))
            fi
            ;;
        codex-hooks)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$codex_hooks_cmds" -- "$cur"))
            fi
            ;;
        gemini-hooks)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "$gemini_hooks_cmds" -- "$cur"))
            fi
            ;;
        completions)
            if [[ $cword -eq 2 ]]; then
                COMPREPLY=($(compgen -W "bash zsh fish" -- "$cur"))
            fi
            ;;
        add|launch|try)
            # Complete directories and files
            COMPREPLY=($(compgen -d -- "$cur"))
            ;;
        *)
            # Default to file/directory completion
            COMPREPLY=($(compgen -f -- "$cur"))
            ;;
    esac
}

complete -F _agent_deck_completions agent-deck
`
	fmt.Print(script)
}

func generateZshCompletion() {
	// ZSH completion script using _arguments
	script := `#compdef agent-deck

_agent_deck() {
    local line state

    _arguments -C \
        '(-p --profile)'{-p,--profile}'[Use specific profile]:profile:' \
        '(-g --group)'{-g,--group}'[Launch TUI scoped to group]:group:' \
        '--select[Launch TUI with cursor on session]:session:' \
        '1: :_agent_deck_commands' \
        '*::arg:->args'

    case $line[1] in
        session)
            _agent_deck_session
            ;;
        mcp)
            _agent_deck_mcp
            ;;
        skill)
            _agent_deck_skill
            ;;
        plugin)
            _agent_deck_plugin
            ;;
        group)
            _agent_deck_group
            ;;
        worktree|wt)
            _agent_deck_worktree
            ;;
        remote)
            _agent_deck_remote
            ;;
        conductor)
            _agent_deck_conductor
            ;;
        profile)
            _agent_deck_profile
            ;;
        codex-hooks)
            _agent_deck_codex_hooks
            ;;
        gemini-hooks)
            _agent_deck_gemini_hooks
            ;;
        completions)
            _agent_deck_completions
            ;;
        add|launch|try)
            _files -/
            ;;
    esac
}

_agent_deck_commands() {
    local commands=(
        'add:Add a new session'
        'launch:Add, start, and optionally send a message in one step'
        'try:Quick experiment (create/find dated folder + session)'
        'list:List all sessions'
        'ls:List all sessions (alias)'
        'remove:Remove a session'
        'rm:Remove a session (alias)'
        'rename:Rename a session'
        'mv:Rename a session (alias)'
        'status:Show session status summary'
        'session:Manage session lifecycle'
        'mcp:Manage MCP servers'
        'skill:Manage project skills'
        'plugin:Manage plugins'
        'codex-hooks:Manage Codex notify hook integration'
        'gemini-hooks:Manage Gemini hook integration'
        'group:Manage groups'
        'worktree:Manage git worktrees'
        'wt:Manage git worktrees (alias)'
        'web:Start TUI with web UI server'
        'remote:Manage remote agent-deck instances'
        'conductor:Manage conductor meta-agent orchestration'
        'profile:Manage profiles'
        'update:Check for and install updates'
        'costs:Show cost analysis'
        'inbox:Manage inbox'
        'feedback:Provide feedback'
        'watcher:File watcher management'
        'openclaw:OpenClaw integration'
        'oc:OpenClaw integration (alias)'
        'notify-daemon:Notification daemon'
        'hook-handler:Hook handler'
        'codex-notify:Codex notification'
        'hooks:Manage hooks'
        'uninstall:Uninstall Agent Deck'
        'debug-dump:Dump debug ring buffer'
        'completions:Generate shell completion scripts'
        'version:Show version'
        'help:Show help'
    )
    _describe 'command' commands
}

_agent_deck_session() {
    local commands=(
        'start:Start a session'\''s tmux process'
        'stop:Stop session process'
        'restart:Restart session (reload MCPs)'
        'fork:Fork Claude session with context'
        'attach:Attach to session interactively'
        'show:Show session details'
        'send:Send message to session'
        'stream:Stream session output'
        'search:Search sessions'
        'color:Set session color'
        'move:Move session to group'
        'info:Show session info'
        'set-telegram-warnings:Configure Telegram warnings'
    )
    _describe 'session command' commands
}

_agent_deck_mcp() {
    local commands=(
        'list:List available MCPs from config.toml'
        'attached:Show MCPs attached to a session'
        'attach:Attach MCP to session'
        'detach:Detach MCP from session'
    )
    _describe 'mcp command' commands
}

_agent_deck_skill() {
    local commands=(
        'list:List discoverable skills'
        'attached:Show skills attached to a session'
        'attach:Attach skill to session project'
        'detach:Detach skill from session project'
        'source:Manage global skill sources'
    )
    _describe 'skill command' commands
}

_agent_deck_plugin() {
    local commands=(
        'list:List available plugins'
        'install:Install a plugin'
        'uninstall:Uninstall a plugin'
        'update:Update a plugin'
        'info:Show plugin information'
    )
    _describe 'plugin command' commands
}

_agent_deck_group() {
    local commands=(
        'list:List all groups'
        'create:Create a new group'
        'delete:Delete a group'
        'move:Move session to group'
        'rename:Rename a group'
        'set-parent:Set parent group'
    )
    _describe 'group command' commands
}

_agent_deck_worktree() {
    local commands=(
        'list:List worktrees with session associations'
        'info:Show worktree info for a session'
        'cleanup:Find and remove orphaned worktrees/sessions'
    )
    _describe 'worktree command' commands
}

_agent_deck_remote() {
    local commands=(
        'add:Register a remote agent-deck instance'
        'remove:Remove a remote'
        'rm:Remove a remote (alias)'
        'list:List configured remotes'
        'ls:List configured remotes (alias)'
        'sessions:Show sessions on remote(s)'
        'attach:Attach to a remote session'
        'rename:Rename a remote session'
        'update:Install/upgrade agent-deck on remote(s)'
    )
    _describe 'remote command' commands
}

_agent_deck_conductor() {
    local commands=(
        'setup:Set up conductor (Telegram bridge + sessions)'
        'teardown:Stop conductor and remove bridge daemon'
        'status:Show conductor health across profiles'
        'list:List configured conductors'
    )
    _describe 'conductor command' commands
}

_agent_deck_profile() {
    local commands=(
        'list:List all profiles'
        'create:Create a new profile'
        'delete:Delete a profile'
        'default:Show or set default profile'
        'migrate:Migrate profile data'
    )
    _describe 'profile command' commands
}

_agent_deck_codex_hooks() {
    local commands=(
        'install:Install or upgrade Codex notify hook'
        'uninstall:Remove Codex notify hook'
        'status:Show Codex hook install status'
    )
    _describe 'codex-hooks command' commands
}

_agent_deck_gemini_hooks() {
    local commands=(
        'install:Install Gemini hooks'
        'uninstall:Remove Gemini hooks'
        'status:Show Gemini hooks install status'
    )
    _describe 'gemini-hooks command' commands
}

_agent_deck_completions() {
    local shells=(
        'bash:Generate Bash completion script'
        'zsh:Generate ZSH completion script'
        'fish:Generate Fish completion script'
    )
    _describe 'shell' shells
}

_agent_deck "$@"
`
	fmt.Print(script)
}

func generateFishCompletion() {
	// Fish completion script
	script := `# fish completion for agent-deck

# Remove any existing completions
complete -c agent-deck -e

# Global options
complete -c agent-deck -s p -l profile -d 'Use specific profile' -r
complete -c agent-deck -s g -l group -d 'Launch TUI scoped to group' -r
complete -c agent-deck -l select -d 'Launch TUI with cursor on session' -r

# Main commands
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'add' -d 'Add a new session'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'launch' -d 'Add, start, and optionally send a message'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'try' -d 'Quick experiment'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'list' -d 'List all sessions'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'ls' -d 'List all sessions'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'remove' -d 'Remove a session'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'rm' -d 'Remove a session'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'rename' -d 'Rename a session'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'mv' -d 'Rename a session'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'status' -d 'Show session status summary'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'session' -d 'Manage session lifecycle'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'mcp' -d 'Manage MCP servers'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'skill' -d 'Manage project skills'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'plugin' -d 'Manage plugins'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'codex-hooks' -d 'Manage Codex notify hook integration'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'gemini-hooks' -d 'Manage Gemini hook integration'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'group' -d 'Manage groups'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'worktree' -d 'Manage git worktrees'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'wt' -d 'Manage git worktrees'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'web' -d 'Start TUI with web UI server'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'remote' -d 'Manage remote instances'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'conductor' -d 'Manage conductor orchestration'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'profile' -d 'Manage profiles'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'update' -d 'Check for and install updates'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'costs' -d 'Show cost analysis'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'inbox' -d 'Manage inbox'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'feedback' -d 'Provide feedback'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'watcher' -d 'File watcher management'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'openclaw' -d 'OpenClaw integration'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'oc' -d 'OpenClaw integration'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'completions' -d 'Generate shell completion scripts'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'version' -d 'Show version'
complete -c agent-deck -f -n '__fish_use_subcommand' -a 'help' -d 'Show help'

# Session subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'start' -d 'Start a session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'stop' -d 'Stop session process'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'restart' -d 'Restart session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'fork' -d 'Fork Claude session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'attach' -d 'Attach to session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'show' -d 'Show session details'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'send' -d 'Send message to session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'stream' -d 'Stream session output'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'search' -d 'Search sessions'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'color' -d 'Set session color'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'move' -d 'Move session to group'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'info' -d 'Show session info'
complete -c agent-deck -f -n '__fish_seen_subcommand_from session' -a 'set-telegram-warnings' -d 'Configure Telegram warnings'

# MCP subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from mcp' -a 'list' -d 'List available MCPs'
complete -c agent-deck -f -n '__fish_seen_subcommand_from mcp' -a 'attached' -d 'Show MCPs attached to session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from mcp' -a 'attach' -d 'Attach MCP to session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from mcp' -a 'detach' -d 'Detach MCP from session'

# Skill subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from skill' -a 'list' -d 'List discoverable skills'
complete -c agent-deck -f -n '__fish_seen_subcommand_from skill' -a 'attached' -d 'Show skills attached to session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from skill' -a 'attach' -d 'Attach skill to session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from skill' -a 'detach' -d 'Detach skill from session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from skill' -a 'source' -d 'Manage global skill sources'

# Plugin subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from plugin' -a 'list' -d 'List available plugins'
complete -c agent-deck -f -n '__fish_seen_subcommand_from plugin' -a 'install' -d 'Install a plugin'
complete -c agent-deck -f -n '__fish_seen_subcommand_from plugin' -a 'uninstall' -d 'Uninstall a plugin'
complete -c agent-deck -f -n '__fish_seen_subcommand_from plugin' -a 'update' -d 'Update a plugin'
complete -c agent-deck -f -n '__fish_seen_subcommand_from plugin' -a 'info' -d 'Show plugin information'

# Group subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from group' -a 'list' -d 'List all groups'
complete -c agent-deck -f -n '__fish_seen_subcommand_from group' -a 'create' -d 'Create a new group'
complete -c agent-deck -f -n '__fish_seen_subcommand_from group' -a 'delete' -d 'Delete a group'
complete -c agent-deck -f -n '__fish_seen_subcommand_from group' -a 'move' -d 'Move session to group'
complete -c agent-deck -f -n '__fish_seen_subcommand_from group' -a 'rename' -d 'Rename a group'
complete -c agent-deck -f -n '__fish_seen_subcommand_from group' -a 'set-parent' -d 'Set parent group'

# Worktree subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from worktree wt' -a 'list' -d 'List worktrees'
complete -c agent-deck -f -n '__fish_seen_subcommand_from worktree wt' -a 'info' -d 'Show worktree info'
complete -c agent-deck -f -n '__fish_seen_subcommand_from worktree wt' -a 'cleanup' -d 'Remove orphaned worktrees'

# Remote subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'add' -d 'Register a remote instance'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'remove' -d 'Remove a remote'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'rm' -d 'Remove a remote'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'list' -d 'List configured remotes'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'ls' -d 'List configured remotes'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'sessions' -d 'Show sessions on remotes'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'attach' -d 'Attach to remote session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'rename' -d 'Rename remote session'
complete -c agent-deck -f -n '__fish_seen_subcommand_from remote' -a 'update' -d 'Install/upgrade on remotes'

# Conductor subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from conductor' -a 'setup' -d 'Set up conductor'
complete -c agent-deck -f -n '__fish_seen_subcommand_from conductor' -a 'teardown' -d 'Stop conductor'
complete -c agent-deck -f -n '__fish_seen_subcommand_from conductor' -a 'status' -d 'Show conductor health'
complete -c agent-deck -f -n '__fish_seen_subcommand_from conductor' -a 'list' -d 'List configured conductors'

# Profile subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from profile' -a 'list' -d 'List all profiles'
complete -c agent-deck -f -n '__fish_seen_subcommand_from profile' -a 'create' -d 'Create a new profile'
complete -c agent-deck -f -n '__fish_seen_subcommand_from profile' -a 'delete' -d 'Delete a profile'
complete -c agent-deck -f -n '__fish_seen_subcommand_from profile' -a 'default' -d 'Show or set default profile'
complete -c agent-deck -f -n '__fish_seen_subcommand_from profile' -a 'migrate' -d 'Migrate profile data'

# Codex hooks subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from codex-hooks' -a 'install' -d 'Install or upgrade Codex hook'
complete -c agent-deck -f -n '__fish_seen_subcommand_from codex-hooks' -a 'uninstall' -d 'Remove Codex hook'
complete -c agent-deck -f -n '__fish_seen_subcommand_from codex-hooks' -a 'status' -d 'Show Codex hook status'

# Gemini hooks subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from gemini-hooks' -a 'install' -d 'Install Gemini hooks'
complete -c agent-deck -f -n '__fish_seen_subcommand_from gemini-hooks' -a 'uninstall' -d 'Remove Gemini hooks'
complete -c agent-deck -f -n '__fish_seen_subcommand_from gemini-hooks' -a 'status' -d 'Show Gemini hooks status'

# Completions subcommands
complete -c agent-deck -f -n '__fish_seen_subcommand_from completions' -a 'bash' -d 'Generate Bash completion script'
complete -c agent-deck -f -n '__fish_seen_subcommand_from completions' -a 'zsh' -d 'Generate ZSH completion script'
complete -c agent-deck -f -n '__fish_seen_subcommand_from completions' -a 'fish' -d 'Generate Fish completion script'
`
	fmt.Print(script)
}
