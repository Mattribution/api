package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) GetCampaigns(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.CampaignService.Find(mockOwnerID)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
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

func (h *Handler) GetOneCampaign(w http.ResponseWriter, r *http.Request) {
	// Get var from url
	vars := mux.Vars(r)
	campaignIDString := vars["campaign"]
	// Validate given id
	if len(campaignIDString) == 0 {
		http.Error(w, "Must specify id", 400)
		return
	}
	campaignIDInt, err := strconv.Atoi(campaignIDString)
	if err != nil {
		http.Error(w, "ID must be a number", 400)
		return
	}

	kpi, err := h.CampaignService.FindByID(campaignIDInt, mockOwnerID)
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

func (h *Handler) ScanForNewCampaigns(w http.ResponseWriter, r *http.Request) {
	createCount, err := h.CampaignService.ScanForNewCampaigns(mockOwnerID)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	data, err := json.Marshal(createCount)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	// Write data back to client
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
