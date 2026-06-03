package cmd

import (
	"strings"
	"testing"
)

func TestUsageDocumentsLongFlags(t *testing.T) {
	text := usageText()

	for _, want := range []string{
		"--addr ADDR",
		"--html",
		"--version",
		"-o PATH",
		"default 0.0.0.0:3000",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("usageText() missing %q:\n%s", want, text)
		}
	}

	for _, notWant := range []string{"  -addr ADDR", "  -html        "} {
		if strings.Contains(text, notWant) {
			t.Fatalf("usageText() contains old flag spelling %q:\n%s", notWant, text)
		}
	}
}

func TestExecuteInvalidFlagIncludesUsage(t *testing.T) {
	err := Execute([]string{"--not-a-real-flag"}, "test")
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid flag error")
	}

	text := err.Error()
	if !strings.Contains(text, "flag provided but not defined") {
		t.Fatalf("Execute() error missing flag parse detail: %v", err)
	}
	if !strings.Contains(text, "Usage:") {
		t.Fatalf("Execute() error missing usage text: %v", err)
	}
}
