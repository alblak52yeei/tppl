package main

import (
	"encoding/binary"
	"math"
	"testing"
	"time"
)

func TestParseServer1Data(t *testing.T) {
	// Create test data
	timestamp := time.Now().Unix() * TimestampMultiplier
	temperature := float32(25.5)
	pressure := int16(1013)
	
	data := make([]byte, Server1RecordSize)
	
	// Write timestamp (8 bytes, little-endian)
	binary.BigEndian.PutUint64(data[0:8], uint64(timestamp))
	
	// Write temperature (4 bytes, big-endian, float32)
	binary.BigEndian.PutUint32(data[8:12], math.Float32bits(temperature))
	
	// Write pressure (2 bytes, big-endian, signed integer)
	binary.BigEndian.PutUint16(data[12:14], uint16(pressure))
	
	// Calculate checksum
	checksum := byte(0)
	for i := 0; i < Server1RecordSize-1; i++ {
		checksum += data[i]
	}
	data[Server1RecordSize-1] = checksum
	
	// Parse
	record, err := ParseServer1Data(data)
	if err != nil {
		t.Fatalf("Failed to parse server1 data: %v", err)
	}
	
	// Verify
	expectedTime := time.Unix(timestamp/TimestampMultiplier, (timestamp%TimestampMultiplier)*1000)
	if record.Timestamp.Unix() != expectedTime.Unix() {
		t.Errorf("Timestamp mismatch: got %v, expected %v", record.Timestamp, expectedTime)
	}
	
	if math.Abs(float64(record.Temperature-temperature)) > 0.001 {
		t.Errorf("Temperature mismatch: got %f, expected %f", record.Temperature, temperature)
	}
	
	if record.Pressure != pressure {
		t.Errorf("Pressure mismatch: got %d, expected %d", record.Pressure, pressure)
	}
	
	if record.Source != "server1" {
		t.Errorf("Source mismatch: got %s, expected server1", record.Source)
	}
}

func TestParseServer1Data_ChecksumError(t *testing.T) {
	data := make([]byte, Server1RecordSize)
	data[Server1RecordSize-1] = 0xFF // Wrong checksum
	
	_, err := ParseServer1Data(data)
	if err == nil {
		t.Error("Expected checksum error, got nil")
	}
}

func TestParseServer2Data(t *testing.T) {
	// Create test data
	timestamp := time.Now().Unix() * TimestampMultiplier
	x := int32(100)
	y := int32(-200) // Test negative value
	z := int32(300)
	
	data := make([]byte, Server2RecordSize)
	
	// Write timestamp (8 bytes, big-endian)
	binary.BigEndian.PutUint64(data[0:8], uint64(timestamp))
	
	// Write X, Y, Z (each 4 bytes, big-endian, signed integers)
	binary.BigEndian.PutUint32(data[8:12], uint32(x))
	binary.BigEndian.PutUint32(data[12:16], uint32(y))
	binary.BigEndian.PutUint32(data[16:20], uint32(z))
	
	// Calculate checksum
	checksum := byte(0)
	for i := 0; i < Server2RecordSize-1; i++ {
		checksum += data[i]
	}
	data[Server2RecordSize-1] = checksum
	
	// Parse
	record, err := ParseServer2Data(data)
	if err != nil {
		t.Fatalf("Failed to parse server2 data: %v", err)
	}
	
	// Verify
	expectedTime := time.Unix(timestamp/TimestampMultiplier, (timestamp%TimestampMultiplier)*1000)
	if record.Timestamp.Unix() != expectedTime.Unix() {
		t.Errorf("Timestamp mismatch: got %v, expected %v", record.Timestamp, expectedTime)
	}
	
	if record.X != x {
		t.Errorf("X mismatch: got %d, expected %d", record.X, x)
	}
	
	if record.Y != y {
		t.Errorf("Y mismatch: got %d, expected %d", record.Y, y)
	}
	
	if record.Z != z {
		t.Errorf("Z mismatch: got %d, expected %d", record.Z, z)
	}
	
	if record.Source != "server2" {
		t.Errorf("Source mismatch: got %s, expected server2", record.Source)
	}
}

func TestParseServer2Data_ChecksumError(t *testing.T) {
	data := make([]byte, Server2RecordSize)
	data[Server2RecordSize-1] = 0xFF // Wrong checksum
	
	_, err := ParseServer2Data(data)
	if err == nil {
		t.Error("Expected checksum error, got nil")
	}
}

func TestParseServer1Data_InsufficientData(t *testing.T) {
	data := make([]byte, Server1RecordSize-1)
	
	_, err := ParseServer1Data(data)
	if err == nil {
		t.Error("Expected error for insufficient data, got nil")
	}
}

func TestParseServer2Data_InsufficientData(t *testing.T) {
	data := make([]byte, Server2RecordSize-1)
	
	_, err := ParseServer2Data(data)
	if err == nil {
		t.Error("Expected error for insufficient data, got nil")
	}
}

