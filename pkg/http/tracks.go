package http

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"gopkg.in/oleiade/reflections.v1"

	"github.com/mattribution/api/pkg/api"
	"github.com/mattribution/api/pkg/utils"
)

func (h *Handler) NewTrack(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	v := r.URL.Query()
	rawEvent := v.Get("data")
	data, err := base64.StdEncoding.DecodeString(rawEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	// Unmarshal
	track := api.Track{}
	if err := json.Unmarshal(data, &track); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}

	// TODO: real auth
	track.OwnerID = mockOwnerID

	// Grab IP
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	track.IP = ip

	// Store raw track
	// TODO: Make this all a transaction to fail if one can or can't
	newTrackID, err := h.TrackService.Store(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	track.ID = newTrackID

	// TODO: there must be a better way to do this...
	kpis, err := h.KPIService.Find(mockOwnerID)
	for _, kpi := range kpis {
		// Get the field name that matches the column patter from the kpi
		fieldName := utils.GetFieldName(kpi.Column, "db", track)
		if fieldName == "" { // field wasn't found
			http.Error(w, "Internal error", http.StatusInternalServerError)
			panic(err)
		}

		// Get the value for that field from the track
		value, err := reflections.GetField(track, fieldName)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			panic(err)
		}

		// If match
		if value == kpi.Value {
			// Create conversion
			conversion := api.Conversion{
				OwnerID: mockOwnerID,
				TrackID: track.ID,
				KPIID:   kpi.ID,
			}
			h.ConversionService.Store(conversion)

			// Get all tracks before to distribute weights
			tracks, err := h.TrackService.GetAllBySameUserBefore(track)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				panic(err)
			}

			if len(tracks) > 0 {
				// First Touch
				h.WeightService.Store(api.Weight{
					OwnerID:   track.OwnerID,
					KPIID:     kpi.ID,
					ModelName: "First Touch",
					Key:       tracks[0].CampaignName,
					Value:     1,
				})
			}

		}
	}

	// Write gif back to client
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/gif")
	w.Write(gif)
}

func (h *Handler) DailyVisits(w http.ResponseWriter, r *http.Request) {
	// TODO: Auth and get info on what data to look at

	// Query
	dailyVisits, err := h.TrackService.GetCountsFromColumn(30, `date_trunc('day', tracks.received_at)`, "tracks")
	if err != nil {
		panic(err)
	}

	// Marshall response
	js, err := json.Marshal(dailyVisits)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write json back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *Handler) TopPages(w http.ResponseWriter, r *http.Request) {
	// TODO: Auth and get info on what data to look at

	// Query
	topPages, err := h.TrackService.GetTopValuesFromColumn(30, "page_title", "tracks", "")
	if err != nil {
		log.Printf("ERORR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshall response
	js, err := json.Marshal(topPages)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write json back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *Handler) MostActiveCampaigns(w http.ResponseWriter, r *http.Request) {
	// TODO: Auth and get info on what data to look at

	// Query
	activeCampaigns, err := h.TrackService.GetTopValuesFromColumn(30, "campaign_name", "tracks", `WHERE campaign_name <> ''`)
	if err != nil {
		log.Printf("ERORR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshall response
	js, err := json.Marshal(activeCampaigns)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write json back to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
