package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunReadsFromStdin(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run(nil, strings.NewReader("field:\nattribute = value"), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("got exit code %d stderr %q", exitCode, stderr.String())
	}

	if got := stdout.String(); got != "field:\nattribute = value" {
		t.Fatalf("got %q", got)
	}
}

func TestRunReportsParseErrors(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run(nil, strings.NewReader("= value"), &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("got exit code %d", exitCode)
	}
	if stdout.Len() != 0 {
		t.Fatalf("unexpected stdout %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "has no key") {
		t.Fatalf("got stderr %q", stderr.String())
	}
}
