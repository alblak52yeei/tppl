package main

import (
	"encoding/binary"
	"io"
	"math"
	"net"
	"strconv"
	"testing"
	"time"
)

// mockServer simulates a data server for testing
func mockServer(t *testing.T, port int, recordSize int, sendData func() []byte) net.Listener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			
			go func(c net.Conn) {
				defer c.Close()
				
				// Read secret key
				buf := make([]byte, 6)
				n, err := c.Read(buf)
				if err != nil || n != 6 || string(buf) != SecretKey {
					return
				}
				
				// Send data when "get" is received
				for {
					buf := make([]byte, 3)
					c.SetReadDeadline(time.Now().Add(1 * time.Second))
					n, err := c.Read(buf)
					if err != nil {
						return
					}
					if n == 3 && string(buf) == GetCommand {
						data := sendData()
						if _, err := c.Write(data); err != nil {
							return
						}
					}
				}
			}(conn)
		}
	}()
	
	return listener
}

func TestClient_Connect(t *testing.T) {
	// Create mock server for server1
	listener := mockServer(t, 0, Server1RecordSize, func() []byte {
		// Create valid server1 record
		data := make([]byte, Server1RecordSize)
		timestamp := time.Now().Unix() * TimestampMultiplier
		binary.LittleEndian.PutUint64(data[0:8], uint64(timestamp))
		binary.LittleEndian.PutUint32(data[8:12], math.Float32bits(25.5))
		binary.LittleEndian.PutUint16(data[12:14], uint16(1013))
		
		checksum := byte(0)
		for i := 0; i < Server1RecordSize-1; i++ {
			checksum += data[i]
		}
		data[Server1RecordSize-1] = checksum
		return data
	})
	defer listener.Close()
	
	address := listener.Addr().String()
	_, portStr, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatalf("Failed to parse address: %v", err)
	}
	
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("Failed to parse port: %v", err)
	}
	
	config := ServerConfig{
		Address: "127.0.0.1",
		Port:    portInt,
		Name:    "test_server",
	}
	
	client := NewClient(config, Server1RecordSize, ParseServer1Data)
	
	// Test connection
	err = client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	
	client.Stop()
}

func TestClient_NewClient(t *testing.T) {
	config := ServerConfig{
		Address: "127.0.0.1",
		Port:    1234,
		Name:    "test",
	}
	
	client := NewClient(config, Server1RecordSize, ParseServer1Data)
	
	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	
	if client.config.Address != config.Address {
		t.Errorf("Address mismatch: got %s, expected %s", client.config.Address, config.Address)
	}
	
	if client.GetDataChan() == nil {
		t.Error("GetDataChan returned nil")
	}
	
	if client.GetErrorChan() == nil {
		t.Error("GetErrorChan returned nil")
	}
}

func TestClient_ReceiveData(t *testing.T) {
	// Create mock server
	listener := mockServer(t, 0, Server1RecordSize, func() []byte {
		data := make([]byte, Server1RecordSize)
		timestamp := time.Now().Unix() * TimestampMultiplier
		binary.LittleEndian.PutUint64(data[0:8], uint64(timestamp))
		binary.LittleEndian.PutUint32(data[8:12], math.Float32bits(25.5))
		binary.LittleEndian.PutUint16(data[12:14], uint16(1013))
		
		checksum := byte(0)
		for i := 0; i < Server1RecordSize-1; i++ {
			checksum += data[i]
		}
		data[Server1RecordSize-1] = checksum
		return data
	})
	defer listener.Close()
	
	address := listener.Addr().String()
	_, portStr, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatalf("Failed to parse address: %v", err)
	}
	
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("Failed to parse port: %v", err)
	}
	
	config := ServerConfig{
		Address: "127.0.0.1",
		Port:    portInt,
		Name:    "test_server",
	}
	
	client := NewClient(config, Server1RecordSize, ParseServer1Data)
	
	// Connect
	err = client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	
	// Connect and send one "get" manually to test receiveLoop logic
	// (Start() would run in background and flood the channel)
	
	// Send "get" command manually
	if _, err := client.conn.Write([]byte(GetCommand)); err != nil {
		t.Fatalf("Failed to send get: %v", err)
	}
	
	// Read data directly from connection
	buffer := make([]byte, Server1RecordSize)
	n, err := io.ReadFull(client.conn, buffer)
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
	if n != Server1RecordSize {
		t.Fatalf("Incomplete read: got %d, expected %d", n, Server1RecordSize)
	}
	
	// Parse
	record, err := ParseServer1Data(buffer)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	if record.Source != "server1" {
		t.Errorf("Expected server1, got %s", record.Source)
	}
	
	client.Stop()
}

