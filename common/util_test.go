package common

import (
	"testing"
)

func TestHash(t *testing.T) {
	s := "abc"
	h := Hash(s)

	// Expected rolling hash values
	// hash[0] = 0
	// hash[1] = 0*131 + 'a' = 97
	// hash[2] = 97*131 + 'b' = 12805
	// hash[3] = 12805*131 + 'c' = 1677554
	expected := []uint64{0, 97, 12805, 1677554}

	if len(h.hash) != len(expected) {
		t.Fatalf("expected hash length %d, got %d", len(expected), len(h.hash))
	}

	for i := range expected {
		if h.hash[i] != expected[i] {
			t.Errorf("at index %d, expected %d, got %d", i, expected[i], h.hash[i])
		}
	}
}

func TestGetRangeHash(t *testing.T) {
	s := "abcde"
	h := Hash(s)

	tests := []struct {
		l, r     int
		expected uint64
	}{
		{0, 0, uint64('a')},
		{0, 1, uint64('a')*BASE + uint64('b')},
		{1, 3, (uint64('b')*quickPow(BASE, 2) + uint64('c')*BASE + uint64('d'))},
		{2, 4, (uint64('c')*quickPow(BASE, 2) + uint64('d')*BASE + uint64('e'))},
		{3, 3, uint64('d')},
	}

	for _, test := range tests {
		got := GetRangeHash(&h, test.l, test.r)
		if got != test.expected {
			t.Errorf("GetRangeHash(s, %d, %d): expected %d, got %d", test.l, test.r, test.expected, got)
		}
	}
}

func TestQuickPow(t *testing.T) {
	if quickPow(2, 10) != 1024 {
		t.Error("quickPow(2, 10) should be 1024")
	}
	if quickPow(5, 0) != 1 {
		t.Error("quickPow(5, 0) should be 1")
	}
	if quickPow(3, 1) != 3 {
		t.Error("quickPow(3, 1) should be 3")
	}
}


func TestGetReverseComplement(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"AAAA", "TTTT"},
		{"ACTG", "CAGT"},
		{"GCTA", "TAGC"},
		{"", ""},
		{"AT", "AT"},
		{"CG", "CG"},
	}

	for _, test := range tests {
		got := GetReverseComplement(test.input)
		if got != test.expected {
			t.Errorf("GetReverseComplement(%q): expected %q, got %q", test.input, test.expected, got)
		}
	}
}