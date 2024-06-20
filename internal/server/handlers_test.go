package server

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shorter/internal/models"
	"url-shorter/internal/server/services"
	"url-shorter/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateLink(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*services.LinkServiceMock)
		expectedStatus int
	}{
		{
			name:        "success",
			requestBody: `{"url":"https://example.com","short":"exmpl"}`,
			mockSetup: func(m *services.LinkServiceMock) {
				m.On("CreateLink", mock.AnythingOfType("models.Link")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "bad request on invalid json",
			requestBody:    `{"url":"https://example.com",}`,
			mockSetup:      func(m *services.LinkServiceMock) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "bad request on link already exists",
			requestBody: `{"url":"https://example.com","short":"exmpl"}`,
			mockSetup: func(m *services.LinkServiceMock) {
				m.On("CreateLink", mock.AnythingOfType("models.Link")).Return(models.ErrLinkAlreadyExists)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "internal server error on service failure",
			requestBody: `{"url":"https://example.com","short":"exmpl"}`,
			mockSetup: func(m *services.LinkServiceMock) {
				m.On("CreateLink", mock.AnythingOfType("models.Link")).Return(errors.New("unexpected error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLinkService := new(services.LinkServiceMock)
			linkServer := &LinkServer{
				logger:      &log.MockLogger{},
				linkService: mockLinkService,
			}

			tt.mockSetup(mockLinkService)

			body := bytes.NewBufferString(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/create", body)
			w := httptest.NewRecorder()

			linkServer.handleCreateLink(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockLinkService.AssertExpectations(t)
		})
	}
}
