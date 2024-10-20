package cli

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/saurabhSPatel/p2p-file-sharing/internal/discovery"
	"github.com/saurabhSPatel/p2p-file-sharing/internal/file"
	"github.com/saurabhSPatel/p2p-file-sharing/internal/network"
)

// CLI represents the command-line interface for file sharing.
type CLI struct {
	host        host.Host
	discovery   *discovery.Discovery
	sharedDir   string
	downloadDir string
	ctx         context.Context
}

// NewCLI initializes a new CLI instance.
func NewCLI(h host.Host, d *discovery.Discovery, sharedDir string, downloadDir string, ctx context.Context) *CLI {
	return &CLI{host: h, discovery: d, sharedDir: sharedDir, downloadDir: downloadDir, ctx: ctx}
}

// Run starts the CLI to listen for user commands.
func (c *CLI) Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the P2P File Sharing CLI!")
	fmt.Println("Available commands: list, download, upload, quit")
	for {
		select {
		case <-c.ctx.Done():
			log.Println("Shutting down CLI...")
			return
		default:
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			parts := strings.Fields(input)

			if len(parts) == 0 {
				continue
			}

			switch parts[0] {
			case "list":
				c.listFiles()
			case "download":
				if len(parts) < 2 {
					log.Println("Usage: download <filename>")
					continue
				}
				c.downloadFile(parts[1])
			case "upload":
				if len(parts) < 2 {
					log.Println("Usage: upload <filename>")
					continue
				}
				c.uploadFile(parts[1])
			case "exit":
				return
			default:
				log.Println("Unknown command. Available commands: list, download, upload, exit")
			}
		}
	}
}

// listFiles displays all files available in the shared directory.
func (c *CLI) listFiles() {
	files, err := file.ListFiles(c.sharedDir)
	if err != nil {
		log.Printf("Error listing files in directory '%s': %v\n", c.sharedDir, err)
		return
	}
	if len(files) == 0 {
		log.Println("No files available.")
		return
	}
	log.Println("Available files:")
	for _, f := range files {
		log.Println(f)
	}
}

// downloadFile retrieves a file from a peer and saves it to the download directory.
func (c *CLI) downloadFile(filename string) {
	ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
	defer cancel()

	peerChan, err := c.discovery.DiscoverPeers(ctx)
	if err != nil {
		log.Printf("Error discovering peers: %v\n", err)
		return
	}

	for peer := range peerChan {
		stream, err := c.host.NewStream(context.Background(), peer.ID)
		if err != nil {
			log.Printf("Error creating stream with peer %s: %v\n", peer.ID, err)
			continue
		}

		if _, err := stream.Write([]byte(filename + "\n")); err != nil {
			log.Printf("Error requesting file: %v\n", err)
			continue
		}

		receivedFilename, data, err := network.ReceiveFile(stream)
		if err != nil {
			log.Printf("Error receiving file: %v\n", err)
			continue
		}

		if err := stream.Close(); err != nil {
			log.Printf("Error closing stream: %v\n", err)
		}

		if receivedFilename != filename {
			log.Printf("Received unexpected file: %s\n", receivedFilename)
			continue
		}

		savePath := c.downloadDir + "/" + filename
		if err := file.WriteFile(savePath, data); err != nil {
			log.Printf("Error saving file to '%s': %v\n", savePath, err)
			continue
		}

		log.Printf("File %s downloaded successfully to %s\n", filename, savePath)
		return
	}

	log.Println("File not found on any peer")
}

// uploadFile sends a file to a discovered peer.
func (c *CLI) uploadFile(filename string) {
	filePath := c.sharedDir + "/" + filename
	data, err := file.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file '%s': %v\n", filePath, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	peerChan, err := c.discovery.DiscoverPeers(ctx)
	if err != nil {
		log.Printf("Error discovering peers: %v\n", err)
		return
	}

	for peer := range peerChan {
		if err := network.SendFile(c.ctx, c.host, peer.ID, filename, data); err != nil {
			log.Printf("Error sending file to peer %s: %v\n", peer.ID, err)
			continue
		}
		log.Printf("File %s uploaded successfully to peer %s\n", filename, peer.ID)
		return
	}

	log.Println("No peers available to upload the file")
}
