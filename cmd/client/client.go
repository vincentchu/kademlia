package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	// "time"

	dht "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht"
	dhtopts "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht/opts"
	// net "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"
	multiaddr "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	// protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	peerstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	peer "gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"

	crypto "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"

	libp2p "github.com/libp2p/go-libp2p"
	// record "gx/ipfs/QmVsp2KdPYE6M8ryzCk5KHLo3zprcY5hBDaYx6uPCFUdxA/go-libp2p-record"
)

// NullValidator is a validator that does no valiadtion
type NullValidator struct{}

// Validate always returns success
func (nv NullValidator) Validate(key string, value []byte) error {
	log.Printf("NullValidator Validate: %s - %v\n", key, value)
	return nil
}

// Select always selects the first record
func (nv NullValidator) Select(key string, values [][]byte) (int, error) {
	log.Printf("NullValidator Select: %s - %v\n", key, values)
	log.Printf("NullValidator Select: %d", len(values))

	return 0, nil
}

func makeHost(ctx context.Context) host.Host {
	randBytes := rand.New(rand.NewSource(999))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randBytes)

	h, err := libp2p.New(ctx, libp2p.Identity(prvKey))
	if err != nil {
		log.Fatalf("Err on creating host: %v\n", err)
	}

	return h
}

func makePeer(dest string) (peer.ID, multiaddr.Multiaddr) {
	ipfsAddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		log.Fatalf("Err on creating host: %v\n", err)
	}
	log.Printf("Parsed: ipfsAddr = %s\n", ipfsAddr)

	peerIDStr, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		log.Fatalf("Err on creating peerIDStr: %v\n", err)
	}
	log.Printf("Parsed: PeerIDStr = %s\n", peerIDStr)

	peerID, err := peer.IDB58Decode(peerIDStr)
	if err != nil {
		log.Fatalf("Err on decoding %s: %v\n", peerIDStr, err)
	}
	log.Printf("Created peerID = %s\n", peerID)

	targetPeerAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerID)))
	log.Printf("Created targetPeerAddr = %v\n", targetPeerAddr)

	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
	log.Printf("Decapsuated = %v\n", targetAddr)

	return peerID, targetAddr
}

func parseCmd(tokens []string) (string, string, string) {
	switch len(tokens) {
	case 2:
		return tokens[0], tokens[1], ""
	case 3:
		return tokens[0], tokens[1], tokens[2]
	default:
		log.Fatalf("Improper command format: %v\n", tokens)
		return "", "", ""
	}
}

func main() {
	dest := flag.String("dest", "", "Destination to connect to")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h := makeHost(ctx)
	destID, destAddr := makePeer(*dest)

	h.Peerstore().AddAddr(destID, destAddr, peerstore.PermanentAddrTTL)

	kad, err := dht.New(ctx, h, dhtopts.Client(true), dhtopts.Validator(NullValidator{}))
	if err != nil {
		log.Fatalf("Error creating DHT: %v\n", err)
	}

	cmd, key, val := parseCmd(flag.Args())
	switch cmd {
	case "put":
		log.Printf("PUT %s => %s\n", key, val)
		err = kad.PutValue(ctx, key, []byte(val), dht.Quorum(0))
		if err != nil {
			log.Fatalf("Error on PUT: %v\n", err)
		}

	case "get":
		log.Printf("GET %s", key)
		fetchedBytes, err := kad.GetValue(ctx, key, dht.Quorum(0))
		if err != nil {
			log.Fatalf("Error on GET: %v\n", err)
		}
		log.Printf("RESULT: %s\n", string(fetchedBytes))

	default:
		log.Fatalf("Command %s unrecognized\n", cmd)
	}
}
