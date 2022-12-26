package interceptors

import (
	"GophKeeper/pkg/md_ctx"
	"context"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"GophKeeper/pkg/token"
)

// ValidateInterceptor - Перехватчик для gRPC, который отвечает за проверку подлинности JWT.
type ValidateInterceptor struct {
	// secretKey - Секретный ключ дял проверки подлинности JWT.
	secretKey string
}

// NewValidateInterceptor - Создание экземпляра перехватчика для валидации JWT.
func NewValidateInterceptor(key string) grpc.ServerOption {
	v := &ValidateInterceptor{secretKey: key}
	return grpc.UnaryInterceptor(middleware.ChainUnaryServer(v.ValidateTokenInterceptor))
}

// ValidateTokenInterceptor - Проверяет подлинность JWT.
// Токен должен быть в метаданных ctx по ключу "token".
// Если токен не найден или он не прошел проверку подлинности, то возвращается
// codes.PermissionDenied.
// Если токен валидный, то создается новый context на базе ctx, а в его метаданные
// записывается email пользователя (из токена) и новый context передается дальше в handler.
//
// При запросе Register или Login токен не проверяется.
func (inter ValidateInterceptor) ValidateTokenInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	// При регистрации и авторизации не проверяем токен
	if info.FullMethod == "/auth.AuthService/Register" || info.FullMethod == "/auth.AuthService/Login" {
		return handler(ctx, req)
	}

	tokenStr, ok := md_ctx.ValueFromContext(ctx, "token")
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "Failed read token")
	}

	jwtToken, err := token.VerifyJWT(tokenStr, inter.secretKey)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Invalid token")
	}

	email := jwtToken.Claims.(*token.Token).Email
	ctx = md_ctx.ValueToContext(ctx, "email", email)

	return handler(ctx, req)
}
