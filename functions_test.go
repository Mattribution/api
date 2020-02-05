package functions

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	b64 "encoding/base64"

	"github.com/golang/mock/gomock"
	"github.com/mattribution/api/internal/app"
	mock_app "github.com/mattribution/api/internal/pkg/mock"
)

const (
	invalidJsonBase64Encoded = "ew=="
)

func TestNewTrack(t *testing.T) {

	t.Run("NewTrack responds with 400 BASE 64 ERROR if data is not base64 encoded", func(t *testing.T) {
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

		ctrl.Finish()
	})

	t.Run("NewTrack responds with 400 DATA FORMAT ERROR if base64 decoded data is invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		handler := Handler{
			Tracks: mock_app.NewMockTracks(ctrl),
		}

		data := invalidJsonBase64Encoded
		expectedStatus := http.StatusBadRequest
		expectedBody := invalidRequestError

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

		ctrl.Finish()
	})

	t.Run("NewTrack responds with 200 AND GIF if data is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		tracks := mock_app.NewMockTracks(ctrl)
		handler := Handler{
			Tracks: tracks,
		}

		track := app.Track{
			CampaignName: "My Campaign",
		}
		trackJSON, _ := json.Marshal(track)
		trackBase64 := b64.StdEncoding.EncodeToString(trackJSON)
		data := string(trackBase64)
		expectedStatus := http.StatusOK
		expectedBody := string(gif)

		tracks.EXPECT().Store(track).Times(1)

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

		ctrl.Finish()
	})

}
