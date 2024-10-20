package network

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

const ProtocolID = "/p2p-file-sharing/1.0.0"

// SetupHost initializes a libp2p host with basic configuration
func SetupHost(ctx context.Context) (host.Host, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.EnableRelay(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}

	pingService := ping.NewPingService(h)
	h.SetStreamHandler(ProtocolID, HandleStream)

	log.Println("Host created with ID:", h.ID().String())
	for _, addr := range h.Addrs() {
		log.Println("Listening on", addr)
	}

	h.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(net network.Network, conn network.Conn) {
			log.Printf("Connected to peer: %s\n", conn.RemotePeer().String())
			go pingService.Ping(ctx, conn.RemotePeer())
		},
	})

	return h, nil
}

// SendFile initiates a stream to a peer and sends the file's name and content
func SendFile(ctx context.Context, h host.Host, peerID peer.ID, filename string, data []byte) error {

	stream, err := h.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		return fmt.Errorf("error creating new stream: %w", err)
	}
	defer func(stream network.Stream) {
		err := stream.Close()
		if err != nil {
			log.Printf("error closing stream: %s", err.Error())
		}
	}(stream)

	writer := bufio.NewWriter(stream)
	// Send the filename first
	_, err = writer.WriteString(filename + "\n")
	if err != nil {
		return fmt.Errorf("error writing filename: %w", err)
	}

	// Send the file data
	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("error writing file data: %w", err)
	}

	// Flush the buffer to ensure all data is sent
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing data: %w", err)
	}

	log.Printf("File '%s' sent to peer %s\n", filename, peerID.String())
	return nil
}

// ReceiveFile handles an incoming stream and reads the file's name and content
func ReceiveFile(stream network.Stream) (string, []byte, error) {
	reader := bufio.NewReader(stream)

	// Read the filename
	filename, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, fmt.Errorf("error reading filename: %w", err)
	}
	filename = filename[:len(filename)-1]

	var data []byte
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", nil, fmt.Errorf("error reading file data: %w", err)
		}
		data = append(data, buf[:n]...)
	}

	log.Printf("Received file '%s' of size %d bytes\n", filename, len(data))
	return filename, data, nil
}

// HandleStream is the stream handler for incoming file transfer streams
func HandleStream(stream network.Stream) {
	defer func(stream network.Stream) {
		err := stream.Close()
		if err != nil {
			log.Printf("error closing stream")
		}
	}(stream)

	log.Printf("New stream opened with peer %s\n", stream.Conn().RemotePeer().String())

	// Receive file
	filename, data, err := ReceiveFile(stream)
	if err != nil {
		log.Printf("Error receiving file: %s\n", err)
		return
	}

	log.Printf("Successfully received file: %s, Size: %d bytes\n", filename, len(data))
}
