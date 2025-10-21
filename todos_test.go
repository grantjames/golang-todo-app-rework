package todos

import (
	"fmt"
	"os"
	"testing"
)

func TestStoreConcurrency(t *testing.T) {
	// Clear out test file first
	testFile := "test_todos.json"
	f, _ := os.Create(testFile)
	f.Write([]byte("{}"))
	defer os.Remove(testFile)
	f.Close()

	//store := NewStore(testFile)
	StartStore(testFile)
	const N = 100
	ids := make([]string, N)

	t.Run("Creating todos", func(t *testing.T) {
		for i := 0; i < N; i++ {
			i := i
			t.Run("Testing parallel add", func(t *testing.T) {
				t.Parallel()
				ids[i] = Create(fmt.Sprintf("task-%d", i))
			})
		}
	})

	if got := len(List()); got != N {
		t.Fatalf("want %d todos, got %d", N, got)
	}

	t.Run("deletes", func(t *testing.T) {
		for i := N - 50; i < N; i++ {
			id := ids[i]
			t.Run(fmt.Sprintf("delete-%s", id), func(t *testing.T) {
				t.Parallel()
				if !Delete(id) {
					t.Fatalf("failed to delete %s", id)
				}
			})
		}
	})

	if got := len(List()); got != N-50 {
		t.Fatalf("final size: want %d, got %d", N-200, got)
	}
}
