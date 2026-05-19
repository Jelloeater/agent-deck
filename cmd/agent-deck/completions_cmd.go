package main

import (
"fmt"
"os"
"strings"
)

// subCmdEntry describes a single subcommand used by completion generators.
type subCmdEntry struct {
name string
desc string
}

// nestedEntry describes a group of sub-subcommands under a parent subcommand.
type nestedEntry struct {
parentSub string
subs      []subCmdEntry
}

// cmdEntry is the single source of truth for one top-level CLI command and all
// its completable subcommands. All three shell generators read from this table;
// updating a subcommand list here fixes Bash, Zsh, and Fish simultaneously.
type cmdEntry struct {
name    string
aliases []string
desc    string
subs    []subCmdEntry
nested  []nestedEntry // ordered sub-subcommand groups (e.g. "source" under "skill")
}

// completionCommands mirrors the dispatch table in main.go. Keep it in sync
// with the case branches in handleSession, handleMCP, etc.
var completionCommands = []cmdEntry{
{name: "add", desc: "Add a new session"},
{name: "launch", desc: "Add, start, and optionally send a message in one step"},
{name: "try", desc: "Quick experiment (create/find dated folder + session)"},
{name: "list", aliases: []string{"ls"}, desc: "List all sessions"},
{name: "remove", aliases: []string{"rm"}, desc: "Remove a session"},
{name: "rename", aliases: []string{"mv"}, desc: "Rename a session"},
{name: "status", desc: "Show session status summary"},
{
name: "session",
desc: "Manage session lifecycle",
subs: []subCmdEntry{
{"start", "Start a session's tmux process"},
{"stop", "Stop session process"},
{"remove", "Remove a session"},
{"restart", "Restart session (reload MCPs)"},
{"revive", "Revive a stopped session"},
{"fork", "Fork session with context"},
{"attach", "Attach to session interactively"},
{"show", "Show session details"},
{"current", "Show current session"},
{"set-parent", "Set parent group"},
{"unset-parent", "Unset parent group"},
{"update", "Update session properties"},
{"set-transition-notify", "Configure transition notifications"},
{"set-title-lock", "Configure title lock"},
{"set", "Set session property"},
{"move", "Move session to group"},
{"mv", "Move session to group (alias)"},
{"send", "Send message to session"},
{"output", "Show session output"},
{"search", "Search sessions"},
},
},
{
name: "mcp",
desc: "Manage MCP servers",
subs: []subCmdEntry{
{"list", "List available MCPs from config.toml"},
{"ls", "List available MCPs (alias)"},
{"attached", "Show MCPs attached to a session"},
{"attach", "Attach MCP to session"},
{"detach", "Detach MCP from session"},
{"server", "Manage background MCP server process"},
},
nested: []nestedEntry{
{
parentSub: "server",
subs: []subCmdEntry{
{"start", "Start MCP server"},
{"stop", "Stop MCP server"},
{"status", "Show MCP server status"},
},
},
},
},
{name: "plugin", desc: "Manage plugins", subs: []subCmdEntry{
{"list", "List available plugins"},
{"ls", "List available plugins (alias)"},
{"attached", "Show plugins attached to session"},
{"attach", "Attach plugin to session"},
{"detach", "Detach plugin from session"},
}},
{
name: "skill",
desc: "Manage project skills",
subs: []subCmdEntry{
{"list", "List discoverable skills"},
{"ls", "List discoverable skills (alias)"},
{"attached", "Show skills attached to a session"},
{"attach", "Attach skill to session project"},
{"detach", "Detach skill from session project"},
{"source", "Manage global skill sources"},
},
nested: []nestedEntry{
{
parentSub: "source",
subs: []subCmdEntry{
{"list", "List skill sources"},
{"ls", "List skill sources (alias)"},
{"add", "Add a skill source"},
{"remove", "Remove a skill source"},
{"rm", "Remove a skill source (alias)"},
},
},
},
},
{name: "mcp-proxy", desc: "MCP proxy server"},
{
name: "group",
desc: "Manage groups",
subs: []subCmdEntry{
{"list", "List all groups"},
{"ls", "List all groups (alias)"},
{"create", "Create a new group"},
{"new", "Create a new group (alias)"},
{"update", "Update a group"},
{"set", "Update a group (alias)"},
{"delete", "Delete a group"},
{"rm", "Delete a group (alias)"},
{"remove", "Delete a group (alias)"},
{"move", "Move session to group"},
{"mv", "Move session to group (alias)"},
{"change", "Change session group"},
{"reparent", "Change session group (alias)"},
{"reorder", "Reorder groups"},
{"sort", "Reorder groups (alias)"},
},
},
{
name:    "worktree",
aliases: []string{"wt"},
desc:    "Manage git worktrees",
subs: []subCmdEntry{
{"list", "List worktrees with session associations"},
{"ls", "List worktrees (alias)"},
{"info", "Show worktree info for a session"},
{"cleanup", "Find and remove orphaned worktrees/sessions"},
{"finish", "Mark worktree task as finished"},
},
},
{name: "web", desc: "Start TUI with web UI server"},
{
name: "remote",
desc: "Manage remote agent-deck instances",
subs: []subCmdEntry{
{"add", "Register a remote agent-deck instance"},
{"remove", "Remove a remote"},
{"rm", "Remove a remote (alias)"},
{"list", "List configured remotes"},
{"ls", "List configured remotes (alias)"},
{"sessions", "Show sessions on remote(s)"},
{"attach", "Attach to a remote session"},
{"rename", "Rename a remote session"},
{"update", "Install/upgrade agent-deck on remote(s)"},
},
},
{
name: "conductor",
desc: "Manage conductor meta-agent orchestration",
subs: []subCmdEntry{
{"setup", "Set up conductor (Telegram bridge + sessions)"},
{"teardown", "Stop conductor and remove bridge daemon"},
{"status", "Show conductor health across profiles"},
{"list", "List configured conductors"},
{"move", "Move conductor to another profile"},
},
},
{name: "watcher", desc: "File watcher management"},
{name: "openclaw", aliases: []string{"oc"}, desc: "OpenClaw integration"},
{name: "costs", desc: "Show cost analysis"},
{name: "inbox", desc: "Manage inbox"},
{name: "feedback", desc: "Provide feedback"},
{name: "notify-daemon", desc: "Notification daemon"},
{name: "hook-handler", desc: "Hook handler"},
{name: "codex-notify", desc: "Codex notification"},
{name: "hooks", desc: "Manage hooks"},
{
name: "codex-hooks",
desc: "Manage Codex notify hook integration",
subs: []subCmdEntry{
{"install", "Install or upgrade Codex notify hook"},
{"uninstall", "Remove Codex notify hook"},
{"status", "Show Codex hook install status"},
},
},
{
name: "gemini-hooks",
desc: "Manage Gemini hook integration",
subs: []subCmdEntry{
{"install", "Install Gemini hooks"},
{"uninstall", "Remove Gemini hooks"},
{"status", "Show Gemini hooks install status"},
},
},
{name: "profile", desc: "Manage profiles", subs: []subCmdEntry{
{"list", "List all profiles"},
{"ls", "List all profiles (alias)"},
{"create", "Create a new profile"},
{"new", "Create a new profile (alias)"},
{"delete", "Delete a profile"},
{"rm", "Delete a profile (alias)"},
{"default", "Show or set default profile"},
}},
{name: "update", desc: "Check for and install updates"},
{name: "debug-dump", desc: "Dump debug ring buffer"},
{name: "uninstall", desc: "Uninstall Agent Deck"},
{
name: "completions",
desc: "Generate shell completion scripts",
subs: []subCmdEntry{
{"bash", "Generate Bash completion script"},
{"zsh", "Generate ZSH completion script"},
{"fish", "Generate Fish completion script"},
},
},
{name: "version", desc: "Show version"},
{name: "help", desc: "Show help"},
}

// handleCompletions dispatches completions subcommands.
func handleCompletions(args []string) {
if code := runCompletions(args); code != 0 {
os.Exit(code)
}
}

// runCompletions is the testable core of handleCompletions.
// It returns 0 on success and 1 on error, without calling os.Exit.
func runCompletions(args []string) int {
if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
printCompletionsHelp()
return 0
}
switch strings.ToLower(args[0]) {
case "bash":
generateBashCompletion()
case "zsh":
generateZshCompletion()
case "fish":
generateFishCompletion()
default:
fmt.Fprintf(os.Stderr, "Error: unsupported shell %q\n", args[0])
fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, fish")
return 1
}
return 0
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
fmt.Println("Bash (system-wide, requires sudo):")
fmt.Println("  agent-deck completions bash | sudo tee /etc/bash_completion.d/agent-deck > /dev/null")
fmt.Println()
fmt.Println("Bash (user-level, add to ~/.bashrc):")
fmt.Println("  mkdir -p ~/.bash_completion.d")
fmt.Println("  agent-deck completions bash > ~/.bash_completion.d/agent-deck")
fmt.Println("  # Then add to ~/.bashrc: source ~/.bash_completion.d/agent-deck")
fmt.Println()
fmt.Println("ZSH (add directory to fpath BEFORE compinit):")
fmt.Println("  mkdir -p ~/.zsh/completions")
fmt.Println("  agent-deck completions zsh > ~/.zsh/completions/_agent-deck")
fmt.Println("  # Add to ~/.zshrc (before compinit):")
fmt.Println("  #   fpath=(~/.zsh/completions $fpath)")
fmt.Println("  #   autoload -U compinit && compinit")
fmt.Println()
fmt.Println("Fish (auto-loaded from completions directory):")
fmt.Println("  mkdir -p ~/.config/fish/completions")
fmt.Println("  agent-deck completions fish > ~/.config/fish/completions/agent-deck.fish")
}

// generateBashCompletion emits a Bash completion script derived from completionCommands.
// The script uses a fallback for systems without bash-completion installed, and skips
// global flags (-p/--profile, -g/--group, --select, --color) when determining which
// main command is active.
func generateBashCompletion() {
var sb strings.Builder

sb.WriteString("# bash completion for agent-deck\n")
sb.WriteString("# Generated by: agent-deck completions bash\n\n")

sb.WriteString("_agent_deck_completions() {\n")
sb.WriteString("    local cur prev words cword\n")
// Fallback: _init_completion is supplied by bash-completion (optional package).
// Use COMP_WORDS/COMP_CWORD directly when it is not available.
sb.WriteString("    if declare -F _init_completion > /dev/null 2>&1; then\n")
sb.WriteString("        _init_completion || return\n")
sb.WriteString("    else\n")
sb.WriteString("        words=(\"${COMP_WORDS[@]}\")\n")
sb.WriteString("        cword=$COMP_CWORD\n")
sb.WriteString("        cur=\"${COMP_WORDS[COMP_CWORD]}\"\n")
sb.WriteString("        prev=\"${COMP_WORDS[COMP_CWORD-1]}\"\n")
sb.WriteString("    fi\n\n")

// Build top-level command list (names + aliases).
var allNames []string
for _, c := range completionCommands {
allNames = append(allNames, c.name)
allNames = append(allNames, c.aliases...)
}
fmt.Fprintf(&sb, "    local commands=%q\n\n", strings.Join(allNames, " "))

// Emit subcommand list variables.
for _, c := range completionCommands {
if len(c.subs) == 0 {
continue
}
var subNames []string
for _, s := range c.subs {
subNames = append(subNames, s.name)
}
varName := bashVarName(c.name)
fmt.Fprintf(&sb, "    local %s_cmds=%q\n", varName, strings.Join(subNames, " "))
for _, n := range c.nested {
var nestedNames []string
for _, ns := range n.subs {
nestedNames = append(nestedNames, ns.name)
}
nestedVar := varName + "_" + bashVarName(n.parentSub)
fmt.Fprintf(&sb, "    local %s_cmds=%q\n", nestedVar, strings.Join(nestedNames, " "))
}
}

// Main completion logic: skip global flags to find the actual command.
sb.WriteString(`
    # Find main command and subcommand, skipping global flags and their values.
    local main_cmd="" sub_cmd=""
    local i
    for ((i = 1; i < cword; i++)); do
        local w="${words[$i]}"
        case "$w" in
            -p|--profile|-g|--group|--select|--color)
                ((i++)) # skip the flag's value argument
                ;;
            -*)
                : # skip standalone flags with no value
                ;;
            *)
                if [[ -z "$main_cmd" ]]; then
                    main_cmd="$w"
                elif [[ -z "$sub_cmd" ]]; then
                    sub_cmd="$w"
                fi
                ;;
        esac
    done

    # No main command typed yet: complete top-level commands.
    if [[ -z "$main_cmd" ]]; then
        COMPREPLY=($(compgen -W "$commands" -- "$cur"))
        return
    fi

    # Main command seen but no subcommand yet: complete subcommands.
    if [[ -z "$sub_cmd" ]]; then
        case "$main_cmd" in
`)

// Per-command subcommand cases.
for _, c := range completionCommands {
if len(c.subs) == 0 {
continue
}
varName := bashVarName(c.name)
names := append([]string{c.name}, c.aliases...)
fmt.Fprintf(&sb, "            %s)\n", strings.Join(names, "|"))
fmt.Fprintf(&sb, "                COMPREPLY=($(compgen -W \"$%s_cmds\" -- \"$cur\"))\n", varName)
sb.WriteString("                ;;\n")
}

sb.WriteString(`            add|launch|try)
                COMPREPLY=($(compgen -d -- "$cur"))
                ;;
        esac
        return
    fi

    # Both main command and subcommand are known: handle nested subcommands.
    case "$main_cmd" in
`)

// Nested subcommand cases.
for _, c := range completionCommands {
if len(c.nested) == 0 {
continue
}
varName := bashVarName(c.name)
names := append([]string{c.name}, c.aliases...)
fmt.Fprintf(&sb, "        %s)\n", strings.Join(names, "|"))
sb.WriteString("            case \"$sub_cmd\" in\n")
for _, n := range c.nested {
nestedVar := varName + "_" + bashVarName(n.parentSub)
fmt.Fprintf(&sb, "                %s)\n", n.parentSub)
fmt.Fprintf(&sb, "                    COMPREPLY=($(compgen -W \"$%s_cmds\" -- \"$cur\"))\n", nestedVar)
sb.WriteString("                    ;;\n")
}
sb.WriteString("            esac\n")
sb.WriteString("            ;;\n")
}

sb.WriteString(`    esac
}

complete -F _agent_deck_completions agent-deck
`)

fmt.Print(sb.String())
}

// generateZshCompletion emits a Zsh completion script.
// Each subcommand function uses _arguments -C with a state machine so that
// subcommand suggestions stop after a subcommand has already been entered.
func generateZshCompletion() {
var sb strings.Builder

sb.WriteString("#compdef agent-deck\n\n")
sb.WriteString("# Generated by: agent-deck completions zsh\n\n")

// Top-level function.
sb.WriteString("_agent_deck() {\n")
sb.WriteString("    local line state\n\n")
sb.WriteString("    _arguments -C \\\n")
sb.WriteString("        '(-p --profile)'{-p,--profile}'[Use specific profile]:profile:' \\\n")
sb.WriteString("        '(-g --group)'{-g,--group}'[Launch TUI scoped to group]:group:' \\\n")
sb.WriteString("        '--select[Launch TUI with cursor on session]:session:' \\\n")
sb.WriteString("        '1: :_agent_deck_commands' \\\n")
sb.WriteString("        '*::arg:->args'\n\n")
sb.WriteString("    case $line[1] in\n")
for _, c := range completionCommands {
if len(c.subs) == 0 {
continue
}
funcName := "_agent_deck_" + zshFuncName(c.name)
names := append([]string{c.name}, c.aliases...)
fmt.Fprintf(&sb, "        %s)\n            %s\n            ;;\n",
strings.Join(names, "|"), funcName)
}
sb.WriteString("        add|launch|try)\n            _files -/\n            ;;\n")
sb.WriteString("    esac\n}\n\n")

// Top-level commands list.
sb.WriteString("_agent_deck_commands() {\n    local commands=(\n")
for _, c := range completionCommands {
fmt.Fprintf(&sb, "        %s\n", zshCmdSpec(c.name, c.desc))
for _, alias := range c.aliases {
fmt.Fprintf(&sb, "        %s\n", zshCmdSpec(alias, c.desc+" (alias)"))
}
}
sb.WriteString("    )\n    _describe 'command' commands\n}\n\n")

// Per-command subcommand functions with state machines.
for _, c := range completionCommands {
if len(c.subs) == 0 {
continue
}
funcName := "_agent_deck_" + zshFuncName(c.name)
fmt.Fprintf(&sb, "%s() {\n", funcName)
sb.WriteString("    local state\n")
sb.WriteString("    _arguments -C \\\n")
sb.WriteString("        '1: :->subcmd' \\\n")
sb.WriteString("        '*::arg:->args'\n")
sb.WriteString("    case $state in\n")
sb.WriteString("        subcmd)\n")
sb.WriteString("            local commands=(\n")
for _, sub := range c.subs {
fmt.Fprintf(&sb, "                %s\n", zshCmdSpec(sub.name, sub.desc))
}
sb.WriteString("            )\n")
fmt.Fprintf(&sb, "            _describe '%s command' commands\n", c.name)
if len(c.nested) > 0 {
sb.WriteString("            ;;\n")
sb.WriteString("        args)\n")
sb.WriteString("            case $line[1] in\n")
for _, n := range c.nested {
fmt.Fprintf(&sb, "                %s)\n", n.parentSub)
sb.WriteString("                    local sub_commands=(\n")
for _, ns := range n.subs {
fmt.Fprintf(&sb, "                        %s\n", zshCmdSpec(ns.name, ns.desc))
}
sb.WriteString("                    )\n")
fmt.Fprintf(&sb, "                    _describe '%s %s command' sub_commands\n",
c.name, n.parentSub)
sb.WriteString("                    ;;\n")
}
sb.WriteString("            esac\n")
sb.WriteString("            ;;\n")
} else {
sb.WriteString("            ;;\n")
sb.WriteString("        args) ;;\n")
}
sb.WriteString("    esac\n}\n\n")
}

sb.WriteString("_agent_deck \"$@\"\n")
fmt.Print(sb.String())
}

// generateFishCompletion emits a Fish completion script.
// Top-level commands are only offered when no command has been entered yet
// (__fish_use_subcommand). Subcommands are guarded so they stop appearing once
// a subcommand has already been selected.
func generateFishCompletion() {
var sb strings.Builder

sb.WriteString("# fish completion for agent-deck\n")
sb.WriteString("# Generated by: agent-deck completions fish\n\n")
sb.WriteString("# Remove any existing completions\n")
sb.WriteString("complete -c agent-deck -e\n\n")
sb.WriteString("# Global options\n")
sb.WriteString("complete -c agent-deck -s p -l profile -d 'Use specific profile' -r\n")
sb.WriteString("complete -c agent-deck -s g -l group -d 'Launch TUI scoped to group' -r\n")
sb.WriteString("complete -c agent-deck -l select -d 'Launch TUI with cursor on session' -r\n\n")

sb.WriteString("# Top-level commands (only when no command has been given yet)\n")
for _, c := range completionCommands {
fmt.Fprintf(&sb, "complete -c agent-deck -f -n __fish_use_subcommand -a %s -d '%s'\n",
c.name, fishEscapeDesc(c.desc))
for _, alias := range c.aliases {
fmt.Fprintf(&sb, "complete -c agent-deck -f -n __fish_use_subcommand -a %s -d '%s (alias)'\n",
alias, fishEscapeDesc(c.desc))
}
}
sb.WriteString("\n")

// Subcommands: guarded by "parent seen AND no sibling subcommand seen yet".
for _, c := range completionCommands {
if len(c.subs) == 0 {
continue
}

// Build the sibling-subcommand list for the "not yet selected" guard.
var subNames []string
for _, sub := range c.subs {
subNames = append(subNames, sub.name)
}
notGuard := "not __fish_seen_subcommand_from " + strings.Join(subNames, " ")

// Parent condition includes the command and all its aliases.
allParentNames := append([]string{c.name}, c.aliases...)
parentCond := "__fish_seen_subcommand_from " + strings.Join(allParentNames, " ")
guard := parentCond + "; and " + notGuard

fmt.Fprintf(&sb, "# %s subcommands\n", c.name)
for _, sub := range c.subs {
fmt.Fprintf(&sb, "complete -c agent-deck -f -n '%s' -a %s -d '%s'\n",
guard, sub.name, fishEscapeDesc(sub.desc))
}
sb.WriteString("\n")

// Nested sub-subcommands.
for _, n := range c.nested {
var nestedNames []string
for _, ns := range n.subs {
nestedNames = append(nestedNames, ns.name)
}
nestedNotGuard := "not __fish_seen_subcommand_from " + strings.Join(nestedNames, " ")
nestedGuard := parentCond + "; and __fish_seen_subcommand_from " + n.parentSub +
"; and " + nestedNotGuard

fmt.Fprintf(&sb, "# %s %s subcommands\n", c.name, n.parentSub)
for _, ns := range n.subs {
fmt.Fprintf(&sb, "complete -c agent-deck -f -n '%s' -a %s -d '%s'\n",
nestedGuard, ns.name, fishEscapeDesc(ns.desc))
}
sb.WriteString("\n")
}
}

fmt.Print(sb.String())
}

// bashVarName converts a command name to a valid Bash variable name fragment
// by replacing hyphens with underscores.
func bashVarName(s string) string {
return strings.ReplaceAll(s, "-", "_")
}

// zshFuncName converts a command name to a Zsh function name fragment.
func zshFuncName(s string) string {
return strings.ReplaceAll(s, "-", "_")
}

// zshCmdSpec formats a Zsh completion spec string 'name:description',
// escaping single quotes in the description.
func zshCmdSpec(name, desc string) string {
escaped := strings.ReplaceAll(desc, "'", "'\\''")
return fmt.Sprintf("'%s:%s'", name, escaped)
}

// fishEscapeDesc escapes single quotes in Fish completion descriptions.
func fishEscapeDesc(s string) string {
return strings.ReplaceAll(s, "'", "\\'")
}
