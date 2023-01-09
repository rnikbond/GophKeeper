package rpc_services

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/cred"
	"GophKeeper/internal/server/app_services"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/data/credential"
)

type CredServiceRPC struct {
	pb.CredentialServiceServer

	credApp app_services.CredentialApp
	logger  *zap.Logger
}

// NewCredServiceRPC - Создание эклемпляра gRPC сервиса дял хранения данных в виде логина и пароля.
func NewCredServiceRPC(credApp app_services.CredentialApp) *CredServiceRPC {
	serv := &CredServiceRPC{
		credApp: credApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *CredServiceRPC) Create(ctx context.Context, in *pb.CreateRequest) (*pb.Empty, error) {

	data := cred.CredentialFull{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
		Password: in.Password,
	}

	err := serv.credApp.Create(data)
	if err != nil {
		if err == errs.ErrAlreadyExist {
			return &pb.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())
		}

		serv.logger.Error("failed create credential data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo),
			zap.String("pwd", in.Password))

		return &pb.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &pb.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *CredServiceRPC) Change(ctx context.Context, in *pb.ChangeRequest) (*pb.Empty, error) {

	data := cred.CredentialFull{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
		Password: in.Password,
	}

	err := serv.credApp.Change(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed change credential data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo),
			zap.String("pwd", in.Password))

		return &pb.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &pb.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *CredServiceRPC) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Empty, error) {

	data := cred.CredentialGet{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
	}

	err := serv.credApp.Delete(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete credential data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo))

		return &pb.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &pb.Empty{}, nil
}

// Get - Получение данных по email и метаданным.
func (serv *CredServiceRPC) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {

	inData := cred.CredentialGet{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
	}

	data, err := serv.credApp.Get(inData)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.GetResponse{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed get credential data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo))

		return &pb.GetResponse{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	out := &pb.GetResponse{
		Email:    data.Email,
		MetaInfo: data.MetaInfo,
		Password: data.Password,
	}

	return out, nil
}
