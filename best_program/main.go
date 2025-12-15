package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	OutputFile = "data.log"
	Server1Address = "95.163.237.76"
	Server1Port = 5123
	Server2Address = "95.163.237.76"
	Server2Port = 5124
)

func main() {
	// Create writer
	writer, err := NewWriter(OutputFile)
	if err != nil {
		log.Fatalf("Failed to create writer: %v", err)
	}
	defer writer.Close()
	
	// Create clients
	server1Config := ServerConfig{
		Address: Server1Address,
		Port:    Server1Port,
		Name:    "server1",
	}
	
	server2Config := ServerConfig{
		Address: Server2Address,
		Port:    Server2Port,
		Name:    "server2",
	}
	
	client1 := NewClient(server1Config, Server1RecordSize, ParseServer1Data)
	client2 := NewClient(server2Config, Server2RecordSize, ParseServer2Data)
	
	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Start clients
	client1.Start()
	client2.Start()
	
	log.Println("Data collection started. Press Ctrl+C to stop.")
	
	// Statistics
	var wg sync.WaitGroup
	startTime := time.Now()
	
	// Start data writer goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		writeData(client1, client2, writer)
	}()
	
	// Start error handler goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleErrors(client1, client2)
	}()
	
	// Start statistics reporter
	wg.Add(1)
	go func() {
		defer wg.Done()
		reportStatistics(writer, startTime, sigChan)
	}()
	
	// Wait for interrupt signal
	<-sigChan
	log.Println("\nShutting down...")
	
	// Stop clients
	client1.Stop()
	client2.Stop()
	
	// Wait for all goroutines to finish
	wg.Wait()
	
	// Final flush
	if err := writer.Flush(); err != nil {
		log.Printf("Error flushing writer: %v", err)
	}
	
	recordCount := writer.GetRecordCount()
	duration := time.Since(startTime)
	log.Printf("Total records written: %d", recordCount)
	log.Printf("Duration: %v", duration)
	log.Printf("Average rate: %.2f records/second", float64(recordCount)/duration.Seconds())
	log.Println("Shutdown complete.")
}

// writeData writes data from both clients to the file
func writeData(client1, client2 *Client, writer *Writer) {
	client1Closed := false
	client2Closed := false
	
	for {
		if client1Closed && client2Closed {
			return
		}
		
		select {
		case record, ok := <-client1.GetDataChan():
			if !ok {
				client1Closed = true
				continue
			}
			if err := writer.WriteRecord(record); err != nil {
				log.Printf("Error writing record from %s: %v", record.Source, err)
			}
		case record, ok := <-client2.GetDataChan():
			if !ok {
				client2Closed = true
				continue
			}
			if err := writer.WriteRecord(record); err != nil {
				log.Printf("Error writing record from %s: %v", record.Source, err)
			}
		}
	}
}

// handleErrors logs errors from clients
func handleErrors(client1, client2 *Client) {
	client1Closed := false
	client2Closed := false
	
	for {
		if client1Closed && client2Closed {
			return
		}
		
		select {
		case err, ok := <-client1.GetErrorChan():
			if !ok {
				client1Closed = true
				continue
			}
			log.Printf("Client1 error: %v", err)
		case err, ok := <-client2.GetErrorChan():
			if !ok {
				client2Closed = true
				continue
			}
			log.Printf("Client2 error: %v", err)
		}
	}
}

// reportStatistics periodically reports statistics
func reportStatistics(writer *Writer, startTime time.Time, sigChan chan os.Signal) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			recordCount := writer.GetRecordCount()
			duration := time.Since(startTime)
			rate := float64(recordCount) / duration.Seconds()
			log.Printf("Statistics: %d records, %.2f records/sec, running for %v",
				recordCount, rate, duration.Round(time.Second))
		case <-sigChan:
			return
		}
	}
}

