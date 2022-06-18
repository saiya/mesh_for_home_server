package credentials

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
	nodeID config.NodeID
	Nonce  string `json:"nonce"`
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
		if err = verifier.pubKeys.AddKey(pubKey); err != nil {
			return nil, fmt.Errorf("failed to add JWK into set: %w", err)
		}
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
		nodeID: nodeID,
		Nonce:  NewNonce(challengeLetters),
	}, nil
}

func (c *Challenger) Answer(challenge *Challenge) (Answer, error) {
	now := time.Now()
	token, err := jwt.NewBuilder().
		JwtID(challenge.Nonce).
		Issuer(string(c.nodeID)).
		IssuedAt(now).
		NotBefore(now).
		Expiration(now.Add(challengeExpire)).
		Build()
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	jws, err := jwt.Sign(token, jwt.WithKey(c.privKey.Algorithm(), c.privKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return Answer(jws), nil
}

func (v *Verifier) Verify(challenge *Challenge, answer Answer) error {
	_, err := jwt.Parse(
		[]byte(answer),
		jwt.WithKeySet(v.pubKeys),
		jwt.WithIssuer(string(challenge.nodeID)),
		jwt.WithJwtID(challenge.Nonce),
		jwt.WithAcceptableSkew(clockSkewMax),
	)
	if err != nil {
		return fmt.Errorf("failed to parse JWS: %w", err)
	}
	return nil
}
