package credentials

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
}

const kidRandomChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewSignKeyPair() KeyPair {
	return NewKeyPair(
		jwk.KeyOperationList{jwk.KeyOpVerify},
		jwk.KeyOperationList{jwk.KeyOpSign},
	)
}

func NewKeyPair(pubOps jwk.KeyOperationList, privOps jwk.KeyOperationList) KeyPair {
	alg := "EdDSA"
	kid := "key-" + SecureRandomString(challengeBytes, kidRandomChars)

	rawPubKey, rawPrivKey, err := ed25519.GenerateKey(rand.Reader)
	neverFail(err)

	privKey, err := jwk.FromRaw(rawPrivKey)
	neverFail(err)

	pubKey, err := jwk.FromRaw(rawPubKey)
	neverFail(err)

	for k, v := range map[string]interface{}{
		jwk.KeyOpsKey:    pubOps,
		jwk.AlgorithmKey: alg,
		jwk.KeyIDKey:     kid,
	} {
		err = pubKey.Set(k, v)
		neverFail(err)
	}

	for k, v := range map[string]interface{}{
		jwk.KeyOpsKey:    privOps,
		jwk.AlgorithmKey: alg,
		jwk.KeyIDKey:     kid,
	} {
		err = privKey.Set(k, v)
		neverFail(err)
	}

	var pair KeyPair
	pair.PublicKey, err = json.Marshal(pubKey)
	neverFail(err)
	pair.PrivateKey, err = json.Marshal(privKey)
	neverFail(err)
	return pair
}
