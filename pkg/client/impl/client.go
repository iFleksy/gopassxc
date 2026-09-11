package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

const IO_BUFFER_SIZE = 4096

type ClientOPTS struct {
	SocketPath        string
	AssociatedName    *string
	IdentificationKey *string
	TriggerUnlock     *bool
}

type Client struct {
	socket     net.Conn
	socketPath string
	crypt      Crypt

	socketMutex sync.Mutex

	// clientID - 24 bytes long random data, base64 encoded. This is used for a single session to identify different browsers if multiple are used with proxy application.
	ClientID string

	IdentificationKey string
	AssociatedName    string
	TriggerUnlock     bool
}

const ClientID string = "gokeexc"

func New(opts ClientOPTS) (*Client, error) {
	crypt := NewCrypto()
	client := &Client{
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

	if opts.TriggerUnlock != nil {
		client.TriggerUnlock = *opts.TriggerUnlock
	} else {
		client.TriggerUnlock = false
	}

	return client, nil
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
	encryptedData, err := c.crypt.EncryptMessage(rawData)
	if err != nil {
		return "", "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData.Nonce), base64.StdEncoding.EncodeToString(encryptedData.EncryptedData), nil
}

func (c *Client) decryptMessage(nonce string, data string) ([]byte, error) {
	encryptedData := nonce + data
	decodedData, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return []byte{}, err
	}

	return c.crypt.DecryptMessage(decodedData)
}

func (c *Client) getTriggerUnlock() string {
	if c.TriggerUnlock {
		return "true"
	}
	return "false"
}

func (c *Client) sendEncryptedMessage(ctx context.Context, action string, data any, obj any) error {
	var response EncryptedResponse
	nonce, encryptedMessage, err := c.encryptMessage(data)

	if err != nil {
		return err
	}

	req := Request{
		Action:        action,
		Message:       encryptedMessage,
		Nonce:         nonce,
		ClientID:      c.ClientID,
		TriggerUnlock: c.getTriggerUnlock(),
	}

	logrus.Debugf("send request: %#v", req)

	rawResponse, err := c.sendMessage(ctx, req)
	logrus.Debugf("%s\n", rawResponse)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(rawResponse, &response); err != nil {
		return err
	}

	if response.ErrorCode != "" {
		errorCode, err := strconv.Atoi(response.ErrorCode)
		if err != nil {
			return errors.New(response.Error)
		}
		logrus.Debugf("error code: %d", errorCode)
		err = HandleErrorCode(errorCode)
		if err != nil {
			return err
		}
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

func (c *Client) sendMessage(ctx context.Context, req any) ([]byte, error) {
	// Блокируем доступ к сокету для предотвращения race conditions
	c.socketMutex.Lock()
	defer c.socketMutex.Unlock()

	content, err := json.Marshal(req)
	if err != nil {
		return []byte{}, err
	}

	// Создаем каналы для результатов
	writeErrChan := make(chan error, 1)
	type readResult struct {
		data []byte
		err  error
	}
	readChan := make(chan readResult, 1)

	go func() {
		_, err := c.socket.Write(content)
		writeErrChan <- err
	}()

	select {
	case <-ctx.Done():
		return []byte{}, ctx.Err()
	case err := <-writeErrChan:
		if err != nil {
			return []byte{}, err
		}
	}

	go func() {
		buff := make([]byte, IO_BUFFER_SIZE)
		count, err := c.socket.Read(buff)
		if err != nil {
			readChan <- readResult{nil, err}
		} else {
			readChan <- readResult{buff[0:count], nil}
		}
	}()

	// Ждем завершения чтения или timeout
	select {
	case <-ctx.Done():
		return []byte{}, ctx.Err()
	case result := <-readChan:
		return result.data, result.err
	}
}
