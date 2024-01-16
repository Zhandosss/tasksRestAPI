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
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
)

func TestHandler_UpdateTask(t *testing.T) {
	type MockBehavior func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string)

	var tests = []struct {
		name                 string
		inputBody            string
		stringTaskID         string
		taskID               int64
		userID               int64
		inputRequest         updateRequest
		mockBehavior         MockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "correct working",
			inputBody:    `{"text":"testText","tags":["testTag1", "testTag2"]}`,
			stringTaskID: "1",
			taskID:       1,
			userID:       1,
			inputRequest: updateRequest{
				Text: "testText",
				Tags: []string{"testTag1", "testTag2"},
			},
			mockBehavior: func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string) {
				s.EXPECT().UpdateTask(taskID, userID, text, tags).Return(nil)
			},
			expectedStatusCode: http.StatusNoContent,
		}, {
			name:                 "incorrect userID",
			userID:               -1,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"failed to get auth id"}`,
		}, {
			name:                 "bad request",
			inputBody:            `{"text":"testText","tags":["testTag1, "testTag2"]}`,
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to decode request"}`,
		}, {
			name:                 "empty text",
			userID:               1,
			inputBody:            `{"tags":["testTag1", "testTag2"]}`,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"there is no text in task"}`,
		}, {
			name:                 "no taskID",
			inputBody:            `{"text":"testText","tags":["testTag1", "testTag2"]}`,
			stringTaskID:         "",
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"failed to get task id from url"}`,
		}, {
			name:                 "incorrect taskID",
			inputBody:            `{"text":"testText","tags":["testTag1", "testTag2"]}`,
			stringTaskID:         "a1",
			userID:               1,
			mockBehavior:         func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"incorrect task id record"}`,
		}, {
			name:         "incorrect taskID return",
			inputBody:    `{"text":"testText","tags":["testTag1", "testTag2"]}`,
			stringTaskID: "1",
			taskID:       1,
			userID:       1,
			inputRequest: updateRequest{
				Text: "testText",
				Tags: []string{"testTag1", "testTag2"},
			},
			mockBehavior: func(s *mock_service.MockTask, taskID, userID int64, text string, tags []string) {
				s.EXPECT().UpdateTask(taskID, userID, text, tags).Return(errors.New("test"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"can't update task"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			task := mock_service.NewMockTask(ctrl)
			test.mockBehavior(task, test.taskID, test.userID, test.inputRequest.Text, test.inputRequest.Tags)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			router := chi.NewRouter()
			router.Put("/task/", Update(logger, task))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/task/", bytes.NewBufferString(test.inputBody))

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
