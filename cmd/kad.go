package main

import (
	"context"
	"fmt"
	"gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	"gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	"gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"log"
	"math/rand"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

func generateHost(ctx context.Context, port int64) host.Host {
	randBytes := rand.New(rand.NewSource(port))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randBytes)
	if err != nil {
		log.Fatal(err)
	}

	hostAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	host, err := libp2p.New(ctx, libp2p.ListenAddrs(hostAddr), libp2p.Identity(prvKey))

	if err != nil {
		log.Fatal(err)
	}

	return host
}

func main() {
	fmt.Println("Kademlia DHT test")

	ctx := context.Background()

	srcHost := generateHost(ctx, 3001)
	fmt.Println(srcHost.ID().Pretty())

	// dataStore := memstore.NewIntMemstore()

	dht.NewDHT
}
