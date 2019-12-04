package http

import (
	"log"
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
	mockToken      = "mock-token"
	mockUserID     = 1
	mocktokenUsers = map[string]int{
		"000000": 1,
	}
)

type DataHandler struct {
	UserService  api.UserService
	TrackService api.TrackService
	KPIService   api.KPIService
}

func DataHandler(userService api.UserService, trackService api.TrackService, kpiService api.KPIService) Handler {
	return DataHandler{
		UserService:  userService,
		TrackService: trackService,
		KPIService:   kpiService,
	}
}

// Middleware function, which will be called for each request
func (dh DataHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")

		if userID, found := mocktokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)

			user, err := dh.UserService.FindByID(userID)
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			// If the user has an active custom presto data store,
			if user.PrestoDataStoreID != nil {

			}

			// TODO:
			// Find the user's data sources
			// if contains presto replacement
			//   replace presto service
			//
			// ENDTODO:

			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
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
