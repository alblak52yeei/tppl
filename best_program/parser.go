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
	
	// Parse timestamp (8 bytes, little-endian, signed int64)
	// Read as uint64 first, then convert to int64 to handle sign correctly
	timestampRaw := binary.LittleEndian.Uint64(data[0:8])
	timestampMicro := int64(timestampRaw)
	// Convert microseconds to seconds and nanoseconds
	// Handle negative values correctly
	var seconds int64
	var nanoseconds int64
	if timestampMicro >= 0 {
		seconds = timestampMicro / TimestampMultiplier
		nanoseconds = (timestampMicro % TimestampMultiplier) * 1000
	} else {
		// For negative timestamps (shouldn't happen for POSIX, but handle it)
		seconds = timestampMicro / TimestampMultiplier
		remainder := timestampMicro % TimestampMultiplier
		if remainder < 0 {
			remainder += TimestampMultiplier
			seconds--
		}
		nanoseconds = remainder * 1000
	}
	timestamp := time.Unix(seconds, nanoseconds)
	
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
	
	// Parse timestamp (8 bytes, little-endian, signed int64)
	// Read as uint64 first, then convert to int64 to handle sign correctly
	timestampRaw := binary.LittleEndian.Uint64(data[0:8])
	timestampMicro := int64(timestampRaw)
	// Convert microseconds to seconds and nanoseconds
	// Handle negative values correctly
	var seconds int64
	var nanoseconds int64
	if timestampMicro >= 0 {
		seconds = timestampMicro / TimestampMultiplier
		nanoseconds = (timestampMicro % TimestampMultiplier) * 1000
	} else {
		// For negative timestamps (shouldn't happen for POSIX, but handle it)
		seconds = timestampMicro / TimestampMultiplier
		remainder := timestampMicro % TimestampMultiplier
		if remainder < 0 {
			remainder += TimestampMultiplier
			seconds--
		}
		nanoseconds = remainder * 1000
	}
	timestamp := time.Unix(seconds, nanoseconds)
	
	// Parse X, Y, Z (each 4 bytes, little-endian, signed integers)
	// Read as uint32 first, then convert to int32 to handle sign correctly
	xRaw := binary.LittleEndian.Uint32(data[8:12])
	yRaw := binary.LittleEndian.Uint32(data[12:16])
	zRaw := binary.LittleEndian.Uint32(data[16:20])
	x := int32(xRaw)
	y := int32(yRaw)
	z := int32(zRaw)
	
	return &DataRecord{
		Timestamp: timestamp,
		Source:    "server2",
		X:         x,
		Y:         y,
		Z:         z,
	}, nil
}

