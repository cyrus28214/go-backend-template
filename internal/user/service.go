package user

import (
	"backend/pkg/apperror"
	"backend/pkg/config"
	"backend/pkg/contextkey"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type IUserService interface {
	Login(ctx context.Context, req *LoginRequest) (string, *apperror.AppError)
}

type UserService struct {
	userRepo  *UserRepository
	jwtConfig *config.JwtConfig
}

func NewUserService(userRepo *UserRepository, jwtConfig *config.JwtConfig) IUserService {
	return &UserService{userRepo: userRepo, jwtConfig: jwtConfig}
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Login
// @error USER_NOT_FOUND
// @error INVALID_CREDENTIALS
// @error INTERNAL_SERVER_ERROR
func (s *UserService) Login(ctx context.Context, req *LoginRequest) (string, *apperror.AppError) {
	l := ctx.Value(contextkey.LoggerKey).(*slog.Logger)

	// 1. 查找用户
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.Debug("user not found", "error", err)
			return "", ErrUserNotFound.Wrap(err)
		}
		l.Error("database error", "error", err)
		return "", apperror.ErrInternal.Wrap(err)
	}

	// 2. 验证密码
	if req.Password != user.Password {
		l.Debug("invalid credentials")
		return "", ErrInvalidCredentials
	}

	// 3. 生成 JWT
	claims := jwtClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(s.jwtConfig.JwtSecretHex))
	if err != nil {
		l.Error("failed to sign JWT token", "error", err)
		return "", apperror.ErrInternal.Wrap(err)
	}

	return tokenStr, nil
}
