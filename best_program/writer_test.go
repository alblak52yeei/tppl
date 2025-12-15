package main

import (
	"os"
	"testing"
	"time"
)

func TestWriter(t *testing.T) {
	// Create temporary file
	tmpfile, err := os.CreateTemp("", "test_writer_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	
	// Create writer
	writer, err := NewWriter(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()
	
	// Write server1 record
	record1 := &DataRecord{
		Timestamp:   time.Now(),
		Source:      "server1",
		Temperature: 25.5,
		Pressure:    1013,
	}
	
	if err := writer.WriteRecord(record1); err != nil {
		t.Fatalf("Failed to write record1: %v", err)
	}
	
	// Write server2 record
	record2 := &DataRecord{
		Timestamp: time.Now(),
		Source:    "server2",
		X:         100,
		Y:         -200,
		Z:         300,
	}
	
	if err := writer.WriteRecord(record2); err != nil {
		t.Fatalf("Failed to write record2: %v", err)
	}
	
	// Flush
	if err := writer.Flush(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}
	
	// Verify record count
	count := writer.GetRecordCount()
	if count != 2 {
		t.Errorf("Expected 2 records, got %d", count)
	}
	
	// Read file and verify content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("File is empty")
	}
	
	// Verify format (basic check)
	contentStr := string(content)
	if len(contentStr) < 10 {
		t.Error("File content too short")
	}
}

func TestWriter_Concurrent(t *testing.T) {
	// Create temporary file
	tmpfile, err := os.CreateTemp("", "test_writer_concurrent_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	
	// Create writer
	writer, err := NewWriter(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()
	
	// Write records concurrently
	done := make(chan bool)
	numGoroutines := 10
	recordsPerGoroutine := 100
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < recordsPerGoroutine; j++ {
				record := &DataRecord{
					Timestamp: time.Now(),
					Source:    "server1",
					Temperature: 25.5,
					Pressure:    1013,
				}
				if err := writer.WriteRecord(record); err != nil {
					t.Errorf("Failed to write record: %v", err)
				}
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// Flush
	if err := writer.Flush(); err != nil {
		t.Fatalf("Failed to flush: %v", err)
	}
	
	// Verify record count
	expectedCount := int64(numGoroutines * recordsPerGoroutine)
	count := writer.GetRecordCount()
	if count != expectedCount {
		t.Errorf("Expected %d records, got %d", expectedCount, count)
	}
}

