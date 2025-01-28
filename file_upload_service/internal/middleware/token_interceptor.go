package middleware

import (
	"context"
	"strings"

	validation "github.com/1abobik1/Cloud-Storage/auth-service/pkg/auth/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TokenInterceptor создаёт интерцептор для проверки токенов
// Пример передачи данных на клиенте metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
func TokenInterceptor(publicKeyPath string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Извлекаем метаданные
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		metadata.Pairs()
		// Получаем токен
		authHeader := md["authorization"]
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid token format")
		}
		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

		// Валидируем токен
		claims, err := validation.ValidateToken(tokenString, publicKeyPath)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Добавляем claims в контекст
		newCtx := context.WithValue(ctx, "claims", claims)
		return handler(newCtx, req)
	}
}
