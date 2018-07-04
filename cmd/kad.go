package main

import (
	"context"
	"fmt"
	"log"

	dht "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht"
	libp2p "gx/ipfs/QmZ86eLPtXkQ1Dfa992Q8NpXArUoWWh3y728JDcWvzRrvC/go-libp2p"
	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"

	"github.com/vincentchu/kademlia/memstore"
)

func generateHost(ctx context.Context, port int64) host.Host {
	// randBytes := rand.New(rand.NewSource(port))
	// prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randBytes)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// hostAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	host, err := libp2p.New(ctx)

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

	dataStore := memstore.NewIntMemstore()

	dht.NewDHT(ctx, srcHost, dataStore)
}
