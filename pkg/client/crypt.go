package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/kevinburke/nacl/scalarmult"
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
	PublicKey  nacl.Key
	PrivateKey nacl.Key

	AssociatedName string
	AssociatedKey  nacl.Key

	PeerKey nacl.Key
}

func (c *Crypt) NewKeys() {
	c.PrivateKey = nacl.NewKey()
	c.PublicKey = scalarmult.Base(c.PrivateKey)
}

func (c *Crypt) B64PublicKey() string {
	return NaclKeyToB64(c.PublicKey)
}

func (c *Crypt) B64PrivateKey() string {
	return NaclKeyToB64(c.PrivateKey)
}

func (c *Crypt) SetPeerKey(k string) {
	c.PeerKey = B64ToNaclKey(k)
}

func (c *Crypt) EncryptMessage(m Message) ([]byte, error) {
	if len(c.PeerKey) == 0 {
		return nil, errors.New("ErrInvalidPeerKey")
	}

	msgData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return box.EasySeal(msgData, c.PeerKey, c.PrivateKey), nil
}

func (c *Crypt) DecryptResponse(encryptedMsg []byte) ([]byte, error) {
	if len(c.PeerKey) == 0 {
		return nil, errors.New("ErrInvalidPeerKey")
	}

	return box.EasyOpen(encryptedMsg, c.PeerKey, c.PrivateKey)
}
