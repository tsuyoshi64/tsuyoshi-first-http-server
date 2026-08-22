package main

import (
	"testing"
)

func TestCleanProfanity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No profanity",
			input:    "This is a clean chirp about birds",
			expected: "This is a clean chirp about birds",
		},
		{
			name:     "Lowercase profanity matching",
			input:    "This is a kerfuffle opinion",
			expected: "This is a **** opinion",
		},
		{
			name:     "Uppercase profanity matching",
			input:    "I love SHARBERT ice cream on FORNAX",
			expected: "I love **** ice cream on ****",
		},
		{
			name:     "Punctuation should skip matching",
			input:    "Sharbert! That was a kerfuffle.",
			expected: "Sharbert! That was a kerfuffle.",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := cleanProfanity(tt.input)
			if actual != tt.expected {
				t.Errorf("\nFAIL: %s\nInput:    %q\nExpected: %q\nActual:   %q", tt.name, tt.input, tt.expected, actual)
			}
		})
	}
}
