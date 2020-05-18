package token

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

//type grpcServer struct {
//	auth *Authentication
//}
//
//func (p *grpcServer) SomeMethod(
//	ctx context.Context, in *certificattedRpc.String,
//) (*certificattedRpc.String, error) {
//	if err := p.auth.Auth(ctx); err != nil {
//		return nil, err
//	}
//
//	return &certificattedRpc.String{Value: "Hello " + in.Value}, nil
//}

func (a *Authentication) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}

	var appid string
	var appkey string

	if val, ok := md["user"]; ok {
		appid = val[0]
	}
	if val, ok := md["password"]; ok {
		appkey = val[0]
	}

	if appid != a.GetAppKey() || appkey != a.GetAppSecret() {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return nil
}

func (a *Authentication) GetAppKey() string {
	return "wxning"
}

func (a *Authentication) GetAppSecret() string {
	return "gopher"
}
