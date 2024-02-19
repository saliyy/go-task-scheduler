package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-scheduler/internal/app/apiserver/handlers/task/save"
	"task-scheduler/internal/app/apiserver/handlers/task/save/mocks"
	dto "task-scheduler/internal/app/dto/task"
	"task-scheduler/internal/app/entities"
	slogdiscard "task-scheduler/internal/lib/logger/logdiscard"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveTaskHandler(t *testing.T) {
	t.Run("Success save", func(t *testing.T) {
		t.Parallel()

		taskSaverMock := mocks.NewTaskSaver(t)

		testingDto := &dto.CreateTaskDTO{
			Name:        "Test task",
			IsCompleted: false,
		}

		returningEntity := &entities.TaskEntity{
			Name:        testingDto.Name,
			IsCompleted: false,
			CreatedAt:   time.DateTime,
		}

		taskSaverMock.On("SaveTask", mock.Anything).Return(returningEntity).Once()

		handler := save.New(slogdiscard.NewDiscardLogger(), taskSaverMock)

		input := fmt.Sprintf(`{"name": "%s"}`, testingDto.Name)

		req, err := http.NewRequest(http.MethodPost, "/task", bytes.NewReader([]byte(input)))
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusOK)

		body := rr.Body.String()

		var resp save.Response

		require.NoError(t, json.Unmarshal([]byte(body), &resp))
	})
}
