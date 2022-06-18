package credentials

import (
	"testing"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	challenger, verifier, challenge, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	now := time.Now()
	answer, err := challenger.Answer(challenge, now)
	if !assert.NoError(t, err) {
		return
	}
	err = verifier.Verify(challenge, answer, now)
	if !assert.NoError(t, err) {
		return
	}
}

func TestInvalidJWK(t *testing.T) {
	nodeID := config.NodeID("test-node")

	_, err := NewChallenger(nodeID, "invalid-jwk")
	assert.ErrorContains(t, err, "failed to parse challenger JWK")

	_, err = NewChallenger(nodeID, `{"alg":"EdDSA","crv":"Ed25519","d":"-4OhFgdUf7k5k6eV_LvFlSdmB5ea0uBfNtYME7tQ8j8","key_ops":["sign"],"kty":"OKP","x":"qhw4VxKZMoHD0kaWma5bK6bTHXhO2-w4V-y8AUmlYLA"}`)
	assert.ErrorContains(t, err, "must have kid")

	_, err = NewChallenger(nodeID, `{"alg":"EdDSA","crv":"Ed25519","d":"-4OhFgdUf7k5k6eV_LvFlSdmB5ea0uBfNtYME7tQ8j8","kid":"key-m2fK7TlbVgwgexqg3hrkLslmKhekGwk5tNR1oUaoeiNwglt7Z7Gi9LijefavkDQz","kty":"OKP","x":"qhw4VxKZMoHD0kaWma5bK6bTHXhO2-w4V-y8AUmlYLA"}`)
	assert.ErrorContains(t, err, "must have use=\"sig\" or key_ops=\"sign\"")

	_, err = NewVerifier([]string{"invalid-jwk"})
	assert.ErrorContains(t, err, "failed to parse verifier JWK")
}

func TestChallengeDifferentNodeID(t *testing.T) {
	challenger, _, challenge, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	challenge.NodeID = config.NodeID("another-node")
	now := time.Now()
	_, err = challenger.Answer(challenge, now.Add(-challengeExpire-time.Second))
	assert.ErrorContains(t, err, "Node ID not match")
}

func TestAnswerExpired(t *testing.T) {
	challenger, verifier, challenge, err := setup(t)
	if !assert.NoError(t, err) {
		return
	}

	now := time.Now()
	answer, err := challenger.Answer(challenge, now.Add(-challengeExpire-time.Second))
	if !assert.NoError(t, err) {
		return
	}
	err = verifier.Verify(challenge, answer, now.Add(clockSkewMax))
	assert.ErrorContains(t, err, "invalid answer (JWS parse/verification not passed): \"exp\" not satisfied")
}

func setup(t *testing.T) (*Challenger, *Verifier, *Challenge, error) {
	keyPair := NewSignKeyPair()

	nodeID := config.NodeID("test-node")
	challenger, err := NewChallenger(nodeID, string(keyPair.PrivateKey))
	if !assert.NoError(t, err) {
		return nil, nil, nil, err
	}
	verifier, err := NewVerifier([]string{string(keyPair.PublicKey)})
	if !assert.NoError(t, err) {
		return nil, nil, nil, err
	}

	challenge, err := verifier.NewChallenge(nodeID)
	if !assert.NoError(t, err) {
		return nil, nil, nil, err
	}

	return challenger, verifier, challenge, nil
}
