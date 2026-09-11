package impl

import (
	"github.com/iFleksy/gopassxc/pkg/client"
)

type Request struct {
	Action        string `json:"action"`
	Message       string `json:"message,omitempty"`
	Nonce         string `json:"nonce"`
	ClientID      string `json:"clientID,omitempty"`
	TriggerUnlock string `json:"triggerUnlock"`
}

type RequestChangePublicKeys struct {
	Action    string `json:"action"`
	PublicKey string `json:"publicKey"`
	Nonce     string `json:"nonce"`
	ClientID  string `json:"clientID,omitempty"`
}

type EncryptedResponse struct {
	Action            string `json:"action"`
	Message           string `json:"message,omitempty"`
	Nonce             string `json:"nonce"`
	Error             string `json:"error,omitempty"`
	ErrorCode         string `json:"errorCode,omitempty"`
	DecryptedResponse []byte `json:"-"`
}

type ChangePubKeysResponse struct {
	Action    string `json:"action"`
	PublicKey string `json:"publicKey,omitempty"`
	Nonce     string `json:"nonce"`
}

type AssociateResponse struct {
	Hash    string `json:"hash"`
	ID      string `json:"id"`
	Nonce   string `json:"nonce"`
	Success string `json:"success"`
	Version string `json:"version"`
}

type EntriesResponse struct {
	Count   int            `json:"count"`
	Entries []client.Login `json:"entries"`
}

type AssociateRequest struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	IDKey  string `json:"idKey"`
}

type ActionRequest struct {
	Action string `json:"action"`
}

type DatabaseBaseBashResponse struct {
	Action  string `json:"action"`
	Hash    string `json:"hash"`
	Version string `json:"version"`
}

type GetLoginKeys struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type GetLoginRequest struct {
	Action string         `json:"action"`
	Url    string         `json:"url"`
	Keys   []GetLoginKeys `json:"keys"`
}

type GetLoginsResponse struct {
	Count   int            `json:"count"`
	Entries []client.Login `json:"entries"`
}

type TestAssociateRequest struct {
	Action        string `json:"action"`
	Id            string `json:"id"`
	Key           string `json:"key"`
	TriggerUnlock string `json:"triggerUnlock,omitempty"`
}

type TestAssociateResponse struct {
	Version string `json:"version"`
	Nonce   string `json:"nonce"`
	Hash    string `json:"hash"`
	ID      string `json:"id"`
	Success string `json:"success"`
}
