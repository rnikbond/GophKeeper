package rpc_services

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/binary"
	"GophKeeper/internal/server/app_services"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/data/binary"
)

type BinaryServiceRPC struct {
	pb.BinaryServiceServer

	credApp app_services.BinaryApp
	logger  *zap.Logger
}

// NewBinaryServiceRPC - Создание эклемпляра gRPC сервиса дял хранения бмнарных данных.
func NewBinaryServiceRPC(credApp app_services.BinaryApp) *BinaryServiceRPC {
	serv := &BinaryServiceRPC{
		credApp: credApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *BinaryServiceRPC) Create(ctx context.Context, in *pb.CreateRequest) (*pb.Empty, error) {

	data := binary.DataFull{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
		Bytes:    in.Data,
	}

	err := serv.credApp.Create(data)
	if err != nil {
		if err == errs.ErrAlreadyExist {
			return &pb.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())
		}

		serv.logger.Error("failed create binary data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo))

		return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	return &pb.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *BinaryServiceRPC) Change(ctx context.Context, in *pb.ChangeRequest) (*pb.Empty, error) {

	data := binary.DataFull{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
		Bytes:    in.Data,
	}

	err := serv.credApp.Change(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed change binary data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo))

		return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	return &pb.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *BinaryServiceRPC) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Empty, error) {

	data := binary.DataGet{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
	}

	err := serv.credApp.Delete(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete binary data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo))

		return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	return &pb.Empty{}, nil
}

// Get - Получение данных по email и метаданным.
func (serv *BinaryServiceRPC) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {

	inData := binary.DataGet{
		Email:    in.Email,
		MetaInfo: in.MetaInfo,
	}

	data, err := serv.credApp.Get(inData)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.GetResponse{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete get data",
			zap.Error(err),
			zap.String("email", in.Email),
			zap.String("meta", in.MetaInfo))

		return &pb.GetResponse{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	out := &pb.GetResponse{
		Email:    data.Email,
		MetaInfo: data.MetaInfo,
		Data:     data.Bytes,
	}

	return out, nil
}
