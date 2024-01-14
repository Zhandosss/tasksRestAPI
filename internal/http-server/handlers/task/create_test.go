package task

import (
	"bytes"
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
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
	"time"
)

func TestHandler_CreateTask(t *testing.T) {
	type MockBehavior func(s *mock_service.MockTask, task model.Task)

	var tests = []struct {
		name                 string
		inputBody            string
		inputTask            model.Task
		userID               int64
		mockBehavior         MockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "correct working with tags",
			inputBody: `{"text":"TestText","tags":["TestTag1","TestTag2"],"date":"1000-10-10T10:10:10.000Z"}`,
			inputTask: model.Task{
				Text:    "TestText",
				Tags:    []string{"TestTag1", "TestTag2"},
				Date:    time.Date(1000, 10, 10, 10, 10, 10, 0, time.UTC),
				OwnerID: 1,
			},
			userID: 1,
			mockBehavior: func(s *mock_service.MockTask, task model.Task) {
				s.EXPECT().CreateTask(task).Return(int64(1), nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"task_id":1}`,
		}, {
			name:      "correct working without tags",
			inputBody: `{"text":"TestText","date":"1000-10-10T10:10:10.000Z"}`,
			inputTask: model.Task{
				Text:    "TestText",
				Date:    time.Date(1000, 10, 10, 10, 10, 10, 0, time.UTC),
				OwnerID: 1,
			},
			userID: 1,
			mockBehavior: func(s *mock_service.MockTask, task model.Task) {
				s.EXPECT().CreateTask(task).Return(int64(1), nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"task_id":1}`,
		}, {
			name:                 "incorrect request",
			inputBody:            `{"text":"TestText,"date":"1000-10-10T10:10:10.000Z"}`,
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, task model.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to decode request"}`,
		}, {
			name:                 "incorrect user id",
			userID:               -1,
			mockBehavior:         func(s *mock_service.MockTask, task model.Task) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"incorrect userID"}`,
		}, {
			name:                 "empty text",
			inputBody:            `{"text":"","tags":["TestTag1","TestTag2"],"date":"1000-10-10T10:10:10.000Z"}`,
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, task model.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"incorrect task information"}`,
		}, {
			name:                 "too late",
			inputBody:            `{"text":"testText","tags":["TestTag1","TestTag2"],"date":"1000-10-10T10:10:09.000Z"}`,
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, task model.Task) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"incorrect task information"}`,
		}, {
			name:      "incorrect CreateTask work",
			inputBody: `{"text":"TestText","tags":["TestTag1","TestTag2"],"date":"1000-10-10T10:10:10.000Z"}`,
			inputTask: model.Task{
				Text:    "TestText",
				Tags:    []string{"TestTag1", "TestTag2"},
				Date:    time.Date(1000, 10, 10, 10, 10, 10, 0, time.UTC),
				OwnerID: 1,
			},
			userID: 1,
			mockBehavior: func(s *mock_service.MockTask, task model.Task) {
				s.EXPECT().CreateTask(task).Return(int64(0), errors.New("test"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"failed to create task"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := mock_service.NewMockTask(ctrl)
			test.mockBehavior(task, test.inputTask)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			router := chi.NewRouter()
			router.Post("/task/", Create(logger, task, time.Date(1000, 10, 10, 10, 10, 10, 0, time.UTC)))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/task/", bytes.NewBufferString(test.inputBody))
			ctx := r.Context()
			ctx = context.WithValue(ctx, "userID", test.userID)
			r = r.WithContext(ctx)

			router.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
