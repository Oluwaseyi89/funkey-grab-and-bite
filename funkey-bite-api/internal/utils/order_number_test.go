package utils

import (
	"strconv"
	"strings"
	"sync"
	"testing"
)

// TestGenerateOrderNumberUniquenessUnderConcurrency spawns 500 goroutines
// that all call GenerateOrderNumber simultaneously and asserts that every
// returned value is unique — catching the previous UnixNano%10000 collision.
func TestGenerateOrderNumberUniquenessUnderConcurrency(t *testing.T) {
	const n = 500

	results := make([]string, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			results[i] = GenerateOrderNumber()
		}()
	}
	wg.Wait()

	seen := make(map[string]struct{}, n)
	for _, num := range results {
		if _, dup := seen[num]; dup {
			t.Errorf("duplicate order number generated under concurrency: %s", num)
		}
		seen[num] = struct{}{}
	}
}

// TestGenerateOrderNumberFormat verifies that every number matches the
// expected FG-YYYY-MM-<seq> prefix pattern.
func TestGenerateOrderNumberFormat(t *testing.T) {
	num := GenerateOrderNumber()
	if !strings.HasPrefix(num, "FG-") {
		t.Errorf("expected order number to start with FG-, got %s", num)
	}
	parts := strings.Split(num, "-")
	// Expected parts: ["FG", "YYYY", "MM", "<seq>"]
	if len(parts) != 4 {
		t.Errorf("expected 4 dash-separated parts, got %d in %q", len(parts), num)
	}
}

// TestGenerateOrderNumberMonotonicallyIncreases verifies that sequential
// calls return strictly increasing sequence numbers by parsing the numeric
// trailing segment.
func TestGenerateOrderNumberMonotonicallyIncreases(t *testing.T) {
	seqOf := func(num string) int64 {
		parts := strings.Split(num, "-")
		if len(parts) != 4 {
			t.Fatalf("unexpected format: %s", num)
		}
		v, err := strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			t.Fatalf("non-numeric sequence in %s: %v", num, err)
		}
		return v
	}

	prev := seqOf(GenerateOrderNumber())
	for i := 0; i < 100; i++ {
		next := seqOf(GenerateOrderNumber())
		if next <= prev {
			t.Errorf("sequence did not increase: prev=%d next=%d", prev, next)
		}
		prev = next
	}
}
