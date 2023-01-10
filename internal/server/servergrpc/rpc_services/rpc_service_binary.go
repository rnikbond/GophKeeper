package rpc_services

import (
	"GophKeeper/internal/server/app_services/app_service_binary"
	binary2 "GophKeeper/pkg/proto/binary"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/binary"
	"GophKeeper/pkg/errs"
)

type BinaryServiceRPC struct {
	binary2.BinaryServiceServer

	credApp app_service_binary.BinaryApp
	logger  *zap.Logger
}

// NewBinaryServiceRPC - Создание эклемпляра gRPC сервиса для хранения бинарных данных.
func NewBinaryServiceRPC(credApp app_service_binary.BinaryApp) *BinaryServiceRPC {
	serv := &BinaryServiceRPC{
		credApp: credApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *BinaryServiceRPC) Create(ctx context.Context, in *binary2.CreateRequest) (*binary2.Empty, error) {

	data := binary.DataFull{
		MetaInfo: in.MetaInfo,
		Bytes:    in.Data,
	}

	err := serv.credApp.Create(data)
	if err != nil {
		if err == errs.ErrAlreadyExist {
			return &binary2.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())
		}

		serv.logger.Error("failed create binary data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &binary2.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &binary2.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *BinaryServiceRPC) Change(ctx context.Context, in *binary2.ChangeRequest) (*binary2.Empty, error) {

	data := binary.DataFull{
		MetaInfo: in.MetaInfo,
		Bytes:    in.Data,
	}

	err := serv.credApp.Change(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &binary2.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed change binary data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &binary2.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &binary2.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *BinaryServiceRPC) Delete(ctx context.Context, in *binary2.DeleteRequest) (*binary2.Empty, error) {

	data := binary.DataGet{
		MetaInfo: in.MetaInfo,
	}

	err := serv.credApp.Delete(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &binary2.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete binary data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &binary2.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &binary2.Empty{}, nil
}

// Get - Получение данных по email и метаданным.
func (serv *BinaryServiceRPC) Get(ctx context.Context, in *binary2.GetRequest) (*binary2.GetResponse, error) {

	inData := binary.DataGet{
		MetaInfo: in.MetaInfo,
	}

	data, err := serv.credApp.Get(inData)
	if err != nil {
		if err == errs.ErrNotFound {
			return &binary2.GetResponse{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed get binary data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &binary2.GetResponse{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	out := &binary2.GetResponse{
		MetaInfo: data.MetaInfo,
		Data:     data.Bytes,
	}

	return out, nil
}
