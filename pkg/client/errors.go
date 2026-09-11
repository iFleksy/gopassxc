package client

import "errors"

var (
	ErrKeepassDatabaseNotOpened          = errors.New("database not opened")
	ErrKeepassDatabaseHashNotReceived    = errors.New("database hash not received")
	ErrKeepassClientPublicKeyNotReceived = errors.New("client public key not received")
	ErrKeepassCannotDecryptMessage       = errors.New("cannot decrypt message")
	ErrKeepassTimeoutOrNotConnected      = errors.New("timeout or not connected")
	ErrKeepassActionCancelledOrDenied    = errors.New("action cancelled or denied")
	ErrKeepassCannotEncryptMessage       = errors.New("cannot encrypt message")
	ErrKeepassAssociationFailed          = errors.New("association failed")
	ErrKeepassKeyChangeFailed            = errors.New("key change failed")
	ErrKeepassEncryptionKeyUnrecognized  = errors.New("encryption key unrecognized")
	ErrKeepassNoSavedDatabasesFound      = errors.New("no saved databases found")
	ErrKeepassIncorrectAction            = errors.New("incorrect action")
	ErrKeepassEmptyMessageReceived       = errors.New("empty message received")
	ErrKeepassNoURLProvided              = errors.New("no URL provided")
	ErrKeepassNoLoginsFound              = errors.New("no logins found")
	ErrKeepassNoGroupsFound              = errors.New("no groups found")
	ErrKeepassCannotCreateNewGroup       = errors.New("cannot create new group")
	ErrKeepassNoValidUUIDProvided        = errors.New("no valid UUID provided")
	ErrKeepassAccessToAllEntriesDenied   = errors.New("access to all entries denied")

	ErrPasskeysAttestationNotSupported = errors.New("attestation not supported")
	ErrPasskeysCredentialIsExcluded    = errors.New("credential is excluded")
	ErrPasskeysRequestCanceled         = errors.New("request canceled")
	ErrPasskeysInvalidUserVerification = errors.New("invalid user verification")
	ErrPasskeysEmptyPublicKey          = errors.New("empty public key")
	ErrPasskeysInvalidURLProvided      = errors.New("invalid URL provided")
	ErrPasskeysOriginNotAllowed        = errors.New("origin not allowed")
	ErrPasskeysDomainIsNotValid        = errors.New("domain is not valid")
	ErrPasskeysDomainRPIDMismatch      = errors.New("domain RPID mismatch")
	ErrPasskeysNoSupportedAlgorithms   = errors.New("no supported algorithms")
	ErrPasskeysWaitForLifetimer        = errors.New("wait for lifetimer")
	ErrPasskeysUnknownError            = errors.New("unknown error")
	ErrPasskeysInvalidChallenge        = errors.New("invalid challenge")
	ErrPasskeysInvalidUserID           = errors.New("invalid user ID")

	ErrUnknownErrorCode = errors.New("unknown error code")
)

func IsAssociationError(err error) bool {
	return err == ErrKeepassAssociationFailed
}
