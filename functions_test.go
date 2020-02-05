package functions

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mock_app "github.com/mattribution/api/internal/pkg/mock"
)

func TestNewTrack(t *testing.T) {

	t.Run("NewTrack responds with 500 BAD REQUEST if data is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		handler := Handler{
			Tracks: mock_app.NewMockTracks(ctrl),
		}

		data := "!"
		expectedStatus := http.StatusBadRequest
		expectedBody := invalidBase64EncodingError

		// Compose request
		req, err := http.NewRequest("GET", "/tracks/new", nil)
		if err != nil {
			t.Fatal(err)
		}
		q := req.URL.Query()
		q.Add("data", data)
		req.URL.RawQuery = q.Encode()

		// Create response recorder and http handler
		rr := httptest.NewRecorder()
		httphandler := http.HandlerFunc(handler.newTrack)

		// Execute request
		httphandler.ServeHTTP(rr, req)

		if status := rr.Code; status != expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, expectedStatus)
		}

		bodyStr := strings.TrimSpace(rr.Body.String())
		if bodyStr != expectedBody {
			t.Errorf("handler returned unexpected body: got %v want %v",
				bodyStr, expectedBody)
		}

	})

}
