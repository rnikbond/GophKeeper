package grpc_service_text

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/text"
	"GophKeeper/internal/server/app_services/app_service_text"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/proto/text"
)

type TextServiceRPC struct {
	text_store.TextServiceServer

	textApp app_service_text.TextApp
	logger  *zap.Logger
}

// NewTextServiceRPC - Создание эклемпляра gRPC сервиса для хранения текстовыъ данных.
func NewTextServiceRPC(textApp app_service_text.TextApp) *TextServiceRPC {
	serv := &TextServiceRPC{
		textApp: textApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *TextServiceRPC) Create(ctx context.Context, in *text_store.CreateRequest) (*text_store.Empty, error) {

	data := text.DataTextFull{
		MetaInfo: in.MetaInfo,
		Text:     string(in.Text),
	}

	err := serv.textApp.Create(data)
	if err != nil {
		if errors.Is(err, errs.ErrAlreadyExist) {
			return &text_store.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())
		}

		serv.logger.Error("failed create text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &text_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &text_store.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *TextServiceRPC) Change(ctx context.Context, in *text_store.ChangeRequest) (*text_store.Empty, error) {

	data := text.DataTextFull{
		MetaInfo: in.MetaInfo,
		Text:     string(in.Text),
	}

	err := serv.textApp.Change(data)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &text_store.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed change text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &text_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &text_store.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *TextServiceRPC) Delete(ctx context.Context, in *text_store.DeleteRequest) (*text_store.Empty, error) {

	data := text.DataTextGet{
		MetaInfo: in.MetaInfo,
	}

	err := serv.textApp.Delete(data)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &text_store.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &text_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &text_store.Empty{}, nil
}

// Get - Получение данных по email и метаданным.
func (serv *TextServiceRPC) Get(ctx context.Context, in *text_store.GetRequest) (*text_store.GetResponse, error) {

	inData := text.DataTextGet{
		MetaInfo: in.MetaInfo,
	}

	data, err := serv.textApp.Get(inData)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &text_store.GetResponse{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed get text data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &text_store.GetResponse{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	out := &text_store.GetResponse{
		MetaInfo: data.MetaInfo,
		Text:     []byte(data.Text),
	}

	return out, nil
}
