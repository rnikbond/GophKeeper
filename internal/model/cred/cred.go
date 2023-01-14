package cred

// CredentialFull - Данные получения логинов и паролей
type CredentialFull struct {
	// MetaInfo - Метаинформация для хранимых данных
	MetaInfo string
	// Email - электронная почта
	Email string
	// Password - Пароль
	Password string
}

// CredentialGet - Данные получения логинов и паролей
type CredentialGet struct {
	// MetaInfo - Метаинформация для хранимых данных
	MetaInfo string
}
