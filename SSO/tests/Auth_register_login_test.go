package tests

import (
	"authService/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	ssov1 "github.com/dima11223432/protos/gen/go/sso"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyApId      = 0
	defaultAppId   = 1
	appSecret      = "test-secret"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    defaultAppId,
	})
	loginTime := time.Now()

	require.NoError(t, err)
	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	assert.Equal(t, respReg.GetUserId(), int64(claims["id"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, defaultAppId, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

}

func TestRegisterLogin_DublicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")

}

func TestRegisterLogin_FailedCases(t *testing.T) {
	ctx, st := suite.New(t)

	testCases := []struct {
		name        string
		email       string
		password    string
		appId       int
		expectedErr string
	}{
		{
			name:        "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			appId:       defaultAppId,
			expectedErr: "password is empty",
		},
		{
			name:        "Register with empty email",
			email:       "",
			password:    randomFakePassword(),
			appId:       defaultAppId,
			expectedErr: "email is empty",
		},
		{
			name:        "Register with both empty fields",
			email:       "",
			password:    "",
			appId:       defaultAppId,
			expectedErr: "email is empty",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(
				ctx,
				&ssov1.RegisterRequest{
					Email:    tc.email,
					Password: tc.password,
				})
			require.Error(t, err)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
