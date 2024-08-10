package gapi

import (
	"context"
	"database/sql"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context,req *pb.LoginUserRequest) (*pb.LoginUserReponse, error) {
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to find user")
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Error(codes.NotFound, "incorrect password")
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		req.Username, server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create access token")
	}
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create refresh token")
	}
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		RefreshToken: refreshToken,
		ID:           refreshPayload.ID,
		Username:     user.Username,
		ExpiresAt:    refreshPayload.ExpiredAt,
		UserAgent:    "", //TODO: get user agent from request,
		ClientIp:     "",          //TODO: get IP address from request,
		IsBlocked:    false,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create session")
	}
	rsb := &pb.LoginUserReponse{
		User: convertUser(user),
		SessionId: session.ID.String(),
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
		RefreshToken: refreshToken,
		AccessToken: accessToken,

	}
	return rsb, nil
}