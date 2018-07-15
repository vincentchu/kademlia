package utils

import (
	"fmt"
	"math/rand"

	multiaddr "gx/ipfs/QmYmsdtJ3HsodkePE3eU3TsCaP2YvPZJ4LoXnNkDE5Tpt7/go-multiaddr"
	logging "gx/ipfs/QmcVVHfdyv15GVPk7NrxdWjh2hLVccXnoD8j2tyQShiXJb/go-log"
	peer "gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	crypto "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

var log = logging.Logger("kadutils")

// MakePeer takes a fully-encapsulated address and converts it to a
// peer ID / Multiaddress pair
func MakePeer(dest string) (peer.ID, multiaddr.Multiaddr) {
	ipfsAddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		log.Fatalf("Err on creating host: %v", err)
	}
	log.Debugf("Parsed: ipfsAddr = %s", ipfsAddr)

	peerIDStr, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		log.Fatalf("Err on creating peerIDStr: %v", err)
	}
	log.Debugf("Parsed: PeerIDStr = %s", peerIDStr)

	peerID, err := peer.IDB58Decode(peerIDStr)
	if err != nil {
		log.Fatalf("Err on decoding %s: %v", peerIDStr, err)
	}
	log.Debugf("Created peerID = %s", peerID)

	targetPeerAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerID)))
	log.Debugf("Created targetPeerAddr = %v", targetPeerAddr)

	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
	log.Debugf("Decapsuated = %v", targetAddr)

	return peerID, targetAddr
}

// NullValidator is a validator that does no valiadtion
type NullValidator struct{}

// Validate always returns success
func (nv NullValidator) Validate(key string, value []byte) error {
	log.Debugf("NullValidator Validate: %s - %s", key, string(value))
	return nil
}

// Select always selects the first record
func (nv NullValidator) Select(key string, values [][]byte) (int, error) {
	strs := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		strs[i] = string(values[i])
	}
	log.Debugf("NullValidator Select: %s - %v", key, strs)

	return 0, nil
}

// GeneratePrivateKey - creates a private key with the given seed
func GeneratePrivateKey(seed int64) crypto.PrivKey {
	randBytes := rand.New(rand.NewSource(seed))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randBytes)

	if err != nil {
		log.Fatalf("Could not generate Private Key: %v", err)
	}

	return prvKey
}
