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

const (
	internalServerErrorMsg = "Internal Error"
	mockOwnerID            = 1
)

type Handler struct {
	TrackService        api.TrackService
	KPIService          api.KPIService
	WeightService       api.WeightService
	ConversionService   api.ConversionService
	BillingEventService api.BillingEventService
	CampaignService     api.CampaignService
}

func NewHandler(trackService api.TrackService, kpiService api.KPIService, conversionService api.ConversionService, billingEventService api.BillingEventService, campaignService api.CampaignService) Handler {
	return Handler{
		TrackService:        trackService,
		KPIService:          kpiService,
		BillingEventService: billingEventService,
		CampaignService:     campaignService,
		ConversionService:   conversionService,
	}
}

// Serve sets up and serves http routes
func (h *Handler) Serve(addr string) error {
	// Setup mux
	r := mux.NewRouter()
	r.HandleFunc("/v1/pixel/track", h.NewTrack).Methods("GET")
	r.HandleFunc("/v1/tracks/daily_visits", h.DailyVisits).Methods("GET")
	r.HandleFunc("/v1/tracks/top_pages", h.TopPages).Methods("GET")
	r.HandleFunc("/v1/tracks/most_active_campaigns", h.MostActiveCampaigns).Methods("GET")

	r.HandleFunc("/v1/kpis", h.NewKPI).Methods("POST")
	r.HandleFunc("/v1/kpis", h.GetKPIs).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}", h.GetOneKPI).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}", h.DeleteKPI).Methods("DELETE")
	// r.HandleFunc("/v1/kpis/{kpi}/daily_conversion_count", h.DailyConversionCountKPI).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}/first_touch", h.FirstTouch).Methods("GET")

	r.HandleFunc("/v1/billing_events", h.NewBillingEvent).Methods("POST")

	r.HandleFunc("/v1/campaigns", h.GetCampaigns).Methods("GET")
	r.HandleFunc("/v1/campaigns/scan", h.ScanForNewCampaigns).Methods("GET")
	r.HandleFunc("/v1/campaigns/{campaign}", h.UpdateCampaign).Methods("PUT")
	r.HandleFunc("/v1/campaigns/{campaign}", h.GetOneCampaign).Methods("GET")
	r.HandleFunc("/v1/campaigns/{campaign}/daily_conversions", h.DailyConversionCountCampaign).Methods("GET")

	return http.ListenAndServe(addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
}
