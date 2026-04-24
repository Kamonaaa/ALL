// main_test.go
package main

import (
	"os"
	"reflect"
	"testing"
)

// ---------- Unit Tests ----------

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "simple words",
			input: "hello world",
			want:  []string{"hello", "world"},
		},
		{
			name:  "punctuation",
			input: "Hello, world!",
			want:  []string{"Hello", ",", "world", "!"},
		},
		{
			name:  "modifier token",
			input: "Hello (up)",
			want:  []string{"Hello", "(up)"},
		},
		{
			name:  "contraction",
			input: "don't",
			want:  []string{"don't"},
		},
		{
			name:  "quotes",
			input: "'hello'",
			want:  []string{"'", "hello", "'"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tokenize(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

// Test the modifier pipeline on a small sentence
func TestProcessAndFormat(t *testing.T) {
	input := "hello world (up) and universe (cap, 1)"
	tokens := tokenize(input)
	processed := process(tokens)
	result := format(processed)
	expected := "HELLO WORLD and Universe"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

// ---------- Integration Test with Golden File ----------

func TestFullProgram(t *testing.T) {
	// Save original args to restore later
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	// Create a temporary output file
	tmpFile, err := os.CreateTemp("", "textmod-output-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name()) // clean up after test
	tmpFile.Close()                 // main will open it via WriteFile

	// Simulate command-line arguments: program name, input file, output file
	os.Args = []string{"textmod", "testdata/input.txt", tmpFile.Name()}

	// Run the actual main
	main()

	// Read the generated output
	gotBytes, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)

	// Read the expected output
	wantBytes, err := os.ReadFile("testdata/expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := string(wantBytes)

	if got != want {
		// Write the actual output to a file for manual inspection
		debugFile := "testdata/actual_output.txt"
		os.WriteFile(debugFile, gotBytes, 0644)
		t.Errorf("output mismatch.\ngot:  %s\nwant: %s\nActual output saved to %s for inspection.",
			got, want, debugFile)
	}
}

// ---------- Optional: Helper to See Output in a Text File ----------
// This test simply runs the program and writes the result to a known file,
// so you can open it and see the result. It's not a strict pass/fail test.
func TestOutputToFileForVisualInspection(t *testing.T) {
	input := "hello world (up) and universe (cap, 1)"
	tokens := tokenize(input)
	processed := process(tokens)
	result := format(processed)

	// Write to a file you can open manually
	outPath := "testdata/visual_inspection.txt"
	err := os.WriteFile(outPath, []byte(result), 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Output written to %s: %s", outPath, result)
}
