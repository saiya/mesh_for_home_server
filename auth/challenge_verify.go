package auth

import (
	"fmt"
	"time"

	"github.com/saiya/mesh_for_home_server/config"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/exp/slices"
)

const challengeBytes = 64
const challengeLetters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const clockSkewMax = time.Second * 10
const challengeExpire = time.Second * 5

type Challenge struct {
	NodeID config.NodeID
	Nonce  string
}

type Answer string

type Challenger struct {
	nodeID  config.NodeID
	privKey jwk.Key
}

type Verifier struct {
	pubKeys jwk.Set
}

func NewChallenger(nodeID config.NodeID, jwkStr string) (*Challenger, error) {
	challenger := Challenger{
		nodeID: nodeID,
	}
	var err error

	challenger.privKey, err = parseJwk(jwkStr, jwk.ForSignature, jwk.KeyOpSign)
	if err != nil {
		return nil, fmt.Errorf("failed to parse challenger JWK: %w", err)
	}

	return &challenger, nil
}

func NewVerifier(jwkStrs []string) (*Verifier, error) {
	var verifier Verifier

	verifier.pubKeys = jwk.NewSet()
	for _, jwkStr := range jwkStrs {
		pubKey, err := parseJwk(jwkStr, jwk.ForSignature, jwk.KeyOpVerify)
		if err != nil {
			return nil, fmt.Errorf("failed to parse verifier JWK: %w", err)
		}
		neverFail(verifier.pubKeys.AddKey(pubKey))
	}

	return &verifier, nil
}

func parseJwk(jwkStr string, use jwk.KeyUsageType, ops jwk.KeyOperation) (jwk.Key, error) {
	key, err := jwk.ParseKey([]byte(jwkStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWK: %w", err)
	}
	if key.KeyID() == "" {
		return nil, fmt.Errorf("must have kid")
	}
	if key.KeyUsage() != string(use) && !slices.Contains(key.KeyOps(), ops) {
		return nil, fmt.Errorf("must have use=\"%v\" or key_ops=\"%v\"", use, ops)
	}
	return key, nil
}

func (v *Verifier) NewChallenge(nodeID config.NodeID) (*Challenge, error) {
	return &Challenge{
		NodeID: nodeID,
		Nonce:  SecureRandomString(challengeBytes, challengeLetters),
	}, nil
}

func (c *Challenger) Answer(challenge *Challenge, now time.Time) (Answer, error) {
	if challenge.NodeID != c.nodeID {
		return "", fmt.Errorf("Node ID not match")
	}

	token, err := jwt.NewBuilder().
		JwtID(challenge.Nonce).
		Issuer(string(c.nodeID)).
		IssuedAt(now).
		NotBefore(now).
		Expiration(now.Add(challengeExpire)).
		Build()
	neverFail(err)

	jws, err := jwt.Sign(token, jwt.WithKey(c.privKey.Algorithm(), c.privKey))
	neverFail(err)
	return Answer(jws), nil
}

func (v *Verifier) Verify(challenge *Challenge, answer Answer, now time.Time) error {
	_, err := jwt.Parse(
		[]byte(answer),
		jwt.WithKeySet(v.pubKeys),
		jwt.WithIssuer(string(challenge.NodeID)),
		jwt.WithJwtID(challenge.Nonce),
		jwt.WithAcceptableSkew(clockSkewMax),
		jwt.WithClock(jwt.ClockFunc(func() time.Time { return now })),
	)
	if err != nil {
		return fmt.Errorf("invalid answer (JWS parse/verification not passed): %w", err)
	}
	return nil
}
