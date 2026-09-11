package impl

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
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

type EncryptedMessage struct {
	EncryptedData []byte `json:"encrypted_data"`
	Nonce         []byte `json:"nonce"`
}

type Crypt struct {
	publicKey     nacl.Key
	privateKey    nacl.Key
	peerPublicKey nacl.Key
}

func NewCrypto() Crypt {
	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	return Crypt{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (c *Crypt) PublicKey() string {
	return NaclKeyToB64(c.publicKey)
}

func (c *Crypt) PrivateKey() string {
	return NaclKeyToB64(c.privateKey)
}

func (c *Crypt) SetPeerKey(k string) {
	c.peerPublicKey = B64ToNaclKey(k)
}

func (c *Crypt) NewNonce() string {
	return base64.StdEncoding.EncodeToString((*nacl.NewNonce())[:])
}

func (c *Crypt) EncryptMessage(data []byte) (EncryptedMessage, error) {
	if len(c.peerPublicKey) == 0 {
		return EncryptedMessage{}, ErrInvalidPeerKey
	}

	encryptedData := box.EasySeal(data, c.peerPublicKey, c.privateKey)

	return EncryptedMessage{
		EncryptedData: encryptedData[nacl.NonceSize:],
		Nonce:         encryptedData[:nacl.NonceSize],
	}, nil
}

func (c *Crypt) DecryptMessage(encryptedData []byte) ([]byte, error) {
	if len(c.peerPublicKey) == 0 {
		return []byte{}, errors.New("ErrInvalidPeerKey")
	}

	return box.EasyOpen(encryptedData, c.peerPublicKey, c.privateKey)
}
