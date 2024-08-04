package usecase

import (
	"context"

	"github.com/bogatyr285/auth-go/internal/auth/entity"
	"github.com/bogatyr285/auth-go/internal/buildinfo"
	"github.com/bogatyr285/auth-go/internal/gateway/http/gen"
	"github.com/golang-jwt/jwt/v5"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, u entity.UserAccount) error
	FindUserByEmail(ctx context.Context, username string) (entity.UserAccount, error)
}

type CryptoPassword interface {
	HashPassword(password string) ([]byte, error)
	ComparePasswords(fromUser, fromDB string) bool
}

type JWTManager interface {
	IssueToken(userID string) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
}

type AuthUseCase struct {
	ur UserRepository
	cp CryptoPassword
	jm JWTManager
	bi buildinfo.BuildInfo
}

func NewUseCase(
	ur UserRepository,
	cp CryptoPassword,
	jm JWTManager,
	bi buildinfo.BuildInfo,
) AuthUseCase {
	return AuthUseCase{
		ur: ur,
		cp: cp,
		jm: jm,
		bi: bi,
	}
}

func (u AuthUseCase) PostLogin(ctx context.Context, request gen.PostLoginRequestObject) (gen.PostLoginResponseObject, error) {
	user, err := u.ur.FindUserByEmail(ctx, request.Body.Username)
	if err != nil {
		return gen.PostLogin500JSONResponse{}, nil
	}

	if !u.cp.ComparePasswords(user.Password, request.Body.Password) {
		return gen.PostLogin401JSONResponse{Error: "unauth"}, nil
	}

	token, err := u.jm.IssueToken(user.Username)
	if err != nil {
		return gen.PostLogin500JSONResponse{}, err
	}

	return gen.PostLogin200JSONResponse{
		AccessToken: token,
	}, nil
}

func (u AuthUseCase) PostRegister(ctx context.Context, request gen.PostRegisterRequestObject) (gen.PostRegisterResponseObject, error) {
	hashedPassword, err := u.cp.HashPassword(request.Body.Password)
	if err != nil {
		return gen.PostRegister500JSONResponse{}, nil
	}

	// TODO with New method
	user := entity.UserAccount{
		Username: request.Body.Username,
		Password: string(hashedPassword),
	}

	err = u.ur.RegisterUser(ctx, user)
	if err != nil {
		return gen.PostRegister500JSONResponse{}, nil
	}
	return gen.PostRegister201JSONResponse{
		Username: request.Body.Username,
	}, nil
}

func (u AuthUseCase) GetBuildinfo(ctx context.Context, request gen.GetBuildinfoRequestObject) (gen.GetBuildinfoResponseObject, error) {
	return gen.GetBuildinfo200JSONResponse{
		Arch:       u.bi.Arch,
		BuildDate:  u.bi.BuildDate,
		CommitHash: u.bi.CommitHash,
		Compiler:   u.bi.Compiler,
		GoVersion:  u.bi.GoVersion,
		Os:         u.bi.OS,
		Version:    u.bi.Version,
	}, nil
}
