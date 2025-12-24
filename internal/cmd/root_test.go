package cmd

import (
	"testing"
)

func TestRootCmdFlags(t *testing.T) {
	// Check persistent flags
	themeFlag := RootCmd.PersistentFlags().Lookup("theme")
	if themeFlag == nil {
		t.Fatal("Theme flag not found")
	}
	if themeFlag.DefValue != "base16" {
		t.Errorf("Expected default theme base16, got %s", themeFlag.DefValue)
	}

	profileFlag := RootCmd.PersistentFlags().Lookup("profile")
	if profileFlag == nil {
		t.Fatal("Profile flag not found")
	}
	if profileFlag.DefValue != "false" {
		t.Errorf("Expected default profile false, got %s", profileFlag.DefValue)
	}
}

func TestInitConfig(t *testing.T) {
	// Just ensure it doesn't panic
	initConfig()
}
