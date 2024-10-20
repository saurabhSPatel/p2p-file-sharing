
# Peer-to-Peer File Sharing Application

A simple peer-to-peer file sharing application built using Golang and libp2p. This application allows peers to discover each other, upload and download files directly without a central server.

## Features

- **Peer Discovery**: Utilize libp2p's mDNS for discovering peers on the network.
- **File Sharing**: Upload and download files between peers using secure streams.
- **Command-Line Interface (CLI)**: User-friendly CLI for interacting with the application.
- **Graceful Shutdown**: Handles shutdown without data loss or corruption.
- **Testing**: Includes unit, Integration tests for critical components.

## Future Enhancements

- **Unit Tests**: Implement unit tests for all components to ensure reliability.
- **Improved Logging**: Replace the current logging mechanism with [zerolog](https://github.com/rs/zerolog) for structured logging and better performance.
- **Security Features**: Add support for encryption and authentication using libp2p's security protocols to enhance data security.
- **Real-Time Data Exchange**: Implement support for peer streaming and real-time data exchange using libp2p pubsub for dynamic interactions.

## Project Structure

```
p2p-file-sharing/
├── cmd/
│   └── p2pfs/
│       └── main.go            # Entry point of the application
├── internal/
│   ├── discovery/              # Peer discovery logic
│   ├── file/                   # File handling utilities
│   ├── network/                # Networking setup and communication
│   └── cli/                    # Command-line interface implementation
├── pkg/
│   └── utils/                  # Utility functions
├── config/
│   └── config.go              # Configuration settings
├── test/
│   └── integration_test.go     # Integration tests
├── go.mod                      # Go module file
├── go.sum                      # Go module dependencies
└── README.md                   # Project documentation
```

## Requirements

- Go (version 1.20 or higher)
- Docker (for deployment)
- Docker Compose (for orchestrating multiple peers)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/saurabhSPatel/p2p-file-sharing.git
   cd p2p-file-sharing
   ```

2. Build the application:

   ```bash
   go build ./cmd/p2pfs
   ```

3. (Optional) To run the application in Docker:

   Ensure Docker and Docker Compose are installed, then run:

   ```bash
   ./deploy.sh <number_of_peers>
   ```

## Usage

1. Run the application:

   If not using Docker, execute the compiled binary:

   ```bash
   ./p2pfs
   ```

2. Use the CLI commands:

   - `list`: List available files in the shared directory.
   - `upload <filename>`: Upload a file to a peer.
   - `download <filename>`: Download a file from a peer.
   - `exit`: Exit the CLI.

## Testing

Run the tests using the following command:

```bash
go test ./...
```

## Troubleshooting

- Ensure all dependencies are installed and up to date.
- Check the logs of each peer for any connection errors.
- Make sure the shared directories exist on each peer.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs or feature requests.
