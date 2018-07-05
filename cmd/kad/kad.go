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
	net "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	multiaddr "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	libp2p "gx/ipfs/QmZ86eLPtXkQ1Dfa992Q8NpXArUoWWh3y728JDcWvzRrvC/go-libp2p"
	peerstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	// record "gx/ipfs/QmVsp2KdPYE6M8ryzCk5KHLo3zprcY5hBDaYx6uPCFUdxA/go-libp2p-record"
)

func dieIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func addrForPort(p string) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%s", p))
}

func streamHandlerFor(name string, h host.Host, kad *dht.IpfsDHT) func(s net.Stream) {
	fn := func(s net.Stream) {
		conn := s.Conn()
		remotePeer := conn.RemotePeer()
		remoteAddr := conn.RemoteMultiaddr()

		kad.Update(context.Background(), remotePeer)
		h.Peerstore().AddAddr(remotePeer, remoteAddr, peerstore.PermanentAddrTTL)

		log.Printf("Opened new stream %s: %v", name, s.Protocol())
		log.Printf("  Local Addr:  %s", conn.LocalMultiaddr().String())
		log.Printf("  Remote Addr: %s", conn.RemoteMultiaddr().String())
		log.Printf("  Remote Peer: %s", conn.RemotePeer().Pretty())

		peerInfo, err := kad.FindPeer(context.Background(), conn.RemotePeer())
		if err != nil {
			log.Printf("Error when finding peer: %v\n", err)
		}

		log.Printf("Found peerinfo: %v, addrs: %v\n", peerInfo.ID, peerInfo.Addrs)

		err = kad.PutValue(context.Background(), "/hello/KEY", []byte{1, 2, 3, 4, 5})
		if err != nil {
			log.Printf("Error on PutValue: %v\n", err)
		}

		get, err := kad.GetValue(context.Background(), "/hello/KEY")
		if err != nil {
			log.Printf("Error when getting: %v\n", err)
		}

		log.Printf("GET result: %v\n", get)

		// ch, err := kad.GetClosestPeers(context.Background(), "/hello/KEY")
		// if err != nil {
		// 	log.Printf("Error when getting closest peers: %v", err)
		// }

		// go func() {
		// 	closestPeer := <-ch
		// 	log.Printf("Closest peerID: %s\n", closestPeer)
		// }()

	}

	return fn
}

func streamHandler(s net.Stream) {
	fmt.Println("StreamHandler")
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

	host.Peerstore().AddAddr(host.ID(), host.Addrs()[0], peerstore.PermanentAddrTTL)
	prettyHostID := host.ID().Pretty()

	kadDHT, err := dht.New(ctx, host, dhtopts.Validator(utils.NullValidator{}))
	dieIfError(err)
	host.SetStreamHandler(dhtopts.ProtocolDHT, streamHandler)

	ipfsAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", prettyHostID))
	fullHostAddr := hostAddr.Encapsulate(ipfsAddr)
	fmt.Println("Generated host: ", host.ID().Pretty())
	fmt.Printf("Host Address: %v\n", fullHostAddr)

	return host, kadDHT
}

func addPeers(ctx context.Context, peerStr string, h host.Host, kad *dht.IpfsDHT) {
	if len(peerStr) == 0 {
		return
	}

	peerStrs := strings.Split(peerStr, ",")
	for i := 0; i < len(peerStrs); i++ {
		peerID, peerAddr := utils.MakePeer(peerStrs[i])

		h.Peerstore().AddAddr(peerID, peerAddr, peerstore.PermanentAddrTTL)
		kad.Update(ctx, peerID)
	}
}

func main() {
	fmt.Println("Kademlia DHT test")

	port := flag.Int64("port", 0, "Port to listen on")
	peers := flag.String("peers", "", "Initial peers")
	flag.Parse()

	ctx := context.Background()
	srvHost, kad := generateHost(ctx, *port)
	addPeers(ctx, *peers, srvHost, kad)

	_ = peers
	_ = srvHost
	_ = kad

	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", srvHost.ID().Pretty()))
	srvAddr := srvHost.Addrs()[0]
	fullSrvAddr := srvAddr.Encapsulate(hostAddr)
	log.Printf("Full Server Address 2: %s\n", fullSrvAddr)

	// srvHost.Peerstore().AddAddr(srvHost.ID(), srvAddr, peerstore.PermanentAddrTTL)
	// kad.Update(ctx, srvHost.ID())

	// _ = kad

	// addPeers(srvHost, *peers)

	fmt.Printf("Listening on %v\n", srvHost.Addrs())
	fmt.Printf("Protocols supported: %v\n", srvHost.Mux().Protocols())

	<-make(chan struct{})

	// srcHost := generateHost(ctx, 3001)
	// fmt.Println(srcHost.ID().Pretty())

	// // dataStore := memstore.NewIntMemstore()

	// kadDHT, _ := dht.New(ctx, srcHost)

	// peers, _ := kadDHT.GetClosestPeers(ctx, "foo")

	// fmt.Println("Close peers ", peers)
}
