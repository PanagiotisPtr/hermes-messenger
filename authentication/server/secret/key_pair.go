package secret

import "crypto/rsa"

type KeyPair struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}
