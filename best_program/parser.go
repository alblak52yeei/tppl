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
	
	// Parse timestamp (8 bytes, big-endian/network byte order, signed int64)
	timestampRaw := binary.BigEndian.Uint64(data[0:8])
	timestampMicro := int64(timestampRaw)
	// Convert microseconds to seconds and nanoseconds
	seconds := timestampMicro / TimestampMultiplier
	nanoseconds := (timestampMicro % TimestampMultiplier) * 1000
	timestamp := time.Unix(seconds, nanoseconds)
	
	// Parse temperature (4 bytes, big-endian, float32)
	tempFloat := math.Float32frombits(binary.BigEndian.Uint32(data[8:12]))
	
	// Parse pressure (2 bytes, big-endian, signed integer)
	pressure := int16(binary.BigEndian.Uint16(data[12:14]))
	
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
	
	// Parse timestamp (8 bytes, big-endian/network byte order, signed int64)
	timestampRaw := binary.BigEndian.Uint64(data[0:8])
	timestampMicro := int64(timestampRaw)
	// Convert microseconds to seconds and nanoseconds
	seconds := timestampMicro / TimestampMultiplier
	nanoseconds := (timestampMicro % TimestampMultiplier) * 1000
	timestamp := time.Unix(seconds, nanoseconds)
	
	// Parse X, Y, Z (each 4 bytes, big-endian, signed integers)
	xRaw := binary.BigEndian.Uint32(data[8:12])
	yRaw := binary.BigEndian.Uint32(data[12:16])
	zRaw := binary.BigEndian.Uint32(data[16:20])
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

