package tests

import (
	"auth/tests/suite"
	"github.com/brianvoe/gofakeit/v7"
	authV1 "github.com/carinfinin/auth_proto/gen/go/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	emptyAppId        = 0
	appID             = 1
	appSecret         = "test-secret"
	passDefaultLength = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {

	ctx, st := suite.New(t)

	email := gofakeit.Email()

	password := gofakeit.Password(true, true, true, true, true, passDefaultLength)

	responseRegister, err := st.AuthClient.Register(ctx, &authV1.RegisterRequest{
		Email:    email,
		Password: password,
	})

	require.NoError(t, err)

	assert.NotEmpty(t, responseRegister.GetId())

	responseLogin, err := st.AuthClient.Login(ctx, &authV1.LoginRequest{
		Email:    email,
		Password: password,
		App:      appID,
	})
	require.NoError(t, err)

	token := responseLogin.GetToken()
	assert.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)

	assert.True(t, ok)

	assert.Equal(t, responseRegister.GetId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
}
