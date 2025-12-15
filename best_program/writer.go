package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// Writer handles writing data records to a file
type Writer struct {
	file     *os.File
	buffer   *bufio.Writer
	mu       sync.Mutex
	recordCount int64
}

// NewWriter creates a new writer for the output file
func NewWriter(filename string) (*Writer, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	
	return &Writer{
		file:   file,
		buffer: bufio.NewWriter(file),
	}, nil
}

// WriteRecord writes a data record to the file
func (w *Writer) WriteRecord(record *DataRecord) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	// Format timestamp as YYYY-MM-DD HH:MM:SS
	timestampStr := record.Timestamp.Format("2006-01-02 15:04:05")
	
	var line string
	if record.Source == "server1" {
		line = fmt.Sprintf("%s,%s,%.6f,%d\n",
			timestampStr,
			record.Source,
			record.Temperature,
			record.Pressure,
		)
	} else {
		line = fmt.Sprintf("%s,%s,%d,%d,%d\n",
			timestampStr,
			record.Source,
			record.X,
			record.Y,
			record.Z,
		)
	}
	
	if _, err := w.buffer.WriteString(line); err != nil {
		return fmt.Errorf("failed to write to buffer: %w", err)
	}
	
	w.recordCount++
	
	// Flush periodically for reliability (every 100 records)
	if w.recordCount%100 == 0 {
		if err := w.buffer.Flush(); err != nil {
			return fmt.Errorf("failed to flush buffer: %w", err)
		}
		if err := w.file.Sync(); err != nil {
			return fmt.Errorf("failed to sync file: %w", err)
		}
	}
	
	return nil
}

// Flush flushes the buffer and syncs the file
func (w *Writer) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	if err := w.buffer.Flush(); err != nil {
		return err
	}
	return w.file.Sync()
}

// Close closes the writer and file
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	
	if err := w.buffer.Flush(); err != nil {
		return err
	}
	if err := w.file.Sync(); err != nil {
		return err
	}
	return w.file.Close()
}

// GetRecordCount returns the number of records written
func (w *Writer) GetRecordCount() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.recordCount
}

