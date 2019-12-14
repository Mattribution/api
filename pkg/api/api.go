package api

import "time"

// Track is event tracking data in our format
type Track struct {
	ID              int       `json:"id" db:"id"`
	OwnerID         int64     `json:"ownerId" db:"owner_id"`
	UserID          string    `json:"userId" db:"user_id"`
	AnonymousID     string    `json:"anonymousId" db:"anonymous_id"` // fingerprint hash
	PageURL         string    `json:"pageURL" db:"page_url"`         // optional (website specific)
	PagePath        string    `json:"pagePath" db:"page_path"`       // optional ()
	PageTitle       string    `json:"pageTitle" db:"page_title"`
	PageReferrer    string    `json:"pageReferrer" db:"page_referrer"`
	Event           string    `json:"event" db:"event"`
	IP              string    `json:"-" db:"sent_at"`
	CampaignSource  string    `json:"campaignSource" db:"campaign_source"`
	CampaignMedium  string    `json:"campaignMedium" db:"campaign_medium"`
	CampaignName    string    `json:"campaignName" db:"campaign_name"`
	CampaignContent string    `json:"campaignContent" db:"campaign_content"`
	ReceivedAt      time.Time `json:"-" db:"received_at"`
	SentAt          time.Time `json:"sentAt" db:"sent_at"`
	Extra           string    `json:"extra" db:"extra"` // (optional) extra json
}

type KPI struct {
	ID        int       `json:"id" db:"id"`
	OwnerID   int       `json:"-" db:"owner_id"`
	Column    string    `json:"column" db:"column_name"`
	Value     string    `json:"value" db:"value"`
	Name      string    `json:"name" db:"name"`
	Target    int       `json:"target" db:"target"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type Conversion struct {
	ID      int `json:"id" db:"id"`
	OwnerID int `json:"-" db:"owner_id"`
	TrackID int `json:"trackId" db:"track_id"`
	KPIID   int `json:"kpiId" db:"kpi_id"`
}

type Campaign struct {
	ID           int       `json:"id" db:"id"`
	OwnerID      int       `json:"ownerId" db:"owner_id"`
	Name         string    `json:"name" db:"name"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	CostPerMonth *float64  `json:"costPerMonth" db:"cost_per_month"`
	// Pattern to match
	ColumnName  string `json:"columnName" db:"column_name"`
	ColumnValue string `json:"columnValue" db:"column_value"`
}

type BillingEvent struct {
	ID        int       `json:"id" db:"id"`
	OwnerID   int       `json:"-" db:"owner_id"`
	UserID    int       `json:"userId" db:"user_id"`
	Amount    float32   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"-" db:"created_at"`
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
	Store(t Track) (int, error)
	FindByID(id int) (Track, error)
	GetTopValuesFromColumn(days int, column, table string, extraWheres string) ([]ValueCount, error)
	GetCountsFromColumn(days int, column, table string) ([]ValueCount, error)
	GetFirstTouchCount(kpi KPI) ([]ValueCount, error)
	// DeleteTrack(id int) error
}

type KPIService interface {
	Store(kpi KPI) (int, error)
	Find(ownerID int) ([]KPI, error)
	FindByID(id int) (KPI, error)
	Delete(int) (int64, error)
}

type BillingEventService interface {
	Store(billingEvent BillingEvent) (int, error)
	FindByUserID(id int) (BillingEvent, error)
}

type CampaignService interface {
	Store(campaign Campaign) (int, error)
	Update(campaifn Campaign) error
	Find(ownerID int) ([]Campaign, error)
	FindByID(id int, ownerID int) (Campaign, error)
	ScanForNewCampaigns(ownerID int) (int, error)
}

type ConversionService interface {
	Store(conversion Conversion) (int, error)
	Find(ownerID int) ([]Conversion, error)
	Delete(id int, ownerID int) (int64, error)
	GetDailyByCampaign(campaign Campaign) ([]ValueCount, error)
}
