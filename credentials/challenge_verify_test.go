package credentials

import (
	"testing"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	keyPair, err := NewSignKeyPair()
	if !assert.NoError(t, err) {
		return
	}

	nodeID := config.NodeID("test-node")
	challenger, err := NewChallenger(nodeID, string(keyPair.PrivateKey))
	if !assert.NoError(t, err) {
		return
	}
	verifier, err := NewVerifier([]string{string(keyPair.PublicKey)})
	if !assert.NoError(t, err) {
		return
	}

	challenge, err := verifier.NewChallenge(nodeID)
	if !assert.NoError(t, err) {
		return
	}
	answer, err := challenger.Answer(challenge)
	if !assert.NoError(t, err) {
		return
	}
	err = verifier.Verify(challenge, answer)
	if !assert.NoError(t, err) {
		return
	}
}
