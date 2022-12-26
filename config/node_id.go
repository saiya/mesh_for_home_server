package config

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

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

	randSeed := make([]byte, 8)
	_, err := io.ReadFull(cryptorand.Reader, randSeed)
	if err != nil {
		// Still system time provides some randomness
		logger.Get().Infow("Failed to get randSeed: " + err.Error())
	}
	rng := rand.New(rand.NewSource(time.Now().Unix() ^ int64(binary.BigEndian.Uint64(randSeed))))
	randPart := make([]byte, nodeIDrandLength)
	for i := range randPart {
		randPart[i] = nodeIDLetters[rng.Intn(len(nodeIDLetters))]
	}
	return NodeID(fmt.Sprintf("%s-%s", hostname, string(randPart)))
}
