package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	SecretKey    = "isu_pt"
	GetCommand   = "get"
	ReadTimeout  = 5 * time.Second
	ReconnectDelay = 2 * time.Second
)

// Client handles connection to a data server
type Client struct {
	config     ServerConfig
	conn       net.Conn
	recordSize int
	parseFunc  func([]byte) (*DataRecord, error)
	stopChan   chan struct{}
	dataChan   chan *DataRecord
	errorChan  chan error
}

// NewClient creates a new client for the given server configuration
func NewClient(config ServerConfig, recordSize int, parseFunc func([]byte) (*DataRecord, error)) *Client {
	return &Client{
		config:     config,
		recordSize: recordSize,
		parseFunc:  parseFunc,
		stopChan:   make(chan struct{}),
		dataChan:   make(chan *DataRecord, 1000), // Buffered channel for high throughput
		errorChan:  make(chan error, 10),
	}
}

// Connect establishes connection to the server and authenticates
func (c *Client) Connect() error {
	address := fmt.Sprintf("%s:%d", c.config.Address, c.config.Port)
	
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	
	c.conn = conn
	
	// Send secret key for authentication
	if _, err := conn.Write([]byte(SecretKey)); err != nil {
		conn.Close()
		return fmt.Errorf("failed to send secret key: %w", err)
	}
	
	log.Printf("[%s] Connected and authenticated", c.config.Name)
	return nil
}

// Start begins receiving data from the server
func (c *Client) Start() {
	go c.run()
}

// Stop stops the client
func (c *Client) Stop() {
	close(c.stopChan)
	if c.conn != nil {
		c.conn.Close()
	}
}

// GetDataChan returns the channel for receiving parsed data
func (c *Client) GetDataChan() <-chan *DataRecord {
	return c.dataChan
}

// GetErrorChan returns the channel for receiving errors
func (c *Client) GetErrorChan() <-chan error {
	return c.errorChan
}

// run is the main loop that handles connection, data requests, and reconnection
func (c *Client) run() {
	defer close(c.dataChan)
	defer close(c.errorChan)
	
	for {
		select {
		case <-c.stopChan:
			return
		default:
		}
		
		// Connect to server
		if err := c.Connect(); err != nil {
			c.errorChan <- fmt.Errorf("[%s] connection error: %w", c.config.Name, err)
			select {
			case <-c.stopChan:
				return
			case <-time.After(ReconnectDelay):
				continue
			}
		}
		
		// Start receiving data
		c.receiveLoop()
		
		// Close connection
		if c.conn != nil {
			c.conn.Close()
			c.conn = nil
		}
		
		// Wait before reconnecting
		select {
		case <-c.stopChan:
			return
		case <-time.After(ReconnectDelay):
		}
	}
}

// receiveLoop continuously requests and receives data
func (c *Client) receiveLoop() {
	buffer := make([]byte, c.recordSize)
	
	for {
		select {
		case <-c.stopChan:
			return
		default:
		}
		
		// Send "get" command first
		if _, err := c.conn.Write([]byte(GetCommand)); err != nil {
			c.errorChan <- fmt.Errorf("[%s] failed to send get command: %w", c.config.Name, err)
			return
		}
		
		// Set read timeout AFTER sending command
		c.conn.SetReadDeadline(time.Now().Add(ReadTimeout))
		
		// Read full record directly from connection (no buffering)
		n, err := io.ReadFull(c.conn, buffer)
		if err != nil {
			if err == io.EOF {
				c.errorChan <- fmt.Errorf("[%s] connection closed by server", c.config.Name)
				return
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout is expected if server doesn't respond, continue to send next "get"
				continue
			} else {
				c.errorChan <- fmt.Errorf("[%s] read error: %w", c.config.Name, err)
				return
			}
		}
		
		if n != c.recordSize {
			c.errorChan <- fmt.Errorf("[%s] incomplete record: got %d bytes, expected %d", c.config.Name, n, c.recordSize)
			continue
		}
		
		// Parse the record
		record, err := c.parseFunc(buffer)
		if err != nil {
			c.errorChan <- fmt.Errorf("[%s] parse error: %w", c.config.Name, err)
			continue
		}
		
		// Send parsed record to channel (non-blocking)
		select {
		case c.dataChan <- record:
		case <-c.stopChan:
			return
		default:
			// Channel full, log warning but continue
			log.Printf("[%s] Warning: data channel full, dropping record", c.config.Name)
		}
	}
}


