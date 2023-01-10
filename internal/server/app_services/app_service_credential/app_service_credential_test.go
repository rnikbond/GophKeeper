package app_service_credential

import (
	"testing"

	"github.com/stretchr/testify/require"

	"GophKeeper/internal/model/cred"
	"GophKeeper/internal/storage/credential_store"
	"GophKeeper/pkg/errs"
)

func TestCredentialAppService(t *testing.T) {

	store := credential_store.NewMemoryStorage()
	serv := NewCredentialAppService(store)

	testDataOK := cred.CredentialFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Password: "qwerty",
	}

	testDataChange := cred.CredentialFull{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
		Password: "qwerty123",
	}

	testDataGet := cred.CredentialGet{
		Email:    "test@email.com",
		MetaInfo: "www.ololo.com",
	}

	testDataFail := cred.CredentialGet{
		Email:    "test@email.com",
		MetaInfo: "www.test.com",
	}

	errCreate := serv.Create(testDataOK)
	require.NoError(t, errCreate)

	data, errGet := serv.Get(testDataGet)
	require.NoError(t, errGet)
	require.Equal(t, data, testDataOK)

	_, errGet = serv.Get(testDataFail)
	require.Error(t, errGet, errs.ErrNotFound)

	errChange := serv.Change(testDataChange)
	require.NoError(t, errChange)

	data, errGet = serv.Get(testDataGet)
	require.NoError(t, errGet)
	require.Equal(t, data, testDataChange)

	errDel := serv.Delete(testDataGet)
	require.NoError(t, errDel)

	_, errGet = serv.Get(testDataGet)
	require.Error(t, errGet, errs.ErrNotFound)

	errChange = serv.Change(testDataChange)
	require.Error(t, errGet, errs.ErrNotFound)

	errCreate = serv.Create(testDataOK)
	require.NoError(t, errCreate)

	errCreate = serv.Create(testDataOK)
	require.Error(t, errCreate, errs.ErrAlreadyExist)
}
