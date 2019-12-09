package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	BillingEventService api.BillingEventService
	CampaignService     api.CampaignService
}

func NewHandler(trackService api.TrackService, kpiService api.KPIService, billingEventService api.BillingEventService, campaignService api.CampaignService) Handler {
	return Handler{
		TrackService:        trackService,
		KPIService:          kpiService,
		BillingEventService: billingEventService,
		CampaignService:     campaignService,
	}
}

// Serve http
func (h *Handler) Serve(addr string) error {
	// Setup mux
	r := mux.NewRouter()
	r.HandleFunc("/v1/pixel/track", h.NewTrack).Methods("GET")
	r.HandleFunc("/v1/tracks/daily_visits", h.DailyVisits).Methods("GET")
	r.HandleFunc("/v1/tracks/top_pages", h.TopPages).Methods("GET")
	r.HandleFunc("/v1/tracks/most_active_campaigns", h.MostActiveCampaigns).Methods("GET")

	r.HandleFunc("/v1/kpis", h.NewKPI).Methods("POST")
	r.HandleFunc("/v1/kpis", h.KPIGetAll).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}", h.KPIGetOne).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}", h.KPIDelete).Methods("DELETE")
	r.HandleFunc("/v1/kpis/{kpi}/daily_conversion_count", h.KPIDailyConversionCount).Methods("GET")
	r.HandleFunc("/v1/kpis/{kpi}/first_touch", h.KPIFirstTouch).Methods("GET")

	r.HandleFunc("/v1/billing_events", h.NewBillingEvent).Methods("POST")
	// r.HandleFunc("/v1/kpis/by_user_id", h.KPIGetAll).Methods("GET")

	r.HandleFunc("/v1/campaigns", h.GetCampaigns).Methods("GET")

	return http.ListenAndServe(addr, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~
// =~ KPIs
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~

// NewKPI creates a new KPI
func (h *Handler) NewKPI(w http.ResponseWriter, r *http.Request) {
	// Get KPI data
	decoder := json.NewDecoder(r.Body)
	var kpi api.KPI
	err := decoder.Decode(&kpi)
	if err != nil {
		http.Error(w, "Invalid KPI object delivered. Expected {column, value}", 400)
		return
	}

	// TODO: validate the kpi data

	// Store
	id, err := h.KPIService.Store(kpi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	kpi.ID = id

	data, err := json.Marshal(kpi)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	// Write data back to client
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) KPIGetAll(w http.ResponseWriter, r *http.Request) {
	kpis, err := h.KPIService.Find()
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	// Marshall response
	jsonValue, err := json.Marshal(kpis)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonValue)
}

func (h *Handler) KPIGetOne(w http.ResponseWriter, r *http.Request) {
	// Get var from url
	vars := mux.Vars(r)
	kpiIDString := vars["kpi"]
	// Validate given kpi
	if len(kpiIDString) == 0 {
		http.Error(w, "Must specify kpi", 400)
		return
	}
	kpiIDInt, err := strconv.Atoi(kpiIDString)
	if err != nil {
		http.Error(w, "KPI ID must be a number", 400)
		return
	}

	kpi, err := h.KPIService.FindByID(kpiIDInt)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	// Marshall response
	jsonValue, err := json.Marshal(kpi)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write json back to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonValue)
}

func (h *Handler) KPIDelete(w http.ResponseWriter, r *http.Request) {
	// Get var from url
	vars := mux.Vars(r)
	kpiIDString := vars["kpi"]
	// Validate given kpi
	if len(kpiIDString) == 0 {
		http.Error(w, "Must specify kpi", 400)
		return
	}
	kpiIDInt, err := strconv.Atoi(kpiIDString)
	if err != nil {
		http.Error(w, "KPI ID must be a number", 400)
		return
	}

	count, err := h.KPIService.Delete(kpiIDInt)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	// Write json back to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Marshall response
	js, err := json.Marshal(count)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func (h *Handler) KPIDailyConversionCount(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	vars := mux.Vars(r)
	kpiIDString := vars["kpi"]

	// Validate given kpi
	if len(kpiIDString) == 0 {
		http.Error(w, "Must specify kpi", 400)
		return
	}
	kpiIDInt, err := strconv.Atoi(kpiIDString)
	if err != nil {
		http.Error(w, "KPI ID must be a number", 400)
		return
	}

	kpi, err := h.KPIService.FindByID(kpiIDInt)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}

	// Query
	dailyConversionCounts, err := h.TrackService.GetDailyConversionCountForKPI(kpi)
	if err != nil {
		log.Printf("ERORR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshall response
	js, err := json.Marshal(dailyConversionCounts)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write json back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *Handler) KPIFirstTouch(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	vars := mux.Vars(r)
	kpiIDString := vars["kpi"]

	// Validate given kpi
	if len(kpiIDString) == 0 {
		http.Error(w, "Must specify kpi", 400)
		return
	}
	kpiIDInt, err := strconv.Atoi(kpiIDString)
	if err != nil {
		http.Error(w, "KPI ID must be a number", 400)
		return
	}

	kpi, err := h.KPIService.FindByID(kpiIDInt)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}

	// Query
	firstTouches, err := h.TrackService.GetFirstTouchForKPI(kpi)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	// Marshall response
	js, err := json.Marshal(firstTouches)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	// Write json back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~
// =~ BillingEvents
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~
// NewBillingEvent creates a new billing event
func (h *Handler) NewBillingEvent(w http.ResponseWriter, r *http.Request) {
	// Get KPI data
	decoder := json.NewDecoder(r.Body)
	var billingEvent api.BillingEvent
	err := decoder.Decode(&billingEvent)
	if err != nil {
		http.Error(w, "Invalid KPI object delivered. Expected {column, value}", 400)
		panic(err)
	}

	// TODO: validate the kpi data

	// Store
	id, err := h.BillingEventService.Store(billingEvent)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	billingEvent.ID = id

	data, err := json.Marshal(billingEvent)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	// Write data back to client
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~
// =~ Campaigns
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~

func (h *Handler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.CampaignService.Find(mockOwnerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	data, err := json.Marshal(campaigns)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	// Write data back to client
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
