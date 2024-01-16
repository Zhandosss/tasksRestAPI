package task

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
)

func TestHandler_DeleteAll(t *testing.T) {
	type MockBehavior func(s *mock_service.MockTask, userID int64)

	var tests = []struct {
		name                 string
		userID               int64
		mockBehavior         MockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "correct working",
			userID: 1,

			mockBehavior: func(s *mock_service.MockTask, userID int64) {
				s.EXPECT().DeleteAllByUser(userID).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		}, {
			name:                 "incorrect userID",
			userID:               -1,
			mockBehavior:         func(s *mock_service.MockTask, userID int64) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"failed to get auth id"}`,
		}, {
			name:   "incorrect DeleteAllByUser return",
			userID: 1,

			mockBehavior: func(s *mock_service.MockTask, userID int64) {
				s.EXPECT().DeleteAllByUser(userID).Return(errors.New("test"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"couldn't delete all task"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := mock_service.NewMockTask(ctrl)
			test.mockBehavior(task, test.userID)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			router := chi.NewRouter()
			router.Delete("/tasks/", DeleteAll(logger, task))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/tasks/", nil)

			r = r.WithContext(context.WithValue(r.Context(), "userID", test.userID))

			router.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
