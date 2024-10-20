package network

import (
	"context"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"testing"
	"time"
)

// Test SetupHost function
func TestSetupHost(t *testing.T) {
	ctx := context.Background()

	host, err := SetupHost(ctx)
	if err != nil {
		t.Fatalf("Failed to setup host: %v", err)
	}
	defer host.Close()

	// Check if the host ID is set correctly
	if host.ID().String() == "" {
		t.Fatalf("Host ID is not set")
	}

	// Check if the host is listening on at least one address
	if len(host.Addrs()) == 0 {
		t.Fatalf("Host is not listening on any addresses")
	}
}

// Test SendFile function
func TestSendAndReceiveFile(t *testing.T) {
	ctx := context.Background()

	// Create two hosts
	host1, err := SetupHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host1: %v", err)
	}
	defer host1.Close()

	host2, err := SetupHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host2: %v", err)
	}
	defer host2.Close()

	// Connect host1 to host2
	host1PeerInfo := peer.AddrInfo{
		ID:    host1.ID(),
		Addrs: host1.Addrs(),
	}

	if err := host2.Connect(ctx, host1PeerInfo); err != nil {
		t.Fatalf("Failed to connect host2 to host1: %v", err)
	}

	// Prepare file data to send
	filename := "testfile.txt"
	fileContent := []byte("This is a test file content.")

	go func() {
		time.Sleep(time.Second) // Give host2 some time to set up the stream handler
		err := SendFile(ctx, host1, host2.ID(), filename, fileContent)
		if err != nil {
			t.Errorf("Failed to send file: %v", err)
		}
	}()

	// Set stream handler on host2 to receive the file
	host2.SetStreamHandler(ProtocolID, func(stream network.Stream) {
		receivedFilename, receivedData, err := ReceiveFile(stream)
		if err != nil {
			t.Fatalf("Error receiving file: %v", err)
		}

		if receivedFilename != filename {
			t.Errorf("Expected filename %s, got %s", filename, receivedFilename)
		}

		if string(receivedData) != string(fileContent) {
			t.Errorf("Expected file content %s, got %s", string(fileContent), string(receivedData))
		}
	})

	time.Sleep(2 * time.Second) // Allow time for the transfer and test completion
}

func TestHandleStream(t *testing.T) {
	ctx := context.Background()

	// Set up two hosts (one to send, one to receive)
	host1, err := SetupHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host1: %v", err)
	}
	defer host1.Close()

	host2, err := SetupHost(ctx)
	if err != nil {
		t.Fatalf("Failed to create host2: %v", err)
	}
	defer host2.Close()

	// Connect host1 to host2
	host1PeerInfo := peer.AddrInfo{
		ID:    host1.ID(),
		Addrs: host1.Addrs(),
	}

	if err := host2.Connect(ctx, host1PeerInfo); err != nil {
		t.Fatalf("Failed to connect host2 to host1: %v", err)
	}

	// Set up stream handler for host2 to use HandleStream for incoming files
	host2.SetStreamHandler(ProtocolID, func(stream network.Stream) {
		HandleStream(stream)
	})

	// Prepare file data to send
	filename := "testfile.txt"
	fileContent := []byte("This is a test file content.")

	// Use a goroutine to simulate sending a file from host1 to host2
	go func() {
		time.Sleep(time.Second) // Give host2 time to set up the stream handler
		err := SendFile(ctx, host1, host2.ID(), filename, fileContent)
		if err != nil {
			t.Errorf("Failed to send file: %v", err)
		}
	}()

	// Wait for the transfer to complete
	time.Sleep(3 * time.Second) // Give enough time for the file transfer and stream handling
}
