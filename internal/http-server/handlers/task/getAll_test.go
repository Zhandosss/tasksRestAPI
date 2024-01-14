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
	mock_service "restAPI/internal/service/mocks"
	"strings"
	"testing"
	"time"
)

func TestHandler_GetAll(t *testing.T) {
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
				s.EXPECT().GetAllByUser(userID).Return([]model.Task{
					{
						ID:      1,
						Text:    "TestText",
						Tags:    []string{"testTag", "testTag2"},
						Date:    time.Date(1000, 10, 10, 10, 10, 10, 0, time.UTC),
						OwnerID: 1,
					}, {
						ID:      100,
						Text:    "TestText",
						Tags:    []string{"testTag"},
						Date:    time.Date(2000, 10, 10, 10, 10, 10, 0, time.UTC),
						OwnerID: 1,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"tasks":[{"text":"TestText","tags":["testTag","testTag2"],"date":"1000-10-10T10:10:10Z"},{"text":"TestText","tags":["testTag"],"date":"2000-10-10T10:10:10Z"}]}`,
		}, {
			name:                 "incorrect userID",
			userID:               -1,
			mockBehavior:         func(s *mock_service.MockTask, userID int64) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"failed to get auth id"}`,
		}, {
			name:   "incorrect GetAllByUser return",
			userID: 1,

			mockBehavior: func(s *mock_service.MockTask, userID int64) {
				s.EXPECT().GetAllByUser(userID).Return(nil, errors.New("test"))
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: `{"message":"couldn't get all task"}`,
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
			router.Get("/tasks/", GetAll(logger, task))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks/", nil)

			r = r.WithContext(context.WithValue(r.Context(), "userID", test.userID))

			router.ServeHTTP(w, r)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
