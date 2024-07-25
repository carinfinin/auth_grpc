package auth

import (
	"auth/internal/domain/models"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var ErrorInvalidCredentials = errors.New("invalid credentials")

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, email string) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, id int) (models.App, error)
}

// New returns instanceof the Auth
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login check if user with given credentials exist
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {

	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotFound) {
			a.log.Warn("user not found", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
		}

		a.log.Error("failed get user", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		return "", fmt.Errorf("%s: %w", op, err)

	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		return "", fmt.Errorf("%s: %w", op, ErrorInvalidCredentials)
	}

	return "", nil

}

// RegisterNewUser registers new user in the system and return ID
func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {

	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Error("failed to generate passwordHash", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)

	if err != nil {
		log.Error("failed to save user", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")

	return id, nil
}

// IsAdmin check if user is admin
func (a *Auth) IsAdmin(cnt context.Context, userID int64) (bool, error) {
	panic("not implemented")

}
