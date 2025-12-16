package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "mcpr") {
		t.Error("expected output to contain 'mcpr'")
	}
	if !strings.Contains(output, "add") {
		t.Error("expected output to contain 'add' command")
	}
	if !strings.Contains(output, "client") {
		t.Error("expected output to contain 'client' command")
	}
	if !strings.Contains(output, "list") {
		t.Error("expected output to contain 'list' command")
	}
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	cmds := rootCmd.Commands()
	cmdNames := make(map[string]bool)
	for _, cmd := range cmds {
		cmdNames[cmd.Name()] = true
	}

	expectedCmds := []string{"add", "client", "list", "completion", "help"}
	for _, name := range expectedCmds {
		if !cmdNames[name] {
			t.Errorf("expected subcommand %q to be present", name)
		}
	}
}

func TestAddCmd_Structure(t *testing.T) {
	if addCmd.Use != "add" {
		t.Errorf("expected Use to be 'add', got %q", addCmd.Use)
	}

	if addCmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	if addCmd.Long == "" {
		t.Error("expected Long description to be set")
	}
}

func TestAddCmd_HasSubcommands(t *testing.T) {
	cmds := addCmd.Commands()
	cmdNames := make(map[string]bool)
	for _, cmd := range cmds {
		cmdNames[cmd.Name()] = true
	}

	expectedCmds := []string{"stdio", "http"}
	for _, name := range expectedCmds {
		if !cmdNames[name] {
			t.Errorf("expected subcommand %q to be present", name)
		}
	}
}

func TestAddCmd_PersistentFlags(t *testing.T) {
	flags := addCmd.PersistentFlags()

	flag := flags.Lookup("local")
	if flag == nil {
		t.Error("expected persistent flag 'local' to exist")
	} else if flag.Shorthand != "l" {
		t.Errorf("expected shorthand 'l' for flag 'local', got %q", flag.Shorthand)
	}
}

func TestClientCmd_Structure(t *testing.T) {
	if clientCmd.Use != "client" {
		t.Errorf("expected Use to be 'client', got %q", clientCmd.Use)
	}

	if clientCmd.Short == "" {
		t.Error("expected Short description to be set")
	}
}

func TestClientCmd_HasSubcommands(t *testing.T) {
	cmds := clientCmd.Commands()
	cmdNames := make(map[string]bool)
	for _, cmd := range cmds {
		cmdNames[cmd.Name()] = true
	}

	expectedCmds := []string{"sync", "remove"}
	for _, name := range expectedCmds {
		if !cmdNames[name] {
			t.Errorf("expected subcommand %q to be present", name)
		}
	}
}

func TestClientSyncCmd_Structure(t *testing.T) {
	if clientSyncCmd.Use != "sync [client-name]" {
		t.Errorf("expected Use to be 'sync [client-name]', got %q", clientSyncCmd.Use)
	}

	if clientSyncCmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	if clientSyncCmd.Long == "" {
		t.Error("expected Long description to be set")
	}

	// Check that Long contains supported clients
	supportedClients := []string{"claude-desktop", "claude-code", "cursor", "windsurf", "zed", "cline", "vscode", "continue", "codex", "gemini", "kilo-code", "zencoder"}
	for _, client := range supportedClients {
		if !strings.Contains(clientSyncCmd.Long, client) {
			t.Errorf("expected Long to mention %q", client)
		}
	}
}

func TestClientSyncCmd_Flags(t *testing.T) {
	flags := clientSyncCmd.Flags()

	testCases := []struct {
		name      string
		shorthand string
	}{
		{"servers", "s"},
		{"local", "l"},
	}

	for _, tc := range testCases {
		flag := flags.Lookup(tc.name)
		if flag == nil {
			t.Errorf("expected flag %q to exist", tc.name)
			continue
		}
		if flag.Shorthand != tc.shorthand {
			t.Errorf("expected shorthand %q for flag %q, got %q", tc.shorthand, tc.name, flag.Shorthand)
		}
	}
}

func TestListCmd_Structure(t *testing.T) {
	if listCmd.Use != "list" {
		t.Errorf("expected Use to be 'list', got %q", listCmd.Use)
	}

	if listCmd.Short == "" {
		t.Error("expected Short description to be set")
	}
}

func TestListCmd_Flags(t *testing.T) {
	flags := listCmd.Flags()

	flag := flags.Lookup("clients")
	if flag == nil {
		t.Error("expected flag 'clients' to exist")
	} else if flag.Shorthand != "c" {
		t.Errorf("expected shorthand 'c' for flag 'clients', got %q", flag.Shorthand)
	}
}

func TestRemoveCmd_Structure(t *testing.T) {
	if removeCmd.Use != "remove [server-name]" {
		t.Errorf("expected Use to be 'remove [server-name]', got %q", removeCmd.Use)
	}

	if removeCmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	// Check aliases
	aliases := removeCmd.Aliases
	hasRmAlias := false
	for _, a := range aliases {
		if a == "rm" {
			hasRmAlias = true
			break
		}
	}
	if !hasRmAlias {
		t.Error("expected 'rm' alias to be present")
	}
}

func TestClientRemoveCmd_Structure(t *testing.T) {
	if clientRemoveCmd.Use != "remove [client-name]" {
		t.Errorf("expected Use to be 'remove [client-name]', got %q", clientRemoveCmd.Use)
	}

	if clientRemoveCmd.Short == "" {
		t.Error("expected Short description to be set")
	}
}

func TestClientRemoveCmd_Flags(t *testing.T) {
	flags := clientRemoveCmd.Flags()

	flag := flags.Lookup("local")
	if flag == nil {
		t.Error("expected flag 'local' to exist")
	} else if flag.Shorthand != "l" {
		t.Errorf("expected shorthand 'l' for flag 'local', got %q", flag.Shorthand)
	}
}
