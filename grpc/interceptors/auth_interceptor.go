package interceptors

import (
	"context"
	"log"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var errMetadataIsNotProvided = status.Errorf(codes.Unauthenticated, "metadata is not provided")
var errTokenIsNotProvided = status.Errorf(codes.Unauthenticated, "authorization token is not provided")

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMetadataIsNotProvided
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, errTokenIsNotProvided
	}

	tokenString := values[0]

	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(strings.ReplaceAll(os.Getenv("PUBLIC_KEY"), "\\n", "\n")))

	header, _ := stripBearerPrefixFromTokenString(strings.TrimSpace(tokenString))

	if header == "" {
		return nil, errTokenIsNotProvided
	}

	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		log.Println(err)
		return nil, errTokenIsNotProvided
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx = context.WithValue(ctx, "User", claims["sub"].(string))

		return handler(ctx, req)
	}

	return handler(ctx, req)
}

func stripBearerPrefixFromTokenString(tok string) (string, error) {
	if len(tok) > 6 && strings.ToUpper(tok[0:7]) == "BEARER " {
		return tok[7:], nil
	}
	return tok, nil
}
