package grpc_api

import (
	"context"
	"fmt"
	"strings"

	"github.com/techschool/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authriztionBearer   = "bearer"
)

func (server *Server) athorizeUser(ctx context.Context) (*token.PayLoad, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Missing metadata")
	}
	values := md.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("Missing authorization header")
	}

	authHeader := values[0]
	// Bearer asdasdasd
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}
	authType := strings.ToLower(fields[0])
	if authType != authriztionBearer {
		return nil, fmt.Errorf("Unsupported authorization type")
	}
	accessToken := fields[1]
	payLoad, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid access token: %s", err)
	}

	return payLoad, nil
}
