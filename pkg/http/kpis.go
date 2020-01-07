package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mattribution/api/pkg/api"
)

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

	// TODO: real auth
	kpi.OwnerID = mockOwnerID

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

func (h *Handler) GetKPIs(w http.ResponseWriter, r *http.Request) {
	kpis, err := h.KPIService.Find(mockOwnerID)
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

func (h *Handler) GetOneKPI(w http.ResponseWriter, r *http.Request) {
	// Get var from url
	vars := mux.Vars(r)
	kpiIDString := vars["kpi"]
	// Validate given kpi
	if len(kpiIDString) == 0 {
		http.Error(w, "Must specify kpi", 400)
		return
	}
	kpiIDInt, err := strconv.ParseInt(kpiIDString, 10, 64)
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

func (h *Handler) DeleteKPI(w http.ResponseWriter, r *http.Request) {
	// Get var from url
	vars := mux.Vars(r)
	kpiIDString := vars["kpi"]
	// Validate given kpi
	if len(kpiIDString) == 0 {
		http.Error(w, "Must specify kpi", 400)
		return
	}
	kpiIDInt, err := strconv.ParseInt(kpiIDString, 10, 64)
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

func (h *Handler) FirstTouch(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	vars := mux.Vars(r)
	kpiIDString := vars["kpi"]

	// Validate given kpi
	if len(kpiIDString) == 0 {
		http.Error(w, "Must specify kpi", 400)
		return
	}
	kpiIDInt, err := strconv.ParseInt(kpiIDString, 10, 64)
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
	firstTouches, err := h.TrackService.GetFirstTouchCount(kpi)
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
