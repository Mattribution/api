package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mattribution/api/internal/app"
)

const (
	invalidRequestError              = "The request you sent is invalid. Please reformat the request and try again."
	invalidBase64EncodingError       = "The data sent was not Base64 encoded. Please encode the data and try again."
	invalidJwtError                  = `{"error": "Invalid JWT"}`
	internalError                    = "We experienced an internal error. Please try again later."
	authClaimsDecodingError          = "Couldn't decode auth claims."
	mockOwnerID                int64 = 0
)

var (
	gif = []byte{
		71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
		255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
		1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
	}
	contextKeyClaims ContextKey = "claims"
)

type ContextKey string

type Handler struct {
	service     app.Service
	auth0Domain string
	auth0ApiID  string
}

func NewHandler(service app.Service, auth0Domain, auth0ApiID string) *Handler {
	return &Handler{
		service:     service,
		auth0Domain: auth0Domain,
		auth0ApiID:  auth0ApiID,
	}
}

// ServeHTTP sets up a router and serves http requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// // Setup auth0
	// jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
	// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	// 		// Verify 'aud' claim
	// 		aud := h.auth0ApiID
	// 		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	// 		if !checkAud {
	// 			return token, errors.New("Invalid audience")
	// 		}
	// 		// Verify 'iss' claim
	// 		iss := fmt.Sprintf("https://%s/", h.auth0Domain)
	// 		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	// 		if !checkIss {
	// 			return token, errors.New("Invalid issuer")
	// 		}

	// 		cert, err := getPemCert(token, h.auth0Domain)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}

	// 		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	// 		return result, nil
	// 	},
	// 	SigningMethod: jwt.SigningMethodRS256,
	// })

	router := mux.NewRouter()
	router.HandleFunc("/tracks/new", h.newTrack).Methods("GET")

	s := router.PathPrefix("/").Subrouter()
	s.HandleFunc("/kpis", h.newKpi).Methods("POST")
	s.HandleFunc("/kpis/{id:[0-9]+}", h.deleteKpi).Methods("DELETE")
	s.HandleFunc("/kpis", h.listKpis).Methods("GET")
	s.Use(h.newJwtMiddleware())
	s.Use(h.addJwtTokenClaimsInContextMiddleware)

	router.ServeHTTP(w, r)
}

// ~=~=~=~=~=~=~=~=
// Tracks
// ~=~=~=~=~=~=~=~=

func (h *Handler) newTrack(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	v := r.URL.Query()
	rawEvent := v.Get("data")
	data, err := base64.StdEncoding.DecodeString(rawEvent)
	if err != nil {
		http.Error(w, invalidBase64EncodingError, http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Unmarshal
	track := app.Track{}
	if err := json.Unmarshal(data, &track); err != nil {
		http.Error(w, invalidRequestError, http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Grab IP
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	track.IP = ip

	// Store raw track
	newTrackID, err := h.service.NewTrack(track)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error storing track: ", err)
		return
	}
	track.ID = newTrackID

	// Write gif back to client
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/gif")
	w.Write(gif)
}

// ~=~=~=~=~=~=~=~=
// Kpis
// ~=~=~=~=~=~=~=~=

func (h *Handler) newKpi(w http.ResponseWriter, r *http.Request) {
	var kpi app.Kpi
	claims := r.Context().Value(contextKeyClaims).(customClaims)

	// Parse body
	err := json.NewDecoder(r.Body).Decode(&kpi)
	if err != nil {
		http.Error(w, invalidRequestError, http.StatusBadRequest)
		log.Println(err)
		return
	}
	kpi.OwnerID = claims.UserID

	// Store KPI
	newKpiID, err := h.service.NewKpi(kpi)
	if err != nil {
		http.Error(w, internalError, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Response
	s := strconv.FormatInt(newKpiID, 10)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s)
}

func (h *Handler) deleteKpi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]
	claims := r.Context().Value(contextKeyClaims).(customClaims)

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "id error", http.StatusBadRequest)
		return
	}

	kpi := app.Kpi{
		ID:      id,
		OwnerID: claims.UserID,
	}

	// Delete KPI
	deleted, err := h.service.DeleteKpi(kpi)
	if err != nil {
		http.Error(w, internalError, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Response
	s := strconv.FormatInt(deleted, 10)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, s)
}

func (h *Handler) listKpis(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(contextKeyClaims).(customClaims)

	// Get Kpis
	kpis, err := h.service.GetKpisForUser(claims.UserID)
	if err != nil {
		http.Error(w, internalError, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Printf("Kpis: %+v\n", kpis)

	// Response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kpis)
}
