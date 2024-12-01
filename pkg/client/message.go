package client

type Action string

const (
	// Аутентификация и ассоциация
	ActionAssociate        Action = "associate"          // Ассоциация с браузерным расширением
	ActionTestAssociate    Action = "test-associate"     // Проверка статуса ассоциации
	ActionChangePublicKeys Action = "change-public-keys" // Обмен публичными ключами

	// Операции с базой данных
	ActionGetDatabaseHash    Action = "get-databasehash"     // Получение хеша базы данных
	ActionLockDatabase       Action = "lock-database"        // Блокировка базы данных
	ActionGetDatabaseGroups  Action = "get-database-groups"  // Получение списка групп
	ActionCreateNewGroup     Action = "create-new-group"     // Создание новой группы
	ActionGetDatabaseEntries Action = "get-database-entries" // Получение всех записей

	// Управление логинами и паролями
	ActionGetLogins        Action = "get-logins"        // Получение учетных данных
	ActionSetLogin         Action = "set-login"         // Сохранение/обновление учетных данных
	ActionDeleteEntry      Action = "delete-entry"      // Удаление записи
	ActionGeneratePassword Action = "generate-password" // Генерация нового пароля
	ActionGetTotp          Action = "get-totp"          // Получение TOTP кода

	// Автоматизация
	ActionRequestAutotype Action = "request-autotype" // Запуск автоввода

	// Операции с Passkeys (если скомпилировано с WITH_XC_BROWSER_PASSKEYS)
	ActionPasskeysGet      Action = "passkeys-get"      // Получение учетных данных passkey
	ActionPasskeysRegister Action = "passkeys-register" // Регистрация нового passkey
)

type ErrorCode int

const (
	ErrorKeePassIncorrectAction ErrorCode = iota
	ErrorKeePassEmptyMessageReceived
	ErrorKeePassClientPublicKeyNotReceived
	ErrorKeePassCannotDecryptMessage
	ErrorKeePassDatabaseNotOpened
	ErrorKeePassDatabaseHashNotReceived
	ErrorKeePassAssociationFailed
	ErrorKeePassActionCancelledOrDenied
	ErrorKeePassEncryptionKeyUnrecognized
	ErrorKeePassNoUrlProvided
	ErrorKeePassNoLoginsFound
	ErrorKeePassNoGroupsFound
	ErrorKeePassCannotCreateNewGroup
	ErrorKeePassNoValidUuidProvided
	ErrorKeePassAccessToAllEntriesDenied
	ErrorPasskeysEmptyPublicKey
	ErrorPasskeysInvalidUrlProvided
)

type MessageKeys struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type Message struct {
	Action    Action `json:"action"`
	PublicKey string `json:"publicKey,omitempty"`
	ClientID  string `json:"clientID"`
	Nonce     string `json:"nonce,omitempty"`

	ID string `json:"id,omitempty"`

	IDKey string `json:"idKey,omitempty"`
	Key   string `json:"key,omitempty"`

	Message string `json:"message,omitempty"`

	Keys []*MessageKeys `json:"keys,omitempty"`

	URL string `json:"url,omitempty"`
}
