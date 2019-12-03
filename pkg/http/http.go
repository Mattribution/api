package http

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mattribution/api/pkg/api"
)

var gif = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

var (
	mockToken  = "mock-token"
	mockUserID = 1
)

type Handler struct {
	TrackService api.TrackService
	KPIService   api.KPIService
}

func NewHandler(trackService api.TrackService, kpiService api.KPIService) Handler {
	return Handler{
		TrackService: trackService,
		KPIService:   kpiService,
	}
}

// Serve http
func (h Handler) Serve(addr string) error {
	// Setup auth middleware
	amw := newAuthMiddleware(h)

	// Setup mux
	r := mux.NewRouter()
	r.Use(amw.Middleware)

	r.HandleFunc("/v1/pixel/track", NewTrack).Methods("GET")
	r.HandleFunc("/v1/tracks/daily_visits", DailyVisits).Methods("GET")
	r.HandleFunc("/v1/tracks/top_pages", TopPages).Methods("GET")
	r.HandleFunc("/v1/tracks/most_active_campaigns", MostActiveCampaigns).Methods("GET")

	r.HandleFunc("/v1/kpis", NewKPI).Methods("POST")
	r.HandleFunc("/v1/kpis", KPIGetAll).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}", KPIGetOne).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}", KPIDelete).Methods("DELETE")
	r.HandleFunc("/v1/kpis/{kpi}/daily_conversion_count", KPIDailyConversionCount).Methods("GET")

	// r.HandleFunc("/v1/billing_events", h.NewKPI).Methods("POST")
	// r.HandleFunc("/v1/kpis/by_user_id", h.KPIGetAll).Methods("GET")

	return http.ListenAndServe(addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
}
