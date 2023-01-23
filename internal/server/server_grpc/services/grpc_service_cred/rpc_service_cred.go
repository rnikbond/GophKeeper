//go:generate mockgen -source rpc_service_cred.go -destination mocks/rpc_service_cred_mock.go -package grpc_service_cred
package grpc_service_cred

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/server/model/cred"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/proto/credential"
)

type CredentialApp interface {
	Create(in cred.CredentialFull) error
	Get(in cred.CredentialGet) (cred.CredentialFull, error)
	Delete(in cred.CredentialGet) error
	Change(in cred.CredentialFull) error
}

type CredServiceRPC struct {
	credential.CredentialServiceServer

	credApp CredentialApp
	logger  *zap.Logger
}

// NewCredServiceRPC - Создание эклемпляра gRPC сервиса дял хранения данных в виде логина и пароля.
func NewCredServiceRPC(credApp CredentialApp) *CredServiceRPC {
	serv := &CredServiceRPC{
		credApp: credApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *CredServiceRPC) Create(ctx context.Context, in *credential.CreateRequest) (*credential.Empty, error) {

	data := cred.CredentialFull{
		MetaInfo: in.MetaInfo,
		Email:    string(in.Email),
		Password: string(in.Password),
	}

	err := serv.credApp.Create(data)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExist) {
			return &credential.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())
		}

		serv.logger.Error("failed create credential data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo),
			zap.String("email", string(in.Email)),
			zap.String("pwd", string(in.Password)))

		return &credential.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &credential.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *CredServiceRPC) Change(ctx context.Context, in *credential.ChangeRequest) (*credential.Empty, error) {

	data := cred.CredentialFull{
		MetaInfo: in.MetaInfo,
		Email:    string(in.Email),
		Password: string(in.Password),
	}

	err := serv.credApp.Change(data)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &credential.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed change credential data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo),
			zap.String("email", string(in.Email)),
			zap.String("pwd", string(in.Password)))

		return &credential.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &credential.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *CredServiceRPC) Delete(ctx context.Context, in *credential.DeleteRequest) (*credential.Empty, error) {

	data := cred.CredentialGet{
		MetaInfo: in.MetaInfo,
	}

	err := serv.credApp.Delete(data)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &credential.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete credential data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &credential.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &credential.Empty{}, nil
}

// Get - Получение данных по email и метаданным.
func (serv *CredServiceRPC) Get(ctx context.Context, in *credential.GetRequest) (*credential.GetResponse, error) {

	inData := cred.CredentialGet{
		MetaInfo: in.MetaInfo,
	}

	data, err := serv.credApp.Get(inData)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &credential.GetResponse{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed get credential data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &credential.GetResponse{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	out := &credential.GetResponse{
		Email:    []byte(data.Email),
		Password: []byte(data.Password),
	}

	return out, nil
}
