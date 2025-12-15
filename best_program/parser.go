package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
)

const (
	// Server1 constants
	Server1RecordSize = 8 + 4 + 2 + 1 // timestamp + temp + pressure + checksum = 15 bytes
	
	// Server2 constants
	Server2RecordSize = 8 + 4 + 4 + 4 + 1 // timestamp + X + Y + Z + checksum = 21 bytes
	
	// Timestamp multiplier
	TimestampMultiplier = 1000000 // 10^6
)

// ParseServer1Data parses data from server 1 (temperature/pressure)
func ParseServer1Data(data []byte) (*DataRecord, error) {
	if len(data) < Server1RecordSize {
		return nil, errors.New("insufficient data for server1 record")
	}
	
	// Verify checksum
	checksum := byte(0)
	for i := 0; i < Server1RecordSize-1; i++ {
		checksum += data[i]
	}
	if checksum != data[Server1RecordSize-1] {
		return nil, fmt.Errorf("checksum mismatch: expected %d, got %d", checksum, data[Server1RecordSize-1])
	}
	
	// Parse timestamp (8 bytes, little-endian)
	timestampMicro := int64(binary.LittleEndian.Uint64(data[0:8]))
	timestamp := time.Unix(timestampMicro/TimestampMultiplier, (timestampMicro%TimestampMultiplier)*1000)
	
	// Parse temperature (4 bytes, little-endian, float32)
	tempFloat := math.Float32frombits(binary.LittleEndian.Uint32(data[8:12]))
	
	// Parse pressure (2 bytes, little-endian, signed integer)
	pressure := int16(binary.LittleEndian.Uint16(data[12:14]))
	
	return &DataRecord{
		Timestamp:   timestamp,
		Source:      "server1",
		Temperature: tempFloat,
		Pressure:    pressure,
	}, nil
}

// ParseServer2Data parses data from server 2 (X, Y, Z)
func ParseServer2Data(data []byte) (*DataRecord, error) {
	if len(data) < Server2RecordSize {
		return nil, errors.New("insufficient data for server2 record")
	}
	
	// Verify checksum
	checksum := byte(0)
	for i := 0; i < Server2RecordSize-1; i++ {
		checksum += data[i]
	}
	if checksum != data[Server2RecordSize-1] {
		return nil, fmt.Errorf("checksum mismatch: expected %d, got %d", checksum, data[Server2RecordSize-1])
	}
	
	// Parse timestamp (8 bytes, little-endian)
	timestampMicro := int64(binary.LittleEndian.Uint64(data[0:8]))
	timestamp := time.Unix(timestampMicro/TimestampMultiplier, (timestampMicro%TimestampMultiplier)*1000)
	
	// Parse X, Y, Z (each 4 bytes, little-endian, signed integers)
	x := int32(binary.LittleEndian.Uint32(data[8:12]))
	y := int32(binary.LittleEndian.Uint32(data[12:16]))
	z := int32(binary.LittleEndian.Uint32(data[16:20]))
	
	return &DataRecord{
		Timestamp: timestamp,
		Source:    "server2",
		X:         x,
		Y:         y,
		Z:         z,
	}, nil
}

