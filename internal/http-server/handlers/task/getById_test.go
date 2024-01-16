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
	"restAPI/internal/model"
	"restAPI/internal/repositories"
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
	"time"
)

func TestHandler_GetTask(t *testing.T) {
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
				s.EXPECT().GetTask(taskID, userID).Return(model.Task{
					ID:      1,
					Text:    "TestText",
					Tags:    []string{"testTag1", "testTag2"},
					Date:    time.Date(1000, 10, 10, 10, 10, 10, 0, time.UTC),
					OwnerID: 1,
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"task":{"text":"TestText","tags":["testTag1","testTag2"],"date":"1000-10-10T10:10:10Z"}}`,
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
			name:         "incorrect GetTask return: no task",
			stringTaskID: "1",
			taskID:       1,
			userID:       1,
			mockBehavior: func(s *mock_service.MockTask, taskID, userID int64) {
				s.EXPECT().GetTask(taskID, userID).Return(model.Task{}, repositories.ErrNoTask)
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"message":"there no task with this taskID"}`,
		}, {
			name:         "incorrect GetTask return: internal server problem",
			stringTaskID: "1",
			taskID:       1,
			userID:       1,
			mockBehavior: func(s *mock_service.MockTask, taskID, userID int64) {
				s.EXPECT().GetTask(taskID, userID).Return(model.Task{}, errors.New("test"))
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
			router.Get("/task/", Get(logger, task))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/task/", nil)

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
