package credentials

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

const kidRandomChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewSignKeyPair() (KeyPair, error) {
	return NewKeyPair(
		jwk.KeyOperationList{jwk.KeyOpVerify},
		jwk.KeyOperationList{jwk.KeyOpSign},
	)
}

func NewKeyPair(pubOps jwk.KeyOperationList, privOps jwk.KeyOperationList) (KeyPair, error) {
	alg := "EdDSA"
	kid := "key-" + NewNonce(kidRandomChars)

	rawPubKey, rawPrivKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to generate ECDSA key: %w", err)
	}

	privKey, err := jwk.FromRaw(rawPrivKey)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to convert to JWK: %w", err)
	}

	pubKey, err := jwk.FromRaw(rawPubKey)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to convert to JWK: %w", err)
	}

	err = pubKey.Set(jwk.KeyOpsKey, pubOps)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to set JWK property: %w", err)
	}
	err = privKey.Set(jwk.KeyOpsKey, privOps)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to set JWK property: %w", err)
	}
	for _, k := range []jwk.Key{pubKey, privKey} {
		err = k.Set(jwk.AlgorithmKey, alg)
		if err != nil {
			return KeyPair{}, fmt.Errorf("failed to set JWK property: %w", err)
		}
		err = k.Set(jwk.KeyIDKey, kid)
		if err != nil {
			return KeyPair{}, fmt.Errorf("failed to set JWK property: %w", err)
		}
	}

	var pair KeyPair
	pair.PublicKey, err = json.Marshal(pubKey)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to serialize JWK: %w", err)
	}
	pair.PrivateKey, err = json.Marshal(privKey)
	if err != nil {
		return KeyPair{}, fmt.Errorf("failed to serialize JWK: %w", err)
	}
	return pair, nil
}
