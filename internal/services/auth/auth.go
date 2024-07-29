package auth

import (
	"auth/internal/domain/models"
	"auth/internal/lib/jwt"
	"auth/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var ErrorInvalidCredentials = errors.New("invalid credentials")
var ErrorInvalidApp = errors.New("invalid app id")
var ErrorUserExists = errors.New("user already exists")

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
	IsAdmin(ctx context.Context, userID int64) (bool, error)
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

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Info("failed to generate token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

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

		if errors.Is(err, storage.ErrorUserExists) {

			slog.Warn("user already exists", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			return 0, fmt.Errorf("%s, %t", op, ErrorUserExists)

		}

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
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorAppNotFound) {

			slog.Warn("app not found", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			return false, fmt.Errorf("%s, %t", op, ErrorInvalidApp)

		}
		return false, fmt.Errorf("%s, %t", op, err)
	}
	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
