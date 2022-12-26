package md_ctx

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

// ValueFromContext - Получение значения по ключу key из метаданных ctx.
func ValueFromContext(ctx context.Context, key string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		fmt.Printf("can not get metadata\n")
		return ``, false
	}

	values := md.Get(key)
	if len(values) != 1 {
		fmt.Printf("len values != 1 : %d for key %s\n", len(values), key)
		return ``, false
	}

	return values[0], true
}

// ValueToContext - Запись значения value по ключу key в метаданные ctx.
// При записи значений создается новый контекст. на базе ctx.
func ValueToContext(ctx context.Context, key, value string) context.Context {
	md := metadata.New(map[string]string{key: value})
	return metadata.NewOutgoingContext(ctx, md)
}
