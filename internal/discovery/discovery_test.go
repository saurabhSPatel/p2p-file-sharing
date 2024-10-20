package discovery

import (
	"context"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"reflect"
	"testing"
)

func TestDiscovery_DiscoverPeers(t *testing.T) {
	type fields struct {
		host host.Host
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    <-chan peer.AddrInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discovery{
				host: tt.fields.host,
			}
			got, err := d.DiscoverPeers(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("DiscoverPeers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DiscoverPeers() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiscovery_HandlePeerFound(t *testing.T) {
	type fields struct {
		host host.Host
	}
	type args struct {
		pi peer.AddrInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discovery{
				host: tt.fields.host,
			}
			d.HandlePeerFound(tt.args.pi)
		})
	}
}

func TestDiscovery_SetupDiscovery(t *testing.T) {
	type fields struct {
		host host.Host
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discovery{
				host: tt.fields.host,
			}
			if err := d.SetupDiscovery(); (err != nil) != tt.wantErr {
				t.Errorf("SetupDiscovery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDiscovery_findPeers(t *testing.T) {
	type fields struct {
		host host.Host
	}
	type args struct {
		ctx      context.Context
		peerChan chan<- peer.AddrInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discovery{
				host: tt.fields.host,
			}
			d.findPeers(tt.args.ctx, tt.args.peerChan)
		})
	}
}

func TestNewDiscovery(t *testing.T) {
	type args struct {
		h host.Host
	}
	tests := []struct {
		name string
		args args
		want *Discovery
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDiscovery(tt.args.h); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDiscovery() = %v, want %v", got, tt.want)
			}
		})
	}
}
