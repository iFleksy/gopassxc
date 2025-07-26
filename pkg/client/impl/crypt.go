package impl

import (
	"encoding/base64"
	"errors"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/kevinburke/nacl/scalarmult"
)

var (
	ErrInvalidPeerKey = errors.New("Invalid Peer Key")
)

func NaclKeyToB64(key nacl.Key) string {
	return base64.StdEncoding.EncodeToString((*key)[:])
}

func NaclNonceToB64(nonce nacl.Nonce) string {
	return base64.StdEncoding.EncodeToString((*nonce)[:])
}

func B64ToNaclNonce(b64Nonce string) nacl.Nonce {
	decoded, err := base64.StdEncoding.DecodeString(b64Nonce)
	if err != nil {
		panic(err)
	}
	nonce := new([nacl.NonceSize]byte)
	copy(nonce[:], decoded)
	return nonce
}

func B64ToNaclKey(b64Key string) nacl.Key {
	decoded, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		panic(err)
	}
	key := new([nacl.KeySize]byte)
	copy(key[:], decoded)
	return key
}

type Crypt struct {
	publicKey  nacl.Key
	privateKey nacl.Key
	peerKey    nacl.Key
}

func NewCrypto() Crypt {
	key := nacl.NewKey()
	pubKey := scalarmult.Base(key)

	return Crypt{
		publicKey:  pubKey,
		privateKey: key,
	}
}

func (c *Crypt) PublicKey() string {
	return NaclKeyToB64(c.publicKey)
}

func (c *Crypt) PrivateKey() string {
	return NaclKeyToB64(c.privateKey)
}

func (c *Crypt) SetPeerKey(k string) {
	c.peerKey = B64ToNaclKey(k)
}

func (c *Crypt) NewNonce() string {
	return base64.StdEncoding.EncodeToString((*nacl.NewNonce())[:])
}

func (c *Crypt) EncryptMessage(data []byte) ([]byte, []byte, error) {
	if len(c.peerKey) == 0 {
		return []byte{}, []byte{}, ErrInvalidPeerKey
	}

	encryptedData := box.EasySeal(data, c.peerKey, c.privateKey)

	return encryptedData[:nacl.NonceSize], encryptedData[nacl.NonceSize:], nil
}

func (c *Crypt) DecryptMessage(encryptedData []byte) ([]byte, error) {
	if len(c.peerKey) == 0 {
		return []byte{}, errors.New("ErrInvalidPeerKey")
	}

	return box.EasyOpen(encryptedData, c.peerKey, c.privateKey)
}
