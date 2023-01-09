package rpc_services

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/card"
	"GophKeeper/internal/server/app_services"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/data/card"
)

type CardServiceRPC struct {
	pb.CardServiceServer

	cardApp app_services.CardApp
	logger  *zap.Logger
}

// NewCardServiceRPC - Создание эклемпляра gRPC сервиса для хранения данных банковских карт.
func NewCardServiceRPC(cardApp app_services.CardApp) *CardServiceRPC {
	serv := &CardServiceRPC{
		cardApp: cardApp,
		logger:  zap.L(),
	}

	return serv
}

// Create - Добавление новых данных.
func (serv *CardServiceRPC) Create(ctx context.Context, in *pb.CreateRequest) (*pb.Empty, error) {

	data := card.DataCard{
		MetaInfo: in.MetaInfo,
		Number:   in.Number,
		Period:   in.Period,
		CVV:      in.CVV,
		FullName: in.FullName,
	}

	err := serv.cardApp.Create(data)
	if err != nil {

		switch err {
		case errs.ErrAlreadyExist:
			return &pb.Empty{}, status.Errorf(codes.AlreadyExists, err.Error())

		case
			app_services.ErrInvalidPeriod,
			app_services.ErrInvalidNumber,
			app_services.ErrInvalidCVV,
			app_services.ErrInvalidFullName:
			return &pb.Empty{}, status.Errorf(codes.InvalidArgument, err.Error())

		default:
			serv.logger.Error("failed create card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
		}
	}

	return &pb.Empty{}, nil
}

// Change - Изменение существующих данных.
func (serv *CardServiceRPC) Change(ctx context.Context, in *pb.ChangeRequest) (*pb.Empty, error) {

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
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())

		case
			app_services.ErrInvalidPeriod,
			app_services.ErrInvalidNumber,
			app_services.ErrInvalidCVV,
			app_services.ErrInvalidFullName:
			return &pb.Empty{}, status.Errorf(codes.InvalidArgument, err.Error())

		default:
			serv.logger.Error("failed change card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
		}
	}

	return &pb.Empty{}, nil
}

// Delete - Удаление существующих данных.
func (serv *CardServiceRPC) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.Empty, error) {

	data := card.DataCardGet{
		MetaInfo: in.MetaInfo,
	}

	err := serv.cardApp.Delete(data)
	if err != nil {

		switch err {
		case errs.ErrNotFound:
			return &pb.Empty{}, status.Errorf(codes.NotFound, err.Error())

		default:
			serv.logger.Error("failed delete card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &pb.Empty{}, status.Errorf(codes.Internal, InternalErrorText)
		}
	}

	return &pb.Empty{}, nil
}

// Get - Получение существующих данных.
func (serv *CardServiceRPC) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {

	data := card.DataCardGet{
		MetaInfo: in.MetaInfo,
	}

	get, err := serv.cardApp.Get(data)
	if err != nil {

		switch err {
		case errs.ErrNotFound:
			return &pb.GetResponse{}, status.Errorf(codes.NotFound, err.Error())

		default:
			serv.logger.Error("failed get card data",
				zap.Error(err),
				zap.String("meta", in.MetaInfo))

			return &pb.GetResponse{}, status.Errorf(codes.Internal, InternalErrorText)
		}
	}

	return &pb.GetResponse{
		Number:   get.Number,
		Period:   get.Period,
		CVV:      get.CVV,
		FullName: get.FullName,
	}, nil
}
