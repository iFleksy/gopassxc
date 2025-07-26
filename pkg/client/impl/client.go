package impl

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"

	"github.com/sirupsen/logrus"
)

type ClientOPTS struct {
	SocketPath        string
	AssociatedName    *string
	IdentificationKey *string
}

type Client struct {
	socket     net.Conn
	socketPath string
	crypt      Crypt

	// clientID - 24 bytes long random data, base64 encoded. This is used for a single session to identify different browsers if multiple are used with proxy application.
	ClientID string

	IdentificationKey string
	AssociatedName    string
}

const ClientID string = "gokeexc"

func New(opts ClientOPTS) Client {
	crypt := NewCrypto()
	client := Client{
		socketPath: opts.SocketPath,
		crypt:      crypt,
		ClientID:   ClientID + crypt.NewNonce(),
	}
	if opts.AssociatedName != nil {
		client.AssociatedName = *opts.AssociatedName
	}
	if opts.IdentificationKey != nil {
		client.IdentificationKey = *opts.IdentificationKey
	} else {
		client.IdentificationKey = crypt.NewNonce()
	}
	return client
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

func (c *Client) encryptMessage(data any) (string, string, error) {
	rawData, err := json.Marshal(data)
	if err != nil {
		return "", "", err
	}

	logrus.Debugf("[ RAW MESSAGE ]: %s", string(rawData))
	nonce, encrypted, err := c.crypt.EncryptMessage(rawData)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(nonce), base64.StdEncoding.EncodeToString(encrypted), nil
}

func (c *Client) decryptMessage(nonce string, data string) ([]byte, error) {
	encryptedData := nonce + data
	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)

	if err != nil {
		return []byte{}, err
	}

	return c.crypt.DecryptMessage(decodedData)
}

func (c *Client) sendEncryptedMessage(action string, data any, obj any) error {
	var response EncryptedResponse
	logrus.Debugf("Action: %s\nsend message %#v \n", action, data)
	nonce, encryptedMessage, err := c.encryptMessage(data)

	if err != nil {
		return err
	}

	req := Request{
		Action:   action,
		Message:  encryptedMessage,
		Nonce:    nonce,
		ClientID: c.ClientID,
	}

	logrus.Debugf("send request: %s", req)
	rawResponse, err := c.sendMessage(req)
	logrus.Debugf("%s\n", rawResponse)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(rawResponse, &response); err != nil {
		return err
	}

	if response.Error != "" {
		return errors.New(response.Error)
	}

	logrus.Debugf("%#v\n", response)

	decryptedMsg, err := c.decryptMessage(response.Nonce, response.Message)

	if err != nil {
		return err
	}

	logrus.Debugf("raw encoded response: %s", string(decryptedMsg))
	response.DecryptedResponse = decryptedMsg
	if obj != nil {
		err = json.Unmarshal(response.DecryptedResponse, obj)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) sendMessage(req any) ([]byte, error) {
	content, err := json.Marshal(req)
	if err != nil {
		return []byte{}, err
	}
	_, err = c.socket.Write(content)
	if err != nil {
		return []byte{}, err
	}

	buff := make([]byte, 4096)
	count, err := c.socket.Read(buff)
	if err != nil {
		return []byte{}, err
	}
	return buff[0:count], nil
}

func (c *Client) ChangePublicKeys() (string, error) {
	message := RequestChangePublicKeys{
		Action:    "change-public-keys",
		PublicKey: c.crypt.PublicKey(),
		Nonce:     c.crypt.NewNonce(),
		ClientID:  c.ClientID,
	}

	rawResponse, err := c.sendMessage(message)
	if err != nil {
		return "", err
	}

	resp := ChangePubKeysResponse{}
	err = json.Unmarshal(rawResponse, &resp)
	if err != nil {
		return "", err
	}

	if resp.PublicKey == "" {
		return "", errors.New("change-public-keys failed")
	}

	c.crypt.SetPeerKey(resp.PublicKey)
	return resp.PublicKey, nil
}

func (c *Client) GetDBHash() (string, error) {
	const action = "get-databasehash"

	req := ActionRequest{
		Action: action,
	}

	var data DatabaseBaseBashResponse
	err := c.sendEncryptedMessage(action, req, &data)
	if err != nil {
		return "", err
	}
	return data.Hash, nil
}

func (c *Client) Associate() (string, string, error) {
	const action = "associate"

	msg := AssociateRequest{
		Action: action,
		Key:    c.crypt.PublicKey(),
		IDKey:  c.IdentificationKey,
	}

	data := AssociateResponse{}

	err := c.sendEncryptedMessage(action, msg, &data)

	if err != nil {
		return "", "", err
	}

	c.AssociatedName = data.ID
	return c.AssociatedName, c.IdentificationKey, err
}

func (c *Client) GetLogins(url string) error {
	const action = "get-logins"
	req := GetLoginRequest{
		Action: action,
		Url:    url,
		Keys: []GetLoginKeys{
			{
				Id:  c.AssociatedName,
				Key: c.IdentificationKey,
			},
		},
	}

	var response GetLoginsResponse
	err := c.sendEncryptedMessage(action, req, &response)
	return err

}

func (c *Client) TestAssociate() error {
	const action = "test-associate"
	req := TestAssociateRequest{
		Action: action,
		Id:     c.AssociatedName,
		Key:    c.IdentificationKey,
	}
	var response TestAssociateResponse
	err := c.sendEncryptedMessage(action, req, &response)
	return err
}

func (c *Client) UnlockDatabase() error {
	const action = "database-unlocked"
	req := ActionRequest{
		Action: action,
	}

	err := c.sendEncryptedMessage(action, req, nil)
	if err != nil {
		return err
	}
	return nil
}

// func (c *Client) TestAssociate() error {
// 	msg := Message{
// 		Action: ActionTestAssociate,
// 		Key:    NaclKeyToB64(c.crypt.AssociatedKey),
// 		ID:     c.crypt.AssociatedName,
// 	}

// 	_, err := c.sendEncryptedMessage(msg)
// 	return err
// }

// func (c *Client) GetLogins(url string) ([]*Entry, error) {
// 	msg := Message{
// 		Action: ActionGetLogins,
// 		URL:    url,
// 		Keys: []*MessageKeys{
// 			{
// 				ID:  c.crypt.AssociatedName,
// 				Key: NaclKeyToB64(c.crypt.AssociatedKey),
// 			},
// 		},
// 	}

// 	response, err := c.sendEncryptedMessage(msg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data EntriesResponse

// 	err = json.Unmarshal(response.DecryptedResponse, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return data.Entries, nil
// }
