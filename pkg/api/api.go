package api

import "time"

// Track is event tracking data in our format
type Track struct {
	ID              int       `json:"id" db:"id"`
	OwnerID         int64     `json:"ownerId" db:"owner_id"`
	UserID          string    `json:"userId" db:"user_id"`
	FpHash          string    `json:"fpHash" db:"fp_hash"`     // fingerprint hash
	PageURL         string    `json:"pageURL" db:"page_url"`   // optional (website specific)
	PagePath        string    `json:"pagePath" db:"page_path"` // optional ()
	PageTitle       string    `json:"pageTitle" db:"page_title"`
	PageReferrer    string    `json:"pageReferrer" db:"page_referrer"`
	Event           string    `json:"event" db:"event"`
	IP              string    `db:"sent_at"`
	CampaignSource  string    `json:"campaignSource" db:"campaign_source"`
	CampaignMedium  string    `json:"campaignMedium" db:"campaign_medium"`
	CampaignName    string    `json:"campaignName" db:"campaign_name"`
	CampaignContent string    `json:"campaignContent" db:"campaign_content"`
	ReceivedAt      time.Time `db:"received_at"`
	SentAt          time.Time `json:"sentAt" db:"sent_at"`
	Extra           string    `json:"extra" db:"extra"` // (optional) extra json
}

type KPI struct {
	ID     int    `json:"id" db:"id"`
	Column string `json:"column" db:"column_name"`
	Value  string `json:"value" db:"value"`
	Name   string `json:"name" db:"name"`
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
	StoreTrack(t Track) (int, error)
	FindByID(id int) (Track, error)
	GetTopValuesFromColumn(days int, column, table string) ([]ValueCount, error)
	GetCountsFromColumn(days int, column, table string) ([]ValueCount, error)
	GetDailyConversionCountForKPI(kpi KPI) ([]ValueCount, error)
	// DeleteTrack(id int) error
}

type KPIService interface {
	StoreKPI(kpi KPI) (int, error)
	FindByID(id int) (KPI, error)
}
