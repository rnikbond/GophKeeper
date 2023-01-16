//go:generate mockgen -source rpc_service_card.go -destination mocks/rpc_service_card_mock.go -package grpc_service_card
package grpc_service_card

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/client/client_grpc/services/card_service"
	"GophKeeper/internal/model/card"
	"GophKeeper/internal/server/app_services/app_service_card"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/proto/card"
)

type CardApp interface {
	Create(data card.DataCardFull) error
	Get(in card.DataCardGet) (card.DataCardFull, error)
	Delete(in card.DataCardGet) error
	Change(in card.DataCardFull) error
}

type CardServiceRPC struct {
	card_store.CardServiceServer

	cardApp app_service_card.CardApp
	logger  *zap.Logger
}

// NewCardServiceRPC - Создание эклемпляра gRPC сервиса для хранения данных банковских карт.
func NewCardServiceRPC(cardApp app_service_card.CardApp) *CardServiceRPC {
	serv := &CardServiceRPC{
		cardApp: cardApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *CardServiceRPC) Create(ctx context.Context, in *card_store.CreateRequest) (*card_store.Empty, error) {

	data := card.DataCardFull{
		MetaInfo: in.MetaInfo,
		Number:   string(in.Number),
		Period:   string(in.Period),
		CVV:      string(in.CVV),
		FullName: string(in.FullName),
	}

	err := serv.cardApp.Create(data)
	if err != nil {

		if errors.Is(err, errs.ErrAlreadyExist) {
			return &card_store.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())
		}

		if errors.Is(err, card_service.ErrInvalidPeriod) ||
			errors.Is(err, card_service.ErrInvalidNumber) ||
			errors.Is(err, card_service.ErrInvalidCVV) ||
			errors.Is(err, card_service.ErrInvalidFullName) {
			return &card_store.Empty{}, status.Errorf(codes.InvalidArgument, err.Error())
		}

		serv.logger.Error("failed create card data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &card_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &card_store.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *CardServiceRPC) Change(ctx context.Context, in *card_store.ChangeRequest) (*card_store.Empty, error) {

	data := card.DataCardFull{
		MetaInfo: in.MetaInfo,
		Number:   string(in.Number),
		Period:   string(in.Period),
		CVV:      string(in.CVV),
		FullName: string(in.FullName),
	}

	err := serv.cardApp.Change(data)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &card_store.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		if errors.Is(err, card_service.ErrInvalidPeriod) ||
			errors.Is(err, card_service.ErrInvalidNumber) ||
			errors.Is(err, card_service.ErrInvalidCVV) ||
			errors.Is(err, card_service.ErrInvalidFullName) {
			return &card_store.Empty{}, status.Errorf(codes.InvalidArgument, err.Error())
		}

		serv.logger.Error("failed change card data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &card_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &card_store.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *CardServiceRPC) Delete(ctx context.Context, in *card_store.DeleteRequest) (*card_store.Empty, error) {

	data := card.DataCardGet{
		MetaInfo: in.MetaInfo,
	}

	err := serv.cardApp.Delete(data)
	if err != nil {

		if errors.Is(err, errs.ErrNotFound) {
			return &card_store.Empty{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed delete card data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &card_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &card_store.Empty{}, nil
}

// Get - Получение существующих данных.
func (serv *CardServiceRPC) Get(ctx context.Context, in *card_store.GetRequest) (*card_store.GetResponse, error) {

	data := card.DataCardGet{
		MetaInfo: in.MetaInfo,
	}

	get, err := serv.cardApp.Get(data)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return &card_store.GetResponse{}, status.Errorf(codes.NotFound, err.Error())
		}

		serv.logger.Error("failed get card data",
			zap.Error(err),
			zap.String("meta", in.MetaInfo))

		return &card_store.GetResponse{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
	}

	return &card_store.GetResponse{
		Number:   []byte(get.Number),
		Period:   []byte(get.Period),
		CVV:      []byte(get.CVV),
		FullName: []byte(get.FullName),
	}, nil
}
