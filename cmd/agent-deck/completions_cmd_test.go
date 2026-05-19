package main

import (
"bytes"
"io"
"os"
"strings"
"testing"
)

// captureStdout runs fn and returns whatever it printed to os.Stdout.
func captureStdout(fn func()) string {
oldStdout := os.Stdout
r, w, _ := os.Pipe()
os.Stdout = w
fn()
w.Close()
os.Stdout = oldStdout
var buf bytes.Buffer
io.Copy(&buf, r) //nolint:errcheck
return buf.String()
}

// TestCompletionsHelp verifies that no-arg and help-flag invocations print usage.
func TestCompletionsHelp(t *testing.T) {
for _, args := range [][]string{{}, {"help"}, {"-h"}, {"--help"}} {
args := args
out := captureStdout(func() { handleCompletions(args) })
if !strings.Contains(out, "Usage: agent-deck completions <shell>") {
t.Errorf("args=%v: expected usage message, got: %q", args, out)
}
for _, shell := range []string{"bash", "zsh", "fish"} {
if !strings.Contains(out, shell) {
t.Errorf("args=%v: expected %q in help output", args, shell)
}
}
}
}

// TestCompletionsInvalidShell verifies that an unsupported shell returns exit code 1.
func TestCompletionsInvalidShell(t *testing.T) {
code := runCompletions([]string{"powershell"})
if code != 1 {
t.Errorf("expected exit code 1 for unsupported shell, got %d", code)
}
// Valid shells must return 0.
for _, shell := range []string{"bash", "zsh", "fish"} {
shell := shell
var code int
captureStdout(func() { code = runCompletions([]string{shell}) })
if code != 0 {
t.Errorf("expected exit code 0 for %s, got %d", shell, code)
}
}
}

// TestCompletionsBash verifies structure and correct subcommand lists in Bash output.
func TestCompletionsBash(t *testing.T) {
out := captureStdout(func() { handleCompletions([]string{"bash"}) })

// Structural checks.
mustContain(t, out, "_agent_deck_completions", "Bash function name")
mustContain(t, out, "complete -F _agent_deck_completions agent-deck", "complete directive")
mustContain(t, out, "_init_completion", "_init_completion call")
mustContain(t, out, "COMP_WORDS", "COMP_WORDS fallback")
mustContain(t, out, "--profile", "global --profile flag skip")

// Top-level commands must be present.
for _, cmd := range []string{"session", "mcp", "skill", "plugin", "group", "worktree",
"remote", "conductor", "profile", "codex-hooks", "gemini-hooks"} {
mustContain(t, out, cmd, "top-level command "+cmd)
}

// session: valid subcommands from handleSession.
validSessionSubs := []string{
"start", "stop", "remove", "restart", "revive", "fork", "attach",
"show", "current", "set-parent", "unset-parent", "update",
"set-transition-notify", "set-title-lock", "set", "move", "mv",
"send", "output", "search",
}
for _, sub := range validSessionSubs {
mustContain(t, out, sub, "session subcmd "+sub)
}
// session: stale commands from the old list must not appear in the session_cmds variable.
for _, bad := range []string{"stream", "color", "info", "set-telegram-warnings"} {
mustNotContainInVar(t, out, "session_cmds", bad, "stale session subcmd "+bad)
}

// mcp: must include server and ls.
for _, sub := range []string{"server", "ls", "list", "attached", "attach", "detach"} {
mustContain(t, out, sub, "mcp subcmd "+sub)
}

// plugin: valid subcommands; old invalid list must be gone from plugin_cmds var.
for _, sub := range []string{"attached", "attach", "detach"} {
mustContain(t, out, sub, "plugin subcmd "+sub)
}
for _, bad := range []string{"install", "uninstall", "update", "info"} {
mustNotContainInVar(t, out, "plugin_cmds", bad, "stale plugin subcmd "+bad)
}

// group: valid subcommands; old invalid ones must not appear in group_cmds var.
for _, sub := range []string{"create", "new", "update", "set", "delete", "rm",
"remove", "move", "mv", "change", "reparent", "reorder", "sort"} {
mustContain(t, out, sub, "group subcmd "+sub)
}
for _, bad := range []string{"rename", "set-parent"} {
mustNotContainInVar(t, out, "group_cmds", bad, "stale group subcmd "+bad)
}

// worktree: finish must be present.
mustContain(t, out, "finish", "worktree subcmd finish")

// conductor: move must be present in conductor_cmds.
mustContainInVar(t, out, "conductor_cmds", "move", "conductor subcmd move")

// profile: no migrate; ls, new, rm, default must be present.
for _, sub := range []string{"ls", "new", "rm", "default"} {
mustContain(t, out, sub, "profile subcmd "+sub)
}
mustNotContainInVar(t, out, "profile_cmds", "migrate", "stale profile subcmd migrate")

// skill source: no update; ls, add, remove, rm must be present.
for _, sub := range []string{"add", "remove", "rm"} {
mustContain(t, out, sub, "skill source subcmd "+sub)
}
mustNotContainInVar(t, out, "skill_source_cmds", "update", "stale skill source subcmd update")
}

// TestCompletionsZsh verifies structure and subcommand correctness in Zsh output.
func TestCompletionsZsh(t *testing.T) {
out := captureStdout(func() { handleCompletions([]string{"zsh"}) })

mustContain(t, out, "#compdef agent-deck", "#compdef header")
mustContain(t, out, "_agent_deck", "_agent_deck function")
mustContain(t, out, "_arguments", "_arguments usage")
mustContain(t, out, "Manage session lifecycle", "session description")

// State machine: each subcommand function must use _arguments -C with ->subcmd state.
mustContain(t, out, "->subcmd", "Zsh state machine subcmd state")
mustContain(t, out, "->args", "Zsh state machine args state")

// Verify correct session subcommands appear in the _agent_deck_session block.
sessionBlock := extractZshFuncBlock(out, "_agent_deck_session")
for _, sub := range []string{"start", "stop", "remove", "revive", "current",
"set-parent", "unset-parent", "output", "set-transition-notify"} {
if !strings.Contains(sessionBlock, sub) {
t.Errorf("zsh _agent_deck_session missing subcmd %q", sub)
}
}
// Stale session subcommands must not appear as completion entries in the session block.
for _, bad := range []string{"stream", "color", "set-telegram-warnings"} {
if strings.Contains(sessionBlock, "'"+bad+":") {
t.Errorf("zsh: stale session subcmd %q should not appear in _agent_deck_session block", bad)
}
}
// 'info' is a valid worktree subcommand; verify it is absent from the session block only.
if strings.Contains(sessionBlock, "'info:") {
t.Errorf("zsh: 'info' should not appear as a session subcommand entry in _agent_deck_session block")
}

// conductor move must be present.
mustContain(t, out, "move:Move conductor", "conductor move subcommand")

// worktree finish must be present.
mustContain(t, out, "finish:Mark worktree", "worktree finish subcommand")
}

// TestCompletionsFish verifies structure and subcommand correctness in Fish output.
func TestCompletionsFish(t *testing.T) {
out := captureStdout(func() { handleCompletions([]string{"fish"}) })

mustContain(t, out, "# fish completion for agent-deck", "header comment")
mustContain(t, out, "complete -c agent-deck", "complete directive")
mustContain(t, out, "__fish_use_subcommand", "top-level guard")
mustContain(t, out, "__fish_seen_subcommand_from session", "session guard")

// Previously missing top-level commands must now be present.
for _, cmd := range []string{"notify-daemon", "hook-handler", "codex-notify",
"hooks", "uninstall", "debug-dump"} {
mustContain(t, out, cmd, "Fish top-level command "+cmd)
}

// Subcommand guards must include "not __fish_seen_subcommand_from" to stop
// suggesting subcommands once one is already selected.
mustContain(t, out, "not __fish_seen_subcommand_from", "Fish not-guard for subcommands")

// Correct session subcommands: each should appear on a line guarded by session.
for _, sub := range []string{"start", "stop", "remove", "revive", "current",
"set-parent", "unset-parent", "output"} {
mustContainOnLineWith(t, out, "session", " -a "+sub, "fish session subcmd "+sub)
}

// conductor move must be present on a conductor-guarded line.
mustContainOnLineWith(t, out, "conductor", " -a move", "fish conductor move")

// worktree finish must be present on a worktree-guarded line.
mustContainOnLineWith(t, out, "worktree", " -a finish", "fish worktree finish")

// plugin: stale subcommands must not appear on plugin-guarded lines.
for _, bad := range []string{"install", "info"} {
for _, line := range strings.Split(out, "\n") {
if strings.Contains(line, "__fish_seen_subcommand_from plugin") &&
strings.Contains(line, " -a "+bad) {
t.Errorf("fish: stale plugin subcmd %q should not appear in plugin-guarded line: %s", bad, line)
}
}
}
}

// mustContain fails the test if out does not contain substr.
func mustContain(t *testing.T, out, substr, label string) {
t.Helper()
if !strings.Contains(out, substr) {
t.Errorf("%s: expected %q in output", label, substr)
}
}

// mustNotContainInVar asserts that the variable-assignment line for varName
// does not include bad. This is more precise than a full-string search because
// it only looks at the specific variable's value.
func mustNotContainInVar(t *testing.T, out, varName, bad, label string) {
t.Helper()
for _, line := range strings.Split(out, "\n") {
if strings.Contains(line, varName+"=") && strings.Contains(line, bad) {
t.Errorf("%s: found stale entry %q in %s variable line: %s", label, bad, varName, line)
}
}
}

// mustContainInVar asserts that the variable-assignment line for varName contains substr.
func mustContainInVar(t *testing.T, out, varName, substr, label string) {
t.Helper()
for _, line := range strings.Split(out, "\n") {
if strings.Contains(line, varName+"=") && strings.Contains(line, substr) {
return
}
}
t.Errorf("%s: %q not found in %s variable line", label, substr, varName)
}

// mustContainOnLineWith asserts that at least one line containing lineKey also
// contains substr. Used to verify that a subcommand appears in the right section.
func mustContainOnLineWith(t *testing.T, out, lineKey, substr, label string) {
t.Helper()
for _, line := range strings.Split(out, "\n") {
if strings.Contains(line, lineKey) && strings.Contains(line, substr) {
return
}
}
t.Errorf("%s: no line containing both %q and %q found", label, lineKey, substr)
}

// extractZshFuncBlock returns the text of the named Zsh function body,
// from the function declaration to its closing "}\n".
func extractZshFuncBlock(out, funcName string) string {
start := strings.Index(out, funcName+"() {")
if start == -1 {
return ""
}
rest := out[start:]
// Find the end of this function (first standalone "}" on its own line after opening).
end := strings.Index(rest, "\n}\n")
if end == -1 {
return rest
}
return rest[:end+3]
}
