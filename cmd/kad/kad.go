package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/vincentchu/kademlia/utils"

	dht "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht"
	dhtopts "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht/opts"
	multiaddr "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	libp2p "gx/ipfs/QmZ86eLPtXkQ1Dfa992Q8NpXArUoWWh3y728JDcWvzRrvC/go-libp2p"
	peerstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
)

func addrForPort(p string) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", p))
}

func generateHost(ctx context.Context, port int64) (host.Host, *dht.IpfsDHT) {
	prvKey := utils.GeneratePrivateKey(port)

	hostAddr, err := addrForPort(fmt.Sprintf("%d", port))
	if err != nil {
		log.Fatal(err)
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrs(hostAddr),
		libp2p.Identity(prvKey),
	}

	host, err := libp2p.New(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}

	kadDHT, err := dht.New(ctx, host, dhtopts.Validator(utils.NullValidator{}))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Generated Host: %s/ipfs/%s\n", host.Addrs()[0].String(), host.ID().Pretty())

	return host, kadDHT
}

func addPeers(ctx context.Context, h host.Host, kad *dht.IpfsDHT, peersArg string) {
	if len(peersArg) == 0 {
		return
	}

	peerStrs := strings.Split(peersArg, ",")
	for i := 0; i < len(peerStrs); i++ {
		peerID, peerAddr := utils.MakePeer(peerStrs[i])

		h.Peerstore().AddAddr(peerID, peerAddr, peerstore.PermanentAddrTTL)
		kad.Update(ctx, peerID)
	}
}

func main() {
	log.Println("Kademlia DHT test")

	port := flag.Int64("port", 3001, "Port to listen on")
	peers := flag.String("peers", "", "Initial peers")
	flag.Parse()

	ctx := context.Background()
	srvHost, kad := generateHost(ctx, *port)

	addPeers(ctx, srvHost, kad, *peers)

	log.Printf("Listening on %v\n", srvHost.Addrs())
	log.Printf("Protocols supported: %v\n", srvHost.Mux().Protocols())

	<-make(chan struct{})
}
