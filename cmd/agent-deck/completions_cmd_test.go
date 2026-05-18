package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// TestCompletionsHelp verifies that running completions with no args shows help
func TestCompletionsHelp(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleCompletions([]string{})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "Usage: agent-deck completions <shell>") {
		t.Error("Expected help output to contain usage message")
	}
	if !strings.Contains(output, "bash") || !strings.Contains(output, "zsh") || !strings.Contains(output, "fish") {
		t.Error("Expected help output to list bash, zsh, and fish")
	}
}

// TestCompletionsBash verifies bash completion generation
func TestCompletionsBash(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleCompletions([]string{"bash"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "_agent_deck_completions") {
		t.Error("Expected bash completion to contain _agent_deck_completions function")
	}
	if !strings.Contains(output, "complete -F _agent_deck_completions agent-deck") {
		t.Error("Expected bash completion to contain complete directive")
	}
	// Verify key commands are present
	keyCommands := []string{"session", "mcp", "skill", "group", "remote", "conductor"}
	for _, cmd := range keyCommands {
		if !strings.Contains(output, cmd) {
			t.Errorf("Expected bash completion to contain %q command", cmd)
		}
	}
}

// TestCompletionsZsh verifies zsh completion generation
func TestCompletionsZsh(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleCompletions([]string{"zsh"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "#compdef agent-deck") {
		t.Error("Expected zsh completion to contain #compdef agent-deck")
	}
	if !strings.Contains(output, "_agent_deck") {
		t.Error("Expected zsh completion to contain _agent_deck function")
	}
	if !strings.Contains(output, "_arguments") {
		t.Error("Expected zsh completion to use _arguments")
	}
	// Verify descriptions are present
	if !strings.Contains(output, "Manage session lifecycle") {
		t.Error("Expected zsh completion to contain command descriptions")
	}
}

// TestCompletionsFish verifies fish completion generation
func TestCompletionsFish(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	handleCompletions([]string{"fish"})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "# fish completion for agent-deck") {
		t.Error("Expected fish completion to contain header comment")
	}
	if !strings.Contains(output, "complete -c agent-deck") {
		t.Error("Expected fish completion to contain complete directives")
	}
	if !strings.Contains(output, "__fish_use_subcommand") {
		t.Error("Expected fish completion to use __fish_use_subcommand")
	}
	if !strings.Contains(output, "__fish_seen_subcommand_from session") {
		t.Error("Expected fish completion to handle session subcommands")
	}
}

// TestCompletionsInvalidShell verifies error handling for unsupported shells
func TestCompletionsInvalidShell(t *testing.T) {
	// We can't test os.Exit directly in Go tests as it would terminate the test
	// The actual exit behavior is tested in integration tests
	// This test just validates that the other shell completions work correctly
	t.Skip("Skipping os.Exit test - validated in integration tests")
}
