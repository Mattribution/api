package api

import "time"

// Track is event tracking data in our format
type Track struct {
	ID              int       `json:"id"`
	OwnerID         int64     `json:"ownerId"`
	UserID          string    `json:"userId"`
	FpHash          string    `json:"fpHash"`   // fingerprint hash
	PageURL         string    `json:"pageURL"`  // optional (website specific)
	PagePath        string    `json:"pagePath"` // optional ()
	PageTitle       string    `json:"pageTitle"`
	PageReferrer    string    `json:"pageReferrer"`
	Event           string    `json:"event"`
	CampaignSource  string    `json:"campaignSource"`
	CampaignMedium  string    `json:"campaignMedium"`
	CampaignName    string    `json:"campaignName"`
	CampaignContent string    `json:"campaignContent"`
	SentAt          time.Time `json:"sentAt"`
	IP              string
	Extra           string `json:"extra"` // (optional) extra json
}

type KPI struct {
	Column string `json:"column"`
	Value  string `json:"value"`
	Name   string `json:"name"`
}

func (kpi KPI) IsValid() bool {
	return len(kpi.Column) > 1 && len(kpi.Value) > 1 && len(kpi.Name) > 1
}

// ValueCount holds a count ascociated with a value
type ValueCount struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type TrackService interface {
	FindByID(id int) (Track, error)
	GetTopValuesFromColumn(days int, column, table string) ([]ValueCount, error)
	GetCountsFromColumn(days int, column, table string) ([]ValueCount, error)
	StoreTrack(t Track) (int, error)
	// DeleteTrack(id int) error
}

type KPIService interface {
	StoreKPI(kpi KPI) (int, error)
}
