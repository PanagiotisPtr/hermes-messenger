package secret

import "crypto/rsa"

type Repository interface {
	GenerateRSAKeyPair() (KeyPair, error)
	StoreKeyPair(string, KeyPair) error
	GetPublicKey(string) (*rsa.PublicKey, error)
	GetPrivateKey(string) (*rsa.PrivateKey, error)
}
