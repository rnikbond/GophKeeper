package rpc_services

import (
	"GophKeeper/internal/model/text"
	"GophKeeper/internal/server/app_services"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/data/text"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TextServiceRPC struct {
	pb.TextServiceServer

	textApp app_services.TextApp
	logger  *zap.Logger
}

// NewTextServiceRPC - Создание эклемпляра gRPC сервиса для хранения текстовыъ данных.
func NewTextServiceRPC(textApp app_services.TextApp) *TextServiceRPC {
	serv := &TextServiceRPC{
		textApp: textApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *TextServiceRPC) Create(ctx context.Context, in *pb.CreateRequest) (*pb.Empty, error) {

	data := text.DataTextFull{
		MetaInfo: in.MetaInfo,
		Text:     in.Text,
	}

	err := serv.textApp.Create(data)
	if err != nil {
		if err == errs.ErrAlreadyExist {
			return &pb.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())
		}

		serv.logger.Error("failed create text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	return &pb.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *TextServiceRPC) Change(ctx context.Context, in *pb.ChangeRequest) (*pb.Empty, error) {

	data := text.DataTextFull{
		MetaInfo: in.MetaInfo,
		Text:     in.Text,
	}

	err := serv.textApp.Change(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed change text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	return &pb.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *TextServiceRPC) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Empty, error) {

	data := text.DataTextGet{
		MetaInfo: in.MetaInfo,
	}

	err := serv.textApp.Delete(data)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	return &pb.Empty{}, nil
}

// Get - Получение данных по email и метаданным.
func (serv *TextServiceRPC) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {

	inData := text.DataTextGet{
		MetaInfo: in.MetaInfo,
	}

	data, err := serv.textApp.Get(inData)
	if err != nil {
		if err == errs.ErrNotFound {
			return &pb.GetResponse{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed get text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &pb.GetResponse{}, status.Errorf(codes.Internal, InternalErrorText)
	}

	out := &pb.GetResponse{
		MetaInfo: data.MetaInfo,
		Text:     data.Text,
	}

	return out, nil
}
