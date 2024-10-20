package discovery

import (
	"context"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// Discovery manages peer discovery in the network.
type Discovery struct {
	host host.Host
}

// NewDiscovery creates a new instance of Discovery.
func NewDiscovery(h host.Host) *Discovery {
	return &Discovery{host: h}
}

// SetupDiscovery initializes the mDNS discovery service.
func (d *Discovery) SetupDiscovery() error {
	service := mdns.NewMdnsService(d.host, "p2p-file-sharing", d)
	return service.Start()
}

// DiscoverPeers returns a channel for discovering peers.
func (d *Discovery) DiscoverPeers(ctx context.Context) (<-chan peer.AddrInfo, error) {
	peerChan := make(chan peer.AddrInfo)
	go d.findPeers(ctx, peerChan)
	return peerChan, nil
}

// findPeers periodically discovers peers and sends them to the channel.
func (d *Discovery) findPeers(ctx context.Context, peerChan chan<- peer.AddrInfo) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers := d.host.Network().Peers()
			for _, p := range peers {
				select {
				case peerChan <- d.host.Peerstore().PeerInfo(p):
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// HandlePeerFound connects to a newly discovered peer.
func (d *Discovery) HandlePeerFound(pi peer.AddrInfo) {
	if err := d.host.Connect(context.Background(), pi); err != nil {
		log.Printf("Error connecting to peer %s: %v", pi.ID, err)
	}
}
