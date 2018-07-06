package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vincentchu/kademlia/utils"

	dht "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht"
	dhtopts "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht/opts"
	multiaddr "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	libp2p "gx/ipfs/QmZ86eLPtXkQ1Dfa992Q8NpXArUoWWh3y728JDcWvzRrvC/go-libp2p"
	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	peer "gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

func dieIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func addrForPort(p string) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", p))
}

func generateHost(ctx context.Context, port int64) (host.Host, *dht.IpfsDHT) {
	prvKey := utils.GeneratePrivateKey(port)

	hostAddr, err := addrForPort(fmt.Sprintf("%d", port))
	dieIfError(err)

	opts := []libp2p.Option{
		libp2p.ListenAddrs(hostAddr),
		libp2p.Identity(prvKey),
	}

	host, err := libp2p.New(ctx, opts...)
	dieIfError(err)

	kadDHT, err := dht.New(ctx, host, dhtopts.Validator(utils.NullValidator{}))
	dieIfError(err)

	log.Printf("Generated Host: %s/ipfs/%s\n", host.Addrs()[0].String(), host.ID().Pretty())

	return host, kadDHT
}

func addPeers(h host.Host, peerStr string) {
	if len(peerStr) == 0 {
		return
	}

	portStrs := strings.Split(peerStr, ",")
	for i := 0; i < len(portStrs); i++ {
		addr, err := addrForPort(portStrs[i])
		dieIfError(err)
		pid := "QmcxsSTeHBEfaWBb2QKe5UZWK8ezWJkxJfmcb5rQV374M6" //peer.ID(fmt.Sprintf("QmcxsSTeHBEfaWBb2QKe5UZWK8ezWJkxJfmcb5rQV374M6", portStrs[i]))
		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			fmt.Printf("Decode pid %v\n", err)
		}

		h.Peerstore().AddAddr(peerid, addr, 24*time.Hour)
		_, err = h.NewStream(context.Background(), peerid, "/multistream/1.0.0", "/ipfs/id/1.0.0", "/ipfs/kad/1.0.0", "/ipfs/dht")
		fmt.Printf("Error on new stream: %v\n", err)
	}
}

func main() {
	log.Println("Kademlia DHT test")

	port := flag.Int64("port", 0, "Port to listen on")
	_ = flag.String("peers", "", "Initial peers")
	flag.Parse()

	ctx := context.Background()
	srvHost, kad := generateHost(ctx, *port)
	_ = kad

	// addPeers(srvHost, *peers)

	log.Printf("Listening on %v\n", srvHost.Addrs())
	log.Printf("Protocols supported: %v\n", srvHost.Mux().Protocols())

	<-make(chan struct{})
}
