package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/saurabhSPatel/p2p-file-sharing/internal/cli"
	"github.com/saurabhSPatel/p2p-file-sharing/internal/discovery"
	"github.com/saurabhSPatel/p2p-file-sharing/internal/network"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	// Setup libp2p host
	host, err := network.SetupHost(ctx)
	if err != nil {
		log.Fatalf("Failed to setup host: %v", err)
	}
	defer func() {
		if err := host.Close(); err != nil {
			log.Printf("Error shutting down host: %v", err)
		}
	}()

	// Print the host's addresses
	log.Println("Host ID:", host.ID())
	log.Println("Host Addresses:")
	for _, addr := range host.Addrs() {
		log.Printf("  %s/p2p/%s\n", addr, host.ID())
	}

	// Setup discovery service
	disc := discovery.NewDiscovery(host)
	if err = disc.SetupDiscovery(); err != nil {
		log.Fatalf("Failed to setup discovery: %v", err)
	}

	// Define and Ensure the directories exist
	sharedDir := "./shared"
	downloadDir := "./downloads"

	if err := os.MkdirAll(sharedDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create shared directory: %v", err)
	}
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create download directory: %v", err)
	}

	// Setup CLI
	c := cli.NewCLI(host, disc, sharedDir, downloadDir, ctx)

	// Handle graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		log.Println("Received interrupt signal, shutting down...")
		cancel()
	}()

	// Run the CLI
	fmt.Println("Starting CLI...")
	c.Run()

	log.Println("Shutdown complete")
}
