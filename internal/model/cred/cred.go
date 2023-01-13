package cred

// CredentialFull - Данные получения логинов и паролей
type CredentialFull struct {
	// Email - электронная почта
	Email string
	// MetaInfo - Метаинформация для хранимых данных
	MetaInfo string
	// Password - Пароль
	Password string
}

// CredentialGet - Данные получения логинов и паролей
type CredentialGet struct {
	// Email - электронная почта
	Email string
	// MetaInfo - Метаинформация для хранимых данных
	MetaInfo string
}
