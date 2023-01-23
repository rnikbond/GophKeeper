package card_model

type Card struct {
	// MetaInfo - Метаинформация для хранимых данных
	MetaInfo string
	// Number - Номер карты
	Number []byte
	// Period - Срок действия карты
	Period []byte
	// CVV - Секретый код
	CVV []byte
	// FullName - Полное имя держателя карты
	FullName []byte
}
