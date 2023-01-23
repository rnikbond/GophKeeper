package cred_model

type Credential struct {
	MetaInfo string
	Login    []byte
	Password []byte
}
