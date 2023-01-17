package grpc_service_card

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/client/model/card_model"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/card"
)

type CardService struct {
	rpc    pb.CardServiceClient
	logger *zap.Logger
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn) *CardService {
	return &CardService{
		rpc:    pb.NewCardServiceClient(conn),
		logger: zap.L(),
	}
}

func (serv CardService) Create(data card_model.Card, token string) error {
	dataReq := &pb.CreateRequest{
		MetaInfo: data.MetaInfo,
		Number:   data.Number,
		Period:   data.Period,
		CVV:      data.CVV,
		FullName: data.FullName,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Create(ctx, dataReq); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.AlreadyExists:
				return errs.ErrAlreadyExist

			default:
				serv.logger.Error("unknown gRPC error in card service Create()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}

func (serv CardService) Get(meta string, token string) (card_model.Card, error) {
	data := &pb.GetRequest{
		MetaInfo: meta,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := serv.rpc.Get(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return card_model.Card{}, errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in card service Get()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}

		}

		return card_model.Card{}, errs.ErrInternal
	}

	return card_model.Card{
		Number:   resp.Number,
		Period:   resp.Period,
		CVV:      resp.CVV,
		FullName: resp.FullName,
	}, nil
}

func (serv CardService) Delete(meta string, token string) error {
	data := &pb.DeleteRequest{
		MetaInfo: meta,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Delete(ctx, data); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in card service Delete()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}

func (serv CardService) Change(data card_model.Card, token string) error {
	dataReq := &pb.ChangeRequest{
		MetaInfo: data.MetaInfo,
		Number:   data.Number,
		Period:   data.Period,
		CVV:      data.CVV,
		FullName: data.FullName,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Change(ctx, dataReq)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in card service Change()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}
