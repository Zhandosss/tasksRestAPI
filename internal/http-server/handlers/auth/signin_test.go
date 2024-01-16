package auth

import (
	"bytes"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"restAPI/internal/repositories"
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
)

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, login, password string)

	var tests = []struct {
		name                 string
		inputBody            string
		inputLogin           string
		inputPassword        string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "correct working",
			inputBody:     `{"login":"testLogin","password":"testPassword"}`,
			inputLogin:    "testLogin",
			inputPassword: "testPassword",
			mockBehavior: func(s *mock_service.MockAuthorization, login, password string) {
				s.EXPECT().GenerateToken(login, password).Return("testToken", nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"token":"testToken"}`,
		}, {
			name:                 "empty fields 1",
			inputBody:            `{"login":"testLogin"}`,
			inputLogin:           "testLogin",
			mockBehavior:         func(s *mock_service.MockAuthorization, login, password string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"wrong json fields"}`,
		}, {
			name:                 "empty fields 2",
			inputBody:            `{"papa":"23"}`,
			inputLogin:           "testLogin",
			mockBehavior:         func(s *mock_service.MockAuthorization, login, password string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"wrong json fields"}`,
		}, {
			name:                 "incorrect request",
			inputBody:            `{"login":"testLogin}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, login, password string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to decode request"}`,
		}, {
			name:          "error from GenerateToken: wrong login",
			inputBody:     `{"login":"testLogin","password":"testPassword"}`,
			inputLogin:    "testLogin",
			inputPassword: "testPassword",
			mockBehavior: func(s *mock_service.MockAuthorization, login, password string) {
				s.EXPECT().GenerateToken(login, password).Return("", repositories.ErrNoSuchUser)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"wrong login"}`,
		}, {
			name:          "error from GenerateToken: wrong password",
			inputBody:     `{"login":"testLogin","password":"testPassword"}`,
			inputLogin:    "testLogin",
			inputPassword: "testPassword",
			mockBehavior: func(s *mock_service.MockAuthorization, login, password string) {
				s.EXPECT().GenerateToken(login, password).Return("", repositories.ErrWrongPassword)
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"wrong password"}`,
		}, {
			name:          "error from GenerateToken: internal problem two same login in system",
			inputBody:     `{"login":"testLogin","password":"testPassword"}`,
			inputLogin:    "testLogin",
			inputPassword: "testPassword",
			mockBehavior: func(s *mock_service.MockAuthorization, login, password string) {
				s.EXPECT().GenerateToken(login, password).Return("", repositories.ErrTwoSameLoginInDb)
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to login"}`,
		}, {
			name:          "error from GenerateToken: internal problem",
			inputBody:     `{"login":"testLogin","password":"testPassword"}`,
			inputLogin:    "testLogin",
			inputPassword: "testPassword",
			mockBehavior: func(s *mock_service.MockAuthorization, login, password string) {
				s.EXPECT().GenerateToken(login, password).Return("", errors.New("test"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to login"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuthorization(ctrl)
			test.mockBehavior(auth, test.inputLogin, test.inputPassword)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			router := chi.NewRouter()
			router.Post("/auth/sign-in", SignIn(log, auth))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(test.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
