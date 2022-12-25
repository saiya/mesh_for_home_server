package config

import (
	"fmt"
	"math/rand"
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
		randPart[i] = nodeIDLetters[rand.Intn(len(nodeIDLetters))]
	}
	return NodeID(fmt.Sprintf("%s-%s", hostname, string(randPart)))
}
