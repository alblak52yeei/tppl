package main

import "time"

// DataRecord represents a parsed data record from either server
type DataRecord struct {
	Timestamp   time.Time
	Source      string // "server1" or "server2"
	Temperature float32 // Only for server1
	Pressure    int16   // Only for server1
	X           int32   // Only for server2
	Y           int32   // Only for server2
	Z           int32   // Only for server2
}

// ServerConfig holds configuration for a server
type ServerConfig struct {
	Address string
	Port    int
	Name    string
}



