package client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"

	"github.com/kevinburke/nacl"
	"github.com/sirupsen/logrus"
)

type ServerResponse struct {
	Action            Action          `json:"action"`
	Message           string          `json:"message,omitempty"`
	Nonce             string          `json:"nonce"`
	PublicKey         string          `json:"publicKey,omitempty"`
	Error             string          `json:"error,omitempty"`
	ErrorCode         json.RawMessage `json:"errorCode,omitempty"`
	DecryptedResponse []byte          `json:"-"`
}

type AssociateResponse struct {
	Hash    string `json:"hash"`
	ID      string `json:"id"`
	Nonce   string `json:"nonce"`
	Success string `json:"success"`
	Version string `json:"version"`
}

type Entry struct {
	Group    string `json:"Group"`
	Login    string `json:"Login"`
	Name     string `json:"Name"`
	Password string `json:"password"`
	UUID     string `json"uuid"`
	TOTP     string `json:"totp,omitempty"`
}

type EntriesResponse struct {
	Count   int      `json"count"`
	Entries []*Entry `json"entries"`
}

type Client struct {
	socket     net.Conn
	socketPath string
	crypt      Crypt
	ClientID   string
}

const ClientID string = "gokeexc"

func NewClient(sockPath string, assoName string, assoB64Key string) Client {
	assoKey := nacl.NewKey()
	if len(assoB64Key) != 0 {
		assoKey = B64ToNaclKey(assoB64Key)
	}

	crypt := Crypt{
		AssociatedName: assoName,
		AssociatedKey:  assoKey,
	}

	crypt.NewKeys()
	return Client{
		socketPath: sockPath,
		crypt:      crypt,
		ClientID:   ClientID + NaclNonceToB64(nacl.NewNonce()),
	}
}

func (c *Client) Connect() error {
	var err error
	logrus.Debugf("connect to socket %s", c.socketPath)
	c.socket, err = net.DialUnix("unix", nil, &net.UnixAddr{Name: c.socketPath, Net: "unix"})
	return err
}

func (c *Client) Disconnect() error {
	if c.socket != nil {
		return c.socket.Close()
	}
	return nil
}

func (c *Client) sendEncryptedMessage(msg Message) (ServerResponse, error) {
	encryptedMsg, err := c.crypt.EncryptMessage(msg)
	var response ServerResponse
	if err != nil {
		return response, err
	}

	rawRequest, err := json.Marshal(msg)
	if err != nil {
		return response, err
	}
	logrus.Debugf("send request: %s", string(rawRequest))

	msg = Message{
		Action:  msg.Action,
		Message: base64.StdEncoding.EncodeToString(encryptedMsg[nacl.NonceSize:]),
		Nonce:   base64.StdEncoding.EncodeToString(encryptedMsg[:nacl.NonceSize]),
	}

	response, err = c.sendMessage(msg)
	if err != nil {
		return response, err
	}

	if response.Error != "" {
		return response, errors.New(response.Error)
	}

	decoded, err := base64.StdEncoding.DecodeString(response.Nonce + response.Message)
	if err != nil {
		return response, err
	}

	decryptedMsg, err := c.crypt.DecryptResponse(decoded)
	if err != nil {
		return response, err
	}

	logrus.Debugf("raw encoded response: %s", string(decryptedMsg))
	response.DecryptedResponse = decryptedMsg
	return response, nil
}

func (c *Client) sendMessage(msg Message) (ServerResponse, error) {
	var response ServerResponse
	if msg.Nonce == "" {
		msg.Nonce = NaclNonceToB64(nacl.NewNonce())
	}

	msg.ClientID = c.ClientID

	content, err := json.Marshal(msg)
	if err != nil {
		return response, err
	}

	logrus.Debugf("send message %s", string(content))
	_, err = c.socket.Write(content)
	if err != nil {
		return response, err
	}

	buff := make([]byte, 4096)
	count, err := c.socket.Read(buff)
	if err != nil {
		return response, err
	}
	buff = buff[0:count]
	logrus.Debugf("raw response: %s", string(buff))
	err = json.Unmarshal(buff, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (c *Client) ChangePublicKeys() error {
	message := Message{
		Action:    ActionChangePublicKeys,
		PublicKey: c.crypt.B64PublicKey(),
	}

	resp, err := c.sendMessage(message)
	if err != nil {
		return err
	}

	if resp.PublicKey != "" {
		c.crypt.SetPeerKey(resp.PublicKey)
		return nil
	}

	return errors.New("change-public-keys failed")
}

func (c *Client) Associate() (AssociateResponse, error) {
	msg := Message{
		Action: ActionAssociate,
		Key:    NaclKeyToB64(c.crypt.PublicKey),
		IDKey:  NaclKeyToB64(c.crypt.AssociatedKey),
	}

	var data AssociateResponse

	response, err := c.sendEncryptedMessage(msg)

	if err != nil {
		return data, err
	}

	err = json.Unmarshal(response.DecryptedResponse, &data)
	if err != nil {
		return data, err
	}

	c.crypt.AssociatedName = data.ID

	return data, err
}

func (c *Client) TestAssociate() error {
	msg := Message{
		Action: ActionTestAssociate,
		Key:    NaclKeyToB64(c.crypt.AssociatedKey),
		ID:     c.crypt.AssociatedName,
	}

	_, err := c.sendEncryptedMessage(msg)
	return err
}

func (c *Client) GetAssociatedProfile() (string, string) {
	return c.crypt.AssociatedName, NaclKeyToB64(c.crypt.AssociatedKey)
}

func (c *Client) GetLogins(url string) ([]*Entry, error) {
	msg := Message{
		Action: ActionGetLogins,
		URL:    url,
		Keys: []*MessageKeys{
			{
				ID:  c.crypt.AssociatedName,
				Key: NaclKeyToB64(c.crypt.AssociatedKey),
			},
		},
	}

	response, err := c.sendEncryptedMessage(msg)
	if err != nil {
		return nil, err
	}

	var data EntriesResponse

	err = json.Unmarshal(response.DecryptedResponse, &data)
	if err != nil {
		return nil, err
	}

	return data.Entries, nil
}

func (c *Client) GeneratePassword() error {
	msg := Message{
		Action: ActionGeneratePassword,
		Keys: []*MessageKeys{
			{
				ID:  c.crypt.AssociatedName,
				Key: NaclKeyToB64(c.crypt.AssociatedKey),
			},
		},
	}
	_, err := c.sendEncryptedMessage(msg)
	return err
}
