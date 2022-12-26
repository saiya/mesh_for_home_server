package config

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/saiya/mesh_for_home_server/logger"
)

type NodeID string

const nodeIDrandLength = 32
const nodeIDLetters = "23456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

func GenerateNodeID(hostname string) NodeID {
	if hostname == "" {
		var err error
		hostname, err = os.Hostname()
		if err != nil {
			logger.Get().Infow("Failed to get hostname: " + err.Error())
			hostname = "hostname-NA"
		}
	}

	randPart := make([]byte, nodeIDrandLength)
	for i := range randPart {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(nodeIDLetters))))
		if err != nil {
			// We generate NodeID only in the start of process.
			// And if system's entropy pool not working, anyway we unlikely able to initialize TLS (for gRPC, HTTPS, ...) anyway.
			panic("failed to read from system's secure random source: " + err.Error())
		}
		randPart[i] = nodeIDLetters[num.Int64()]
	}
	return NodeID(fmt.Sprintf("%s-%s", hostname, string(randPart)))
}
