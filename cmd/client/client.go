package main

import (
	"context"
	"flag"
	"log"
	"math/rand"

	"github.com/vincentchu/kademlia/utils"

	dht "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht"
	dhtopts "gx/ipfs/QmNg6M98bwS97SL9ArvrRxKujFps3eV6XvmKgduiYga8Bn/go-libp2p-kad-dht/opts"
	peerstore "gx/ipfs/QmZR2XWVVBCtbgBWnQhWk2xcQfaR3W8faQPriAiaaj7rsr/go-libp2p-peerstore"
	host "gx/ipfs/Qmb8T6YBBsjYsVGfrihQLfCJveczZnneSBqBKkYEBWDjge/go-libp2p-host"
	crypto "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"

	libp2p "github.com/libp2p/go-libp2p"
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

	h.Peerstore().AddAddr(destID, destAddr, peerstore.PermanentAddrTTL)

	kad, err := dht.New(ctx, h, dhtopts.Client(true), dhtopts.Validator(NullValidator{}))
	if err != nil {
		log.Fatalf("Error creating DHT: %v\n", err)
	}

	kad.Update(ctx, destID)

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
		fetchedBytes, err := kad.GetValue(ctx, key, dht.Quorum(1))
		if err != nil {
			log.Fatalf("Error on GET: %v\n", err)
		}
		log.Printf("RESULT: %s\n", string(fetchedBytes))

	default:
		log.Fatalf("Command %s unrecognized\n", cmd)
	}
}
