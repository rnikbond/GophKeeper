package md_ctx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// ValueFromContext - Получение значения по ключу key из метаданных ctx.
func ValueFromContext(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ``, false
	}

	values := md.Get(key)
	if len(values) == 0 {
		return ``, false
	}

	return values[0], true
}
