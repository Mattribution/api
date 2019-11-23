package api

// Track is event tracking data in our format
type Track struct {
	ID              int       `json:"id"`
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

type TrackService interface {
	Track(id int) (*Track, error)
	CreateTrack(t *Track) error
	DeleteTrack(id int) error
}

// // SegmentTrack is event tracking in Segment's format
// type SegmentTrack struct {
// 	ID              int       `json:"id"`
// 	UserID          string    `json:"userId"`
// 	FpHash          string    `json:"fpHash"`   // fingerprint hash
// 	PageURL         string    `json:"pageURL"`  // optional (website specific)
// 	PagePath        string    `json:"pagePath"` // optional ()
// 	PageTitle       string    `json:"pageTitle"`
// 	PageReferrer    string    `json:"pageReferrer"`
// 	Event           string    `json:"event"`
// 	CampaignSource  string    `json:"campaignSource"`
// 	CampaignMedium  string    `json:"campaignMedium"`
// 	CampaignName    string    `json:"campaignName"`
// 	CampaignContent string    `json:"campaignContent"`
// 	SentAt          time.Time `json:"sentAt"`
// 	IP              string
// 	Extra           string `json:"extra"` // (optional) extra json
// }