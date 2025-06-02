package db

import (
	"os"
	"testing"
)

func TestGetSet(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "testdb.db")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	fileName := f.Name()
	defer os.Remove(fileName)
	defer f.Close()

	db, close, err := NewDatabase(fileName)
	if err != nil {
		t.Fatalf("Could not create database: %v", err)
	}
	defer close()

	if err := db.SetKey("test", []byte("value")); err != nil {
		t.Fatalf("Could not set key: %v", err)
	}

	value, err := db.GetKey("test")
	if err != nil {
		t.Fatalf("Could not get key: %v", err)
	}

	if string(value) != "value" {
		t.Errorf("Expected value 'value', got '%s'", value)
	}
}
