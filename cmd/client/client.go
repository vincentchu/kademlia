package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	// "time"

	dht "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht"
	dhtopts "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht/opts"
	// net "gx/ipfs/QmPjvxTpVH8qJyQDnxnsxF9kv9jezKD1kozz1hs3fCGsNh/go-libp2p-net"

	// protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"

	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/vincentchu/kademlia/utils"
	// record "gx/ipfs/QmVsp2KdPYE6M8ryzCk5KHLo3zprcY5hBDaYx6uPCFUdxA/go-libp2p-record"
)

func makeHost(ctx context.Context) host.Host {
	prvKey := utils.GeneratePrivateKey(999)

	h, err := libp2p.New(ctx, libp2p.Identity(prvKey))
	if err != nil {
		log.Fatalf("Err on creating host: %v\n", err)
	}

	return h
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
	destID, destAddr := utils.MakePeer(*dest)

	h.Peerstore().AddAddr(destID, destAddr, 24*time.Hour)
	kad, err := dht.New(ctx, h, dhtopts.Client(true), dhtopts.Validator(utils.NullValidator{}))
	if err != nil {
		log.Fatalf("Error creating DHT: %v\n", err)
	}

	fmt.Println(destID.Pretty())
	h.NewStream(ctx, destID, dhtopts.ProtocolDHT, dhtopts.ProtocolDHTOld)
	kad.Update(ctx, destID)

	cmd, key, val := parseCmd(flag.Args())
	switch cmd {
	case "put":
		log.Printf("PUT %s => %s\n", key, val)
		err = kad.PutValue(ctx, key, []byte(val))
		if err != nil {
			log.Fatalf("Error on PUT: %v\n", err)
		}

	case "get":
		log.Printf("GET %s", key)
		fetchedBytes, err := kad.GetValue(ctx, key, dht.Quorum(1))
		if err != nil {
			log.Fatalf("Error on GET: %v\n", err)
		}
		log.Printf("RESULT: %s\n", string(fetchedBytes))

	default:
		log.Fatalf("Command %s unrecognized\n", cmd)
	}
}
