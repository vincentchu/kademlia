package utils

import (
	"fmt"
	"log"

	multiaddr "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	peer "gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// MakePeer takes a fully-encapsulated address and converts it to a
// peer ID / Multiaddress pair
func MakePeer(dest string) (peer.ID, multiaddr.Multiaddr) {
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
