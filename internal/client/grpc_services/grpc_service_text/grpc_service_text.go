package grpc_service_text

import (
	"GophKeeper/internal/client/model/text_model"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/text"
)

type TextService struct {
	rpc    pb.TextServiceClient
	logger *zap.Logger
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn) *TextService {
	return &TextService{
		rpc:    pb.NewTextServiceClient(conn),
		logger: zap.L(),
	}
}

func (serv TextService) Create(data text_model.Text, token string) error {
	dataReq := &pb.CreateRequest{
		MetaInfo: data.MetaInfo,
		Text:     data.Data,
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Create(ctx, dataReq)
	if err == nil {
		return nil
	}

	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.AlreadyExists:
			return errs.ErrAlreadyExist
		default:
			serv.logger.Error("unknown gRPC error in text service Create()",
				zap.Uint32("gRPC code", uint32(e.Code())),
				zap.String("gRPC text", e.String()))
		}
	}

	return errs.ErrInternal
}

func (serv TextService) Get(meta, token string) (text_model.Text, error) {
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
				return text_model.Text{}, errs.ErrNotFound

			default:
				serv.logger.Error("unknown gRPC error in text service Get()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return text_model.Text{}, errs.ErrInternal
	}

	return text_model.Text{
		MetaInfo: meta,
		Data:     resp.Text,
	}, nil
}

func (serv TextService) Delete(meta, token string) error {
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
				serv.logger.Error("unknown gRPC error in text service Delete()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}

func (serv TextService) Change(data text_model.Text, token string) error {
	dataReq := &pb.ChangeRequest{
		MetaInfo: data.MetaInfo,
		Text:     data.Data,
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
				serv.logger.Error("unknown gRPC error in text service Change()",
					zap.Uint32("gRPC code", uint32(e.Code())),
					zap.String("gRPC text", e.String()))
			}
		}
		return errs.ErrInternal
	}

	return nil
}
