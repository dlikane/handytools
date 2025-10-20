package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandWildcards_NoMatches(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	pattern := filepath.Join(tmp, "*.jpg")

	got := ExpandWildcards([]string{pattern})
	if len(got) != 0 {
		t.Fatalf("expected no matches, got %v", got)
	}
}

func TestExpandWildcards_SinglePattern_MatchesFiles(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	f1 := filepath.Join(tmp, "a.txt")
	f2 := filepath.Join(tmp, "b.txt")
	os.WriteFile(f1, []byte("a"), 0644)
	os.WriteFile(f2, []byte("b"), 0644)

	got := ExpandWildcards([]string{filepath.Join(tmp, "*.txt")})
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d (%v)", len(got), got)
	}

	wantSet := map[string]bool{f1: true, f2: true}
	for _, g := range got {
		if !wantSet[g] {
			t.Fatalf("unexpected file: %s", g)
		}
	}
}

func TestExpandWildcards_MultiplePatterns_CombinedResults(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	f1 := filepath.Join(tmp, "x.log")
	f2 := filepath.Join(tmp, "y.txt")
	os.WriteFile(f1, []byte("x"), 0644)
	os.WriteFile(f2, []byte("y"), 0644)

	got := ExpandWildcards([]string{
		filepath.Join(tmp, "*.log"),
		filepath.Join(tmp, "*.txt"),
	})

	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d (%v)", len(got), got)
	}

	found := make(map[string]bool)
	for _, f := range got {
		found[f] = true
	}
	if !found[f1] || !found[f2] {
		t.Fatalf("expected both %s and %s, got %v", f1, f2, got)
	}
}

func TestExpandWildcards_InvalidPattern_LogsAndContinues(t *testing.T) {
	t.Parallel()

	// filepath.Glob returns error for malformed patterns like "["
	got := ExpandWildcards([]string{"["})
	if len(got) != 0 {
		t.Fatalf("expected no matches for invalid pattern, got %v", got)
	}
}
