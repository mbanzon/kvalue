package storage

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var (
	smallValue json.RawMessage
	largeValue json.RawMessage
	hugeValue  json.RawMessage
)

func init() {
	setupData(10, &smallValue)
	setupData(1024*1024, &largeValue)
	setupData(1024*1024*200, &hugeValue)
}

func setupData(size int, target *json.RawMessage) {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		panic(fmt.Sprintf("error creating random data: %v", err))
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("error marshalling data: %v", err))
	}
	*target = json.RawMessage(jsonData)
}

func TestStorageCreate(t *testing.T) {
	// create a teporary file for testing
	err := os.WriteFile("testdata/empty.json", []byte(`{}`), 0644)
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove("testdata/empty.json")

	// create a new storage
	_, err = New("testdata/empty.json")
	if err != nil {
		t.Fatalf("error creating storage: %v", err)
	}
}

func TestStorageStore(t *testing.T) {
	// create a teporary file for testing
	err := os.WriteFile("testdata/empty.json", []byte(`{}`), 0644)
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove("testdata/empty.json")

	// create a new storage
	s, err := New("testdata/empty.json")
	if err != nil {
		t.Fatalf("error creating storage: %v", err)
	}

	// store a value
	err = s.Store("test", []byte(`"value"`))
	if err != nil {
		t.Fatalf("error storing value: %v", err)
	}
}

func TestStorageGet(t *testing.T) {
	// create a teporary file for testing
	err := os.WriteFile("testdata/empty.json", []byte(`{}`), 0644)
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove("testdata/empty.json")

	// create a new storage
	s, err := New("testdata/empty.json")
	if err != nil {
		t.Fatalf("error creating storage: %v", err)
	}

	// store a value
	err = s.Store("test", []byte(`"value"`))
	if err != nil {
		t.Fatalf("error storing value: %v", err)
	}

	// get the value
	v, err := s.Get("test")
	if err != nil {
		t.Fatalf("error getting value: %v", err)
	}

	if string(v) != `"value"` {
		t.Fatalf("expected value to be 'value', got %s", v)
	}
}

func TestStorageFileWithError(t *testing.T) {
	// create a new storage
	_, err := New("testdata/error.json")
	if err == nil {
		t.Fatalf("expected error creating storage")
	}
}

func TestStorageNonExistingFile(t *testing.T) {
	// create a new storage
	_, err := New("testdata/non-existing.json")
	if err == nil {
		t.Fatalf("expected error creating storage")
	}
}

func TestStorageInsertMultipleValues(t *testing.T) {
	// create a teporary file for testing
	err := os.WriteFile("testdata/empty.json", []byte(`{}`), 0644)
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove("testdata/empty.json")

	// create a new storage
	s, err := New("testdata/empty.json")
	if err != nil {
		t.Fatalf("error creating storage: %v", err)
	}

	// store multiple values
	err = s.Store("test1", []byte(`"value1"`))
	if err != nil {
		t.Fatalf("error storing value: %v", err)
	}

	err = s.Store("test2", []byte(`"value2"`))
	if err != nil {
		t.Fatalf("error storing value: %v", err)
	}

	// get the values
	v1, err := s.Get("test1")
	if err != nil {
		t.Fatalf("error getting value: %v", err)
	}

	if string(v1) != `"value1"` {
		t.Fatalf("expected value to be 'value1', got %s", v1)
	}

	v2, err := s.Get("test2")
	if err != nil {
		t.Fatalf("error getting value: %v", err)
	}

	if string(v2) != `"value2"` {
		t.Fatalf("expected value to be 'value2', got %s", v2)
	}
}

func BenchmarkInsertMultipleSmallValues(b *testing.B) {
	// create a teporary file for testing
	err := os.WriteFile("testdata/empty.json", []byte(`{}`), 0644)
	if err != nil {
		b.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove("testdata/empty.json")

	// create a new storage
	s, err := New("testdata/empty.json")
	if err != nil {
		b.Fatalf("error creating storage: %v", err)
	}

	// store multiple values
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("test-%d", i)
		err = s.Store(key, smallValue)
		if err != nil {
			b.Fatalf("error storing value: %v", err)
		}
	}
}

func BenchmarkInsertMultipleLargeValues(b *testing.B) {
	// create a teporary file for testing
	err := os.WriteFile("testdata/empty.json", []byte(`{}`), 0644)
	if err != nil {
		b.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove("testdata/empty.json")

	// create a new storage
	s, err := New("testdata/empty.json")
	if err != nil {
		b.Fatalf("error creating storage: %v", err)
	}

	// store multiple values
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("test-%d", i)
		err = s.Store(key, largeValue)
		if err != nil {
			b.Fatalf("error storing value: %v", err)
		}
	}
}

func BenchmarkInsertMultipleHugeValues(b *testing.B) {
	// create a teporary file for testing
	err := os.WriteFile("testdata/empty.json", []byte(`{}`), 0644)
	if err != nil {
		b.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove("testdata/empty.json")

	// create a new storage
	s, err := New("testdata/empty.json")
	if err != nil {
		b.Fatalf("error creating storage: %v", err)
	}

	// store multiple values
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("test-%d", i)
		err = s.Store(key, hugeValue)
		if err != nil {
			b.Fatalf("error storing value: %v", err)
		}
	}
}
