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
	"restAPI/internal/model"
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user model.User)

	var tests = []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Correct working",
			inputBody: `{"first_name":"firstnameTest","second_name":"secondnameTest","login":"testLogin","password":"passwordTest"}`,
			inputUser: model.User{
				FirstName:  "firstnameTest",
				SecondName: "secondnameTest",
				Login:      "testLogin",
				Password:   "passwordTest",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(int64(1), nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"user_id":1}`,
		}, {
			name:      "empty fields 1",
			inputBody: `{"first_name":"firstnameTest","login":"testLogin","password":"passwordTest"}`,
			inputUser: model.User{
				FirstName: "firstnameTest",
				Login:     "testLogin",
				Password:  "passwordTest",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user model.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"wrong json fields"}`,
		}, {
			name:                 "empty fields 2",
			inputBody:            `{"test":"test"}`,
			inputUser:            model.User{},
			mockBehavior:         func(s *mock_service.MockAuthorization, user model.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"wrong json fields"}`,
		}, {
			name:      "incorrect request 1",
			inputBody: `{"first_name:"firstnameTest","second_name":"secondnameTest","login":"testLogin","password":"passwordTest"}`,
			inputUser: model.User{
				FirstName:  "firstnameTest",
				SecondName: "secondnameTest",
				Login:      "testLogin",
				Password:   "passwordTest",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user model.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to decode request"}`,
		}, {
			name:      "incorrect request 2",
			inputBody: ``,
			inputUser: model.User{
				ID:         1,
				FirstName:  "firstnameTest",
				SecondName: "secondnameTest",
				Login:      "testLogin",
				Password:   "passwordTest",
			},
			mockBehavior:         func(s *mock_service.MockAuthorization, user model.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to decode request"}`,
		},
		{
			name:      "error from CreateUser",
			inputBody: `{"first_name":"firstnameTest","second_name":"secondnameTest","login":"testLogin","password":"passwordTest"}`,
			inputUser: model.User{
				FirstName:  "firstnameTest",
				SecondName: "secondnameTest",
				Login:      "testLogin",
				Password:   "passwordTest",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(int64(0), errors.New("test"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to register"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			auth := mock_service.NewMockAuthorization(ctrl)
			test.mockBehavior(auth, test.inputUser)

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			router := chi.NewRouter()
			router.Post("/auth/sign-up", SignUp(log, auth))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/sign-up", bytes.NewBufferString(test.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
