package app

import "time"

// Track is event tracking data in our format
type Track struct {
	ID              int64     `json:"id" db:"id"`
	OwnerID         int64     `json:"ownerId" db:"owner_id"`
	UserID          string    `json:"userId" db:"user_id"`
	AnonymousID     string    `json:"anonymousId" db:"anonymous_id"` // fingerprint hash
	PageURL         string    `json:"pageURL" db:"page_url"`         // optional (website specific)
	PagePath        string    `json:"pagePath" db:"page_path"`       // optional ()
	PageTitle       string    `json:"pageTitle" db:"page_title"`
	PageReferrer    string    `json:"pageReferrer" db:"page_referrer"`
	Event           string    `json:"event" db:"event"`
	IP              string    `json:"ip" db:"ip"`
	CampaignSource  string    `json:"campaignSource" db:"campaign_source"`
	CampaignMedium  string    `json:"campaignMedium" db:"campaign_medium"`
	CampaignName    string    `json:"campaignName" db:"campaign_name"`
	CampaignContent string    `json:"campaignContent" db:"campaign_content"`
	SentAt          time.Time `json:"sentAt" db:"sent_at"`
	Extra           string    `json:"extra" db:"extra"` // (optional) extra json
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

// Kpi stores rules that can be matched on and recorded as conversions
type Kpi struct {
	ID                     int64     `json:"id" db:"id"`
	OwnerID                int64     `json:"-" db:"owner_id"`
	Name                   string    `json:"name" db:"name"`
	Target                 int64     `json:"target" db:"target"`
	DataWasChanged         bool      `json:"-" db:"-"`
	PatternMatchColumnName string    `json:"column" db:"pattern_match_column_name"`
	PatternMatchValue      string    `json:"value" db:"pattern_match_row_value"`
	CreatedAt              time.Time `json:"-" db:"created_at"`
}

type Tracks interface {
	Store(t Track) (int64, error)
}

type Kpis interface {
	Store(kpi Kpi) (int64, error)
	// Delete(int64) (int64, error)
}
