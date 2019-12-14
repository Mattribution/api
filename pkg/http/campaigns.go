package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mattribution/api/pkg/api"
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

func (h *Handler) UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	// Get Campaign data
	decoder := json.NewDecoder(r.Body)
	var campaign api.Campaign
	err := decoder.Decode(&campaign)
	if err != nil {
		http.Error(w, "Invalid object delivered.", 400)
		panic(err)
	}

	// Attempt to update
	err = h.CampaignService.Update(campaign)
	if err != nil {
		http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		panic(err)
	}

	// Write data back to client
	w.WriteHeader(http.StatusOK)
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

	campaign, err := h.CampaignService.FindByID(campaignIDInt, mockOwnerID)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, "An error occurred", http.StatusInternalServerError)
		return
	}

	// Marshall response
	jsonValue, err := json.Marshal(campaign)
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

func (h *Handler) DailyConversionCountCampaign(w http.ResponseWriter, r *http.Request) {
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

	// Find by id
	campaign, err := h.CampaignService.FindByID(campaignIDInt, mockOwnerID)
	if err != nil {
		http.Error(w, "Not found", 404)
		return
	}

	// Query for daily conversions
	dailyConversionCounts, err := h.ConversionService.GetDailyByCampaign(campaign)
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
