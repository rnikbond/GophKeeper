package binary

type DataFull struct {
	Email    string
	MetaInfo string
	Bytes    []byte
}

type DataGet struct {
	Email    string
	MetaInfo string
}
