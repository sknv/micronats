package handler

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sknv/micronats/app/auth/internal/models"
	proto "github.com/sknv/micronats/app/auth/rpc"
	xjwt "github.com/sknv/micronats/app/lib/x/jwt"
	xnats "github.com/sknv/micronats/app/lib/x/nats"
)

const (
	accessTokenExpiresIn  = time.Hour
	refreshTokenExpiresIn = time.Hour * 24 * 60 // 60 days

	claimSessionID = "sessionID"
	claimUserID    = "userID"
	claimExp       = "exp"
)

const (
	// TODO: translate
	errCredentialsRequired  = "login and password are required"
	errAuthorize            = "invalid login or password"
	errRefreshTokenRequired = "refreshToken is required"
	errToken                = "invalid token"
)

type Auth struct {
	JWT *xjwt.JWT
}

func (a *Auth) Authorize(ctx context.Context, in *proto.Credentials) (*proto.Tokens, error) {
	if in.Login == "" || in.Password == "" {
		return nil, xnats.StatusError(xnats.InvalidArgument, errCredentialsRequired)
	}

	log.Print("[INFO] login: ", in.Login)

	// TODO: fetch from db
	user := &models.User{
		ID:           "abc123",
		Login:        in.Login,
		PasswordHash: "$2a$10$M1lqPbhTbNzR.HyVvbElPOn0o.Rma.FJTnVzFuTwPSDiw8Q8Lac.G", // password
	}
	if err := user.Authenticate(in.Password); err != nil {
		return nil, xnats.StatusError(xnats.Unauthenticated, errAuthorize)
	}

	accessTkn, accessTknExpAt, refreshTkn := a.issueTokens()
	return &proto.Tokens{
		AccessToken:  accessTkn,
		ExpiresAt:    accessTknExpAt,
		RefreshToken: refreshTkn,
	}, nil
}

func (a *Auth) RefreshTokens(ctx context.Context, in *proto.RefreshToken) (*proto.Tokens, error) {
	if in.RefreshToken == "" {
		return nil, xnats.StatusError(xnats.InvalidArgument, errRefreshTokenRequired)
	}

	// fetch sessionID from claims
	token, decErr := a.JWT.Decode(in.RefreshToken)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, xnats.StatusError(xnats.InvalidArgument, errToken)
	}

	// TODO: fetch from the db
	sessionID := claims[claimSessionID].(string)
	session := &models.Session{
		ID:           sessionID,
		UserID:       "abc123",
		RefreshToken: in.RefreshToken,
	}

	log.Print("[INFO] sessionID: ", sessionID)

	// compare refresh tokens
	if !session.Verify(in.RefreshToken) {
		return nil, xnats.StatusError(xnats.InvalidArgument, errToken)
	}

	// check if provided token is valid
	if decErr != nil {
		// TODO: remove the session from the db
		return nil, xnats.StatusError(xnats.InvalidArgument, errToken)
	}

	// issue new tokens
	accessTkn, accessTknExpAt, refreshTkn := a.issueTokens()

	// TODO: save session to the db
	session.RefreshToken = refreshTkn

	return &proto.Tokens{
		AccessToken:  accessTkn,
		ExpiresAt:    accessTknExpAt,
		RefreshToken: refreshTkn,
	}, nil
}

func (a *Auth) RevokeToken(ctx context.Context, in *proto.RefreshToken) (*proto.Empty, error) {
	if in.RefreshToken == "" {
		return nil, xnats.StatusError(xnats.InvalidArgument, errRefreshTokenRequired)
	}

	// fetch sessionID from claims
	token, _ := a.JWT.Decode(in.RefreshToken)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, xnats.StatusError(xnats.InvalidArgument, errToken)
	}

	// TODO: fetch from the db
	sessionID := claims[claimSessionID].(string)
	session := &models.Session{
		ID:           sessionID,
		UserID:       "abc123",
		RefreshToken: in.RefreshToken,
	}

	log.Print("[INFO] sessionID: ", sessionID)

	// compare refresh tokens
	if !session.Verify(in.RefreshToken) {
		return nil, xnats.StatusError(xnats.InvalidArgument, errToken)
	}

	// TODO: remove the session from the db
	return &proto.Empty{}, nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (a *Auth) issueTokens() (accessToken string, accessTokenExpiresAt int64, refreshToken string) {
	accessTokenExpiresAt = time.Now().Add(accessTokenExpiresIn).Unix()
	accessToken, _ = a.JWT.Encode(jwt.MapClaims{ // issue an access token
		claimUserID: "abc123",
		claimExp:    accessTokenExpiresAt,
	})
	refreshToken, _ = a.JWT.Encode(jwt.MapClaims{ // issue a refresh token
		claimSessionID: "def456",
		claimExp:       time.Now().Add(refreshTokenExpiresIn).Unix(),
	})
	return
}
