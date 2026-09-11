package impl

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/iFleksy/gopassxc/pkg/client"
)

func (c *Client) ChangePublicKeys(ctx context.Context) error {
	message := RequestChangePublicKeys{
		Action:    "change-public-keys",
		PublicKey: c.crypt.PublicKey(),
		Nonce:     c.crypt.NewNonce(),
		ClientID:  c.ClientID,
	}

	rawResponse, err := c.sendMessage(ctx, message)
	if err != nil {
		return err
	}

	resp := ChangePubKeysResponse{}
	err = json.Unmarshal(rawResponse, &resp)
	if err != nil {
		return err
	}

	if resp.PublicKey == "" {
		return errors.New("change-public-keys failed")
	}

	c.crypt.SetPeerKey(resp.PublicKey)
	return nil
}

func (c *Client) GetDBHash(ctx context.Context) (string, error) {
	const action = "get-databasehash"

	req := ActionRequest{
		Action: action,
	}

	var data DatabaseBaseBashResponse
	err := c.sendEncryptedMessage(ctx, action, req, &data)
	if err != nil {
		return "", err
	}
	return data.Hash, nil
}

func (c *Client) Associate(ctx context.Context) (client.AssociateResponse, error) {
	const action = "associate"
	msg := AssociateRequest{
		Action: action,
		Key:    c.crypt.PublicKey(),
		IDKey:  c.IdentificationKey,
	}

	data := AssociateResponse{}

	err := c.sendEncryptedMessage(ctx, action, msg, &data)

	if err != nil {
		return client.AssociateResponse{}, err
	}

	return client.AssociateResponse{
		AssociatedName:    data.ID,
		IdentificationKey: c.IdentificationKey,
	}, nil
}

func (c *Client) GetLogins(ctx context.Context, url string) ([]client.Login, error) {
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
	err := c.sendEncryptedMessage(ctx, action, req, &response)
	if err != nil {
		return nil, err
	}
	return response.Entries, nil
}

func (c *Client) TestAssociate(ctx context.Context, triggerUnlock bool) error {
	const action = "test-associate"
	req := TestAssociateRequest{
		Action:        action,
		Id:            c.AssociatedName,
		Key:           c.IdentificationKey,
		TriggerUnlock: BoolToString(triggerUnlock),
	}
	var response TestAssociateResponse
	err := c.sendEncryptedMessage(ctx, action, req, &response)
	return err
}

func (c *Client) UnlockDatabase(ctx context.Context) error {
	const action = "database-unlocked"
	req := ActionRequest{
		Action: action,
	}

	err := c.sendEncryptedMessage(ctx, action, req, nil)
	if err != nil {
		return err
	}
	return nil
}
