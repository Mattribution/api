package app

import "time"

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

// Tracks handles Track data
type Tracks interface {
	Store(t Track) (int64, error)
}
