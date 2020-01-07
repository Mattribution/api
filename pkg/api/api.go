package api

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jmoiron/sqlx/types"
)

// Track is event tracking data in our format
type Track struct {
	ID              int64      `json:"id" db:"id"`
	OwnerID         int64      `json:"ownerId" db:"owner_id"`
	UserID          *string    `json:"userId" db:"user_id"`
	AnonymousID     *string    `json:"anonymousId" db:"anonymous_id"` // fingerprint hash
	PageURL         *string    `json:"pageURL" db:"page_url"`         // optional (website specific)
	PagePath        *string    `json:"pagePath" db:"page_path"`       // optional ()
	PageTitle       *string    `json:"pageTitle" db:"page_title"`
	PageReferrer    *string    `json:"pageReferrer" db:"page_referrer"`
	Event           *string    `json:"event" db:"event"`
	IP              *string    `json:"ip" db:"ip"`
	CampaignSource  *string    `json:"campaignSource" db:"campaign_source"`
	CampaignMedium  *string    `json:"campaignMedium" db:"campaign_medium"`
	CampaignName    *string    `json:"campaignName" db:"campaign_name"`
	CampaignContent *string    `json:"campaignContent" db:"campaign_content"`
	ReceivedAt      *time.Time `json:"receivedAt" db:"received_at"`
	SentAt          *time.Time `json:"sentAt" db:"sent_at"`
	Extra           *string    `json:"extra" db:"extra"` // (optional) extra json
}

// Campaign holds data usually connected to a campaign found in tracks
type Campaign struct {
	ID           int64     `json:"id" db:"id"`
	OwnerID      int64     `json:"ownerId" db:"owner_id"`
	Name         string    `json:"name" db:"name"`
	CreatedAt    time.Time `json:"-" db:"created_at"`
	CostPerMonth *float64  `json:"costPerMonth" db:"cost_per_month"`
	// Pattern to match
	ColumnName  string `json:"columnName" db:"column_name"`
	ColumnValue string `json:"columnValue" db:"column_value"`
}

// KPI stores rules that can be matched on and recorded as conversions
type KPI struct {
	ID             int64          `json:"id" db:"id"`
	OwnerID        int64          `json:"-" db:"owner_id"`
	Column         string         `json:"column" db:"column_name"`
	Value          string         `json:"value" db:"value"`
	Name           string         `json:"name" db:"name"`
	Data           types.JSONText `json:"data" db:"data"`
	Target         int64          `json:"target" db:"target"`
	CreatedAt      time.Time      `json:"-" db:"created_at"`
	DataWasChanged bool           `json:"-" db:"-"`
}

// AdjustWeight will adjust the value for a weight for a given model and key.
//  Creates the weight if it doesn't already exist.
func (kpi *KPI) AdjustWeight(modelName, attribute, key string, delta float32) (bool, error) {
	// Parse data json
	// modelName->attribute->weights
	data := make(map[string]map[string]map[string]float32)
	if kpi.Data != nil {
		if err := json.Unmarshal([]byte(kpi.Data), &data); err != nil {
			log.Println("ERROR: could not unmarshal kpi weights, resetting entire object")
		}
	}

	// (if entry for model doesn't exist create it)
	modelObj, modelObjExists := data[modelName]
	if !modelObjExists {
		modelObj = make(map[string]map[string]float32)
		data[modelName] = modelObj
		kpi.DataWasChanged = true
	}

	// (if entry for attribute doesn't exist create it)
	attributeObj, attributeObjExists := modelObj[attribute]
	if !attributeObjExists {
		attributeObj = make(map[string]float32)
		modelObj[attribute] = attributeObj
		kpi.DataWasChanged = true
	}

	// Update weight with delta
	attributeObj[key] += delta

	// Remarshal and update string
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	kpi.Data = jsonBytes
	kpi.DataWasChanged = true

	return true, nil
}

// ClearWeightsForModel removes all weights for a given model
func (kpi *KPI) ClearWeightsForModel(modelName string) (bool, error) {
	// Parse data json
	data := make(map[string]ModelData)
	if kpi.Data != nil {
		if err := json.Unmarshal([]byte(kpi.Data), &data); err != nil {
			log.Println("ERROR: could not unmarshal kpi weights, resetting entire object")
		}
	}

	// If the model entry doesn't exist, do nothing
	_, ok := data[modelName]
	if !ok {
		return false, nil
	}

	// If the model entry does exist, remove it
	delete(data, modelName)

	// Remarshal and update string
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return false, err
	}
	kpi.Data = jsonBytes
	kpi.DataWasChanged = true

	return true, nil
}

// IsValid checks if a KPI is valid and ok to be created
func (kpi KPI) IsValid() bool {
	return len(kpi.Column) > 1 && len(kpi.Value) > 1 && len(kpi.Name) > 1
}

// Conversion is a specific match for a KPI
type Conversion struct {
	ID      int64 `json:"id" db:"id"`
	OwnerID int64 `json:"-" db:"owner_id"`
	TrackID int64 `json:"trackId" db:"track_id"`
	KPIID   int64 `json:"kpiId" db:"kpi_id"`
}

// BillingEvent is a structure that holds billing data in order to calculate ROI
type BillingEvent struct {
	ID        int64     `json:"id" db:"id"`
	OwnerID   int64     `json:"-" db:"owner_id"`
	UserID    int64     `json:"userId" db:"user_id"`
	Amount    float32   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

// Weight holds a weight for a model for a given KPI
type Weight struct {
	ID        int64     `json:"id" db:"id"`
	OwnerID   int64     `json:"-" db:"owner_id"`
	KPIID     int64     `json:"kpiId" db:"kpi_id"`
	ModelName string    `json:"modelName" db:"model_name"`
	Key       string    `json:"key" db:"key"`
	Value     float32   `json:"value" db:"value"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

// ValueCount holds a count ascociated with a value
type ValueCount struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type TrackService interface {
	Store(t Track) (int64, error)
	FindByID(id int64) (Track, error)
	GetTopValuesFromColumn(days int, column, table string, extraWheres string) ([]ValueCount, error)
	GetCountsFromColumn(days int, column, table string) ([]ValueCount, error)
	GetFirstTouchCount(kpi KPI) ([]ValueCount, error)
	GetAllBySameUserBefore(Track) ([]Track, error)
	// DeleteTrack(id int) error
}

type KPIService interface {
	Store(kpi KPI) (int64, error)
	Find(ownerID int64) ([]KPI, error)
	FindByID(id int64) (KPI, error)
	UpdateData(KPI) error
	Delete(int64) (int64, error)
}

type WeightService interface {
	Store(weight Weight) (int64, error)
	FindByKpiAndModelName(ownerID int64, kpiID int64, modelName string) ([]Weight, error)
}

type BillingEventService interface {
	Store(billingEvent BillingEvent) (int64, error)
	FindByUserID(id int64) (BillingEvent, error)
}

type CampaignService interface {
	Store(campaign Campaign) (int64, error)
	ScanForNewCampaigns(ownerID int64) (int64, error)
	Update(campaifn Campaign) error
	Find(ownerID int64) ([]Campaign, error)
	FindByID(id int64, ownerID int64) (Campaign, error)
}

type ConversionService interface {
	Store(conversion Conversion) (int64, error)
	Find(ownerID int64) ([]Conversion, error)
	Delete(id int64, ownerID int64) (int64, error)
	GetDailyByCampaign(campaign Campaign) ([]ValueCount, error)
}
