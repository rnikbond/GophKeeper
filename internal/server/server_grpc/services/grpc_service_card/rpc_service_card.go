package grpc_service_card

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/card"
	"GophKeeper/internal/server/app_services/app_service_card"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/proto/card"
)

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

	data := card.DataCard{
		MetaInfo: in.MetaInfo,
		Number:   in.Number,
		Period:   in.Period,
		CVV:      in.CVV,
		FullName: in.FullName,
	}

	// TODO :: errors.Is(...)
	err := serv.cardApp.Create(data)
	if err != nil {

		switch err {
		case errs.ErrAlreadyExist:
			return &card_store.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())

		case
			app_service_card.ErrInvalidPeriod,
			app_service_card.ErrInvalidNumber,
			app_service_card.ErrInvalidCVV,
			app_service_card.ErrInvalidFullName:
			return &card_store.Empty{}, status.Errorf(codes.InvalidArgument, err.Error())

		default:
			serv.logger.Error("failed create card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &card_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
		}
	}

	return &card_store.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *CardServiceRPC) Change(ctx context.Context, in *card_store.ChangeRequest) (*card_store.Empty, error) {

	data := card.DataCard{
		MetaInfo: in.MetaInfo,
		Number:   in.Number,
		Period:   in.Period,
		CVV:      in.CVV,
		FullName: in.FullName,
	}

	err := serv.cardApp.Change(data)
	if err != nil {

		switch err {
		case errs.ErrNotFound:
			return &card_store.Empty{}, status.Errorf(codes.NotFound, err.Error())

		case
			app_service_card.ErrInvalidPeriod,
			app_service_card.ErrInvalidNumber,
			app_service_card.ErrInvalidCVV,
			app_service_card.ErrInvalidFullName:
			return &card_store.Empty{}, status.Errorf(codes.InvalidArgument, err.Error())

		default:
			serv.logger.Error("failed change card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &card_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
		}
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

		switch err {
		case errs.ErrNotFound:
			return &card_store.Empty{}, status.Errorf(codes.NotFound, err.Error())

		default:
			serv.logger.Error("failed delete card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &card_store.Empty{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
		}
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

		switch err {
		case errs.ErrNotFound:
			return &card_store.GetResponse{}, status.Errorf(codes.NotFound, err.Error())

		default:
			serv.logger.Error("failed get card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &card_store.GetResponse{}, status.Errorf(codes.Internal, errs.ErrInternal.Error())
		}
	}

	return &card_store.GetResponse{
		Number:   get.Number,
		Period:   get.Period,
		CVV:      get.CVV,
		FullName: get.FullName,
	}, nil
}
