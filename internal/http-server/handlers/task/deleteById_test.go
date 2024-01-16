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
	"restAPI/internal/repositories"
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
)

func TestHandler_DeleteTask(t *testing.T) {
	type MockBehavior func(s *mock_service.MockTask, taskID, userID int64)

	var tests = []struct {
		name                 string
		stringTaskID         string
		taskID               int64
		userID               int64
		mockBehavior         MockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "correct working",
			stringTaskID: "1",
			taskID:       1,
			userID:       1,
			mockBehavior: func(s *mock_service.MockTask, taskID, userID int64) {
				s.EXPECT().DeleteTask(taskID, userID).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		}, {
			name:                 "incorrect userID",
			userID:               -1,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"message":"failed to get auth id"}`,
		}, {
			name:                 "empty taskID",
			stringTaskID:         "",
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to get task id from url"}`,
		}, {
			name:                 "incorrect taskID",
			stringTaskID:         "a1",
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"incorrect task id record"}`,
		}, {
			name:         "incorrect DeleteTask return: no content",
			stringTaskID: "1",
			taskID:       1,
			userID:       1,
			mockBehavior: func(s *mock_service.MockTask, taskID, userID int64) {
				s.EXPECT().DeleteTask(taskID, userID).Return(repositories.ErrNoTask)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"message":"there no task with this taskID"}`,
		}, {
			name:         "incorrect DeleteTask return: internal server error",
			stringTaskID: "1",
			taskID:       1,
			userID:       1,
			mockBehavior: func(s *mock_service.MockTask, taskID, userID int64) {
				s.EXPECT().DeleteTask(taskID, userID).Return(errors.New("test"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"can't found task"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := mock_service.NewMockTask(ctrl)
			test.mockBehavior(task, test.taskID, test.userID)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			router := chi.NewRouter()
			router.Delete("/task/", Delete(logger, task))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/task/", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("taskId", test.stringTaskID)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			r = r.WithContext(context.WithValue(r.Context(), "userID", test.userID))

			router.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
