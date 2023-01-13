package card

// DataCardFull - Данные карты и их принаждлежности
type DataCardFull struct {
	// MetaInfo - Метаинформация для хранимых данных
	MetaInfo string
	// Number - Номер карты
	Number string
	// Period - Срок действия карты
	Period string
	// CVV - Секретый код
	CVV string
	// FullName - Полное имя держателя карты
	FullName string
}

// DataCardGet - Данные lkz получения
type DataCardGet struct {
	// MetaInfo - Метаинформация для хранимых данных
	MetaInfo string
}
