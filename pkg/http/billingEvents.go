package http

import (
	"encoding/json"
	"net/http"

	"github.com/mattribution/api/pkg/api"
)

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
