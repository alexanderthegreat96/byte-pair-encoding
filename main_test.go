package main

import (
	"strings"
	"testing"
)

func TestVocabularyInitialization(t *testing.T) {
	v := NewVocabulary(VOCABULARY_CAPACITY)

	if len(v.data) < VOCABULARY_CAPACITY {
		t.Errorf("Expected at least 256 base vocabulary tokens, got %d", len(v.data))
	}

	tokenBytes := v.GetTokenAtId(65)
	if tokenBytes == nil || len(tokenBytes) != 1 || tokenBytes[0] != 65 {
		t.Errorf("Expected token 65 to map to ASCII 'A' [65], got %v", tokenBytes)
	}
}

func TestVocabularyTrainingAndEncoding(t *testing.T) {
	v := NewVocabulary(VOCABULARY_CAPACITY)

	corpus := "lower lower lower lowest lowest lowest"

	v.Train(corpus, 265)

	if len(v.merges) == 0 {
		t.Error("Expected vocabulary training to generate merge rules, but none found")
	}

	encoded := v.Encode("lower")
	if len(encoded) == 0 {
		t.Error("Expected non-empty slice of token IDs from Encode")
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	v := NewVocabulary(VOCABULARY_CAPACITY)

	if tokens := v.Encode(""); tokens != nil {
		t.Errorf("Expected nil for empty string encoding, got %v", tokens)
	}

	if decoded := v.Decode([]int{}); decoded != "" {
		t.Errorf("Expected empty string for empty token slice decoding, got %q", decoded)
	}

	original := "hello"
	tokens := v.Encode(original)
	decoded := v.Decode(tokens)

	if decoded != original {
		t.Errorf("Roundtrip failed. Expected %q, got %q", original, decoded)
	}

	if !strings.HasPrefix(v.Decode([]int{104, 101}), "he") {
		t.Errorf("Expected decoded tokens [104, 101] to start with 'he', got %q", v.Decode([]int{104, 101}))
	}
}
