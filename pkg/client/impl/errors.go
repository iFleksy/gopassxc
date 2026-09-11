package impl

import "github.com/iFleksy/gopassxc/pkg/client"

func HandleErrorCode(code int) error {
	switch code {
	case 1:
		return client.ErrKeepassDatabaseNotOpened
	case 2:
		return client.ErrKeepassDatabaseHashNotReceived
	case 3:
		return client.ErrKeepassClientPublicKeyNotReceived
	case 4:
		return client.ErrKeepassCannotDecryptMessage
	case 5:
		return client.ErrKeepassTimeoutOrNotConnected
	case 6:
		return client.ErrKeepassActionCancelledOrDenied
	case 7:
		return client.ErrKeepassCannotEncryptMessage
	case 8:
		return client.ErrKeepassAssociationFailed
	case 9:
		return client.ErrKeepassKeyChangeFailed
	case 10:
		return client.ErrKeepassEncryptionKeyUnrecognized
	case 11:
		return client.ErrKeepassNoSavedDatabasesFound
	case 12:
		return client.ErrKeepassIncorrectAction
	case 13:
		return client.ErrKeepassEmptyMessageReceived
	case 14:
		return client.ErrKeepassNoURLProvided
	case 15:
		return client.ErrKeepassNoLoginsFound
	case 16:
		return client.ErrKeepassNoGroupsFound
	case 17:
		return client.ErrKeepassCannotCreateNewGroup
	case 18:
		return client.ErrKeepassNoValidUUIDProvided
	case 19:
		return client.ErrKeepassAccessToAllEntriesDenied
	case 20:
		return client.ErrPasskeysAttestationNotSupported
	case 21:
		return client.ErrPasskeysCredentialIsExcluded
	case 22:
		return client.ErrPasskeysRequestCanceled
	case 23:
		return client.ErrPasskeysInvalidUserVerification
	case 24:
		return client.ErrPasskeysEmptyPublicKey
	case 25:
		return client.ErrPasskeysInvalidURLProvided
	case 26:
		return client.ErrPasskeysOriginNotAllowed
	case 27:
		return client.ErrPasskeysDomainIsNotValid
	case 28:
		return client.ErrPasskeysDomainRPIDMismatch
	case 29:
		return client.ErrPasskeysNoSupportedAlgorithms
	case 30:
		return client.ErrPasskeysWaitForLifetimer
	case 31:
		return client.ErrPasskeysUnknownError
	case 32:
		return client.ErrPasskeysInvalidChallenge
	case 33:
		return client.ErrPasskeysInvalidUserID
	default:
		return client.ErrUnknownErrorCode
	}
}
