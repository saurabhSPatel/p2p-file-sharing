package test_test

import (
	"context"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/saurabhSPatel/p2p-file-sharing/internal/discovery"
	"github.com/saurabhSPatel/p2p-file-sharing/internal/network"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeerToPeerFileSharing(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Create a temporary shared directory for the test
	sharedDir := os.TempDir()

	// Setup hosts and discovery
	hosts := setupHosts(t, ctx, 3)
	defer func() {
		for _, p2pHost := range hosts {
			err := p2pHost.Close()
			if err != nil {
				t.Errorf("Error closing host: %v", err)
				return
			}
		}
	}()

	setupDiscoveries(t, hosts)

	// Wait for peer discovery
	time.Sleep(5 * time.Second)

	// Prepare test files
	testFiles := prepareTestFiles(t, sharedDir)

	// Run test scenarios
	t.Run("Basic File Transfer", func(t *testing.T) {
		testFileTransfer(t, ctx, hosts[0], hosts[1], testFiles["basic"], sharedDir)
	})

	t.Run("Large File Transfer", func(t *testing.T) {
		testFileTransfer(t, ctx, hosts[1], hosts[2], testFiles["large"], sharedDir)
	})

	t.Run("Multiple File Transfers", func(t *testing.T) {
		testMultipleFileTransfers(t, ctx, hosts, sharedDir)
	})

	t.Run("Concurrent Transfers", func(t *testing.T) {
		testConcurrentTransfers(t, ctx, hosts, sharedDir)
	})
}

func setupHosts(t *testing.T, ctx context.Context, count int) []host.Host {
	hosts := make([]host.Host, count)
	for i := 0; i < count; i++ {
		p2pHost, err := network.SetupHost(ctx)
		require.NoError(t, err, "Failed to setup p2pHost%d", i+1)
		hosts[i] = p2pHost
		t.Logf("Host%d ID: %s", i+1, p2pHost.ID().String())
	}
	return hosts
}

func setupDiscoveries(t *testing.T, hosts []host.Host) {
	for i, p2pHost := range hosts {
		disc := discovery.NewDiscovery(p2pHost)
		err := disc.SetupDiscovery()
		require.NoError(t, err, "Failed to setup discovery for p2pHost%d", i+1)
	}
}

func prepareTestFiles(t *testing.T, sharedDir string) map[string]testFile {
	// Prepare test files inside the temporary shared directory
	files := map[string]testFile{
		"basic": {
			name:    "basic.txt",
			content: []byte("This is a basic test file for integration testing."),
		},
		"large": {
			name:    "large.bin",
			content: make([]byte, 10*1024*1024), // 10MB file
		},
	}

	for _, file := range files {
		filePath := filepath.Join(sharedDir, file.name)
		require.NoError(t, os.WriteFile(filePath, file.content, 0644), "Failed to write test file: %s", file.name)
		t.Logf("Test file '%s' created at '%s'", file.name, filePath)
	}

	return files
}

func testFileTransfer(t *testing.T, ctx context.Context, sender, receiver host.Host, file testFile, sharedDir string) {
	// Sender sends the file
	err := network.SendFile(ctx, sender, receiver.ID(), filepath.Join(sharedDir, file.name), file.content)
	require.NoError(t, err, "Failed to send file")

	// Wait for the file to be processed
	time.Sleep(1 * time.Second)

	// Verify the received file
	receivedPath := filepath.Join(sharedDir, file.name)
	receivedContent, err := os.ReadFile(receivedPath)
	require.NoError(t, err, "Failed to read received file")
	assert.Equal(t, file.content, receivedContent, "File content mismatch")

	t.Logf("File '%s' transferred successfully", file.name)
}

func testMultipleFileTransfers(t *testing.T, ctx context.Context, hosts []host.Host, sharedDir string) {
	files := []testFile{
		{name: "file1.txt", content: []byte("Content of file1")},
		{name: "file2.txt", content: []byte("Content of file2")},
		{name: "file3.txt", content: []byte("Content of file3")},
	}

	for _, file := range files {
		// Write each file to the shared directory
		err := os.WriteFile(filepath.Join(sharedDir, file.name), file.content, 0644)
		require.NoError(t, err, "Failed to write test file: %s", file.name)

		// Test file transfer between hosts
		testFileTransfer(t, ctx, hosts[0], hosts[1], file, sharedDir)
	}
}

func testConcurrentTransfers(t *testing.T, ctx context.Context, hosts []host.Host, sharedDir string) {
	files := []testFile{
		{name: "concurrent1.txt", content: []byte("Content of concurrent1")},
		{name: "concurrent2.txt", content: []byte("Content of concurrent2")},
		{name: "concurrent3.txt", content: []byte("Content of concurrent3")},
	}

	errChan := make(chan error, len(files))
	for _, file := range files {
		// Write each file to the shared directory
		err := os.WriteFile(filepath.Join(sharedDir, file.name), file.content, 0644)
		require.NoError(t, err, "Failed to write test file: %s", file.name)

		// Start concurrent file transfers
		go func(f testFile) {
			err := network.SendFile(ctx, hosts[0], hosts[1].ID(), filepath.Join(sharedDir, f.name), f.content)
			errChan <- err
		}(file)
	}

	for range files {
		require.NoError(t, <-errChan, "Error in concurrent transfer")
	}

	// Wait for all files to be processed
	time.Sleep(2 * time.Second)

	// Verify all files were received correctly
	for _, file := range files {
		receivedPath := filepath.Join(sharedDir, file.name)
		receivedContent, err := os.ReadFile(receivedPath)
		require.NoError(t, err, "Failed to read received file")
		assert.Equal(t, file.content, receivedContent, "File content mismatch for concurrent transfer")
	}
}

type testFile struct {
	name    string
	content []byte
}
