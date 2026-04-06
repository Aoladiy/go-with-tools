package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/config"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/database/queries"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/errs"
	"github.com/Aoladiy/go-with-tools-auth-microservice/internal/helpers"
	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	accessExp  = time.Minute * 15
	refreshExp = time.Hour * 24 * 7
	SignedOut  = "signed-out-token-"
)

type Microservice struct {
	gen.UnimplementedAuthMicroserviceServer
	q   *queries.Queries
	rdb *redis.Client
	p   *pgxpool.Pool
	c   config.Config
}

func New(q *queries.Queries, rdb *redis.Client, p *pgxpool.Pool, c config.Config) *Microservice {
	return &Microservice{q: q, rdb: rdb, p: p, c: c}
}

func (a *Microservice) SignUp(ctx context.Context, request *gen.SignUpRequest) (*gen.JWTResponse, error) {
	if len(request.Password) < 8 {
		return nil, errs.BadRequest(errors.New("password must be at least 8 characters"))
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		return nil, errs.Internal(err)
	}

	var jwtResponse gen.JWTResponse
	appErr := helpers.WithTx(ctx, a.p, a.q, func(timeout context.Context, q *queries.Queries) *errs.AppError {
		adminUser, err := q.CreateAdminUser(timeout, queries.CreateAdminUserParams{
			Email:        request.Email,
			PasswordHash: string(password),
		})
		if err != nil {
			return errs.FromPgErr(err)
		}
		var appErr *errs.AppError
		jwtResponse, appErr = generateJWTResponse(a.c.JwtSecret, strconv.FormatInt(adminUser.ID, 10))
		if appErr != nil {
			return appErr
		}
		return nil
	})
	if appErr != nil {
		return nil, appErr
	}
	return &jwtResponse, nil
}

func (a *Microservice) SignIn(ctx context.Context, request *gen.SignInRequest) (*gen.JWTResponse, error) {
	timeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	adminUser, err := a.q.GetAdminUser(timeout, request.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.Unauthorized(fmt.Errorf("no user with such email %w", err))
		}
		return nil, errs.Internal(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminUser.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, errs.Unauthorized(fmt.Errorf("wrong password %w", err))
	}
	jwtResponse, appErr := generateJWTResponse(a.c.JwtSecret, strconv.FormatInt(adminUser.ID, 10))
	if appErr != nil {
		return nil, appErr
	}
	return &jwtResponse, nil
}

func (a *Microservice) TokenRefresh(ctx context.Context, request *gen.TokenRefreshRequest) (*gen.JWTResponse, error) {
	withClaims, appErr := ParseToken(request.RefreshToken, a.c.JwtSecret)
	if appErr != nil {
		return nil, appErr
	}
	tokenSignedOutResponse, err := a.IsTokenSignedOut(ctx, &gen.IsTokenSignedOutRequest{Token: withClaims.Raw})
	if err != nil {
		return nil, errs.Internal(err)
	}
	if tokenSignedOutResponse.IsTokenSignedOut {
		return nil, errs.Unauthorized(fmt.Errorf("token is in signed out tokens cache"))
	}
	userID, err := withClaims.Claims.GetSubject()
	if err != nil {
		return nil, errs.Internal(err)
	}

	jwtResponse, appErr := generateJWTResponse(a.c.JwtSecret, userID)
	if appErr != nil {
		return nil, appErr
	}
	return &jwtResponse, nil
}

func (a *Microservice) SignOut(ctx context.Context, request *gen.SignOutRequest) (*emptypb.Empty, error) {
	accessToken, appErr := ParseToken(request.AccessToken, a.c.JwtSecret)
	if appErr != nil {
		return &emptypb.Empty{}, appErr
	}
	refreshToken, appErr := ParseToken(request.RefreshToken, a.c.JwtSecret)
	if appErr != nil {
		return &emptypb.Empty{}, appErr
	}
	accessTokenExp, err := accessToken.Claims.GetExpirationTime()
	if err != nil {
		return &emptypb.Empty{}, errs.BadRequest(err)
	}
	refreshTokenExp, err := refreshToken.Claims.GetExpirationTime()
	if err != nil {
		return &emptypb.Empty{}, errs.BadRequest(err)
	}
	set := a.rdb.Set(ctx, SignedOut+accessToken.Raw, true, accessTokenExp.Time.Sub(time.Now()))
	if set.Err() != nil {
		return &emptypb.Empty{}, errs.Internal(set.Err())
	}
	set = a.rdb.Set(ctx, SignedOut+refreshToken.Raw, true, refreshTokenExp.Time.Sub(time.Now()))
	if set.Err() != nil {
		return &emptypb.Empty{}, errs.Internal(set.Err())
	}
	return &emptypb.Empty{}, nil
}

func (a *Microservice) IsTokenSignedOut(ctx context.Context, request *gen.IsTokenSignedOutRequest) (*gen.IsTokenSignedOutResponse, error) {
	err := a.rdb.Get(ctx, SignedOut+request.Token).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return &gen.IsTokenSignedOutResponse{IsTokenSignedOut: false}, err
	} else if errors.Is(err, redis.Nil) {
		return &gen.IsTokenSignedOutResponse{IsTokenSignedOut: false}, nil
	}
	return &gen.IsTokenSignedOutResponse{IsTokenSignedOut: true}, nil
}
