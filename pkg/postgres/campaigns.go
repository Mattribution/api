package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/mattribution/api/pkg/api"
)

type CampaignService struct {
	DB *sqlx.DB
}

// Store stores the track in the db
func (s CampaignService) Store(campaign api.Campaign) (int, error) {
	sqlStatement :=
		`INSERT INTO public.campaigns (owner_id, name, column_name, column_value, cost_per_month, created_at)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id`

	id := 0
	err := s.DB.QueryRow(sqlStatement, campaign.OwnerID, campaign.Name, campaign.ColumnName, campaign.ColumnValue, campaign.CostPerMonth, time.Now().Format(time.RFC3339)).Scan(&id)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

	return id, err
}

func (s CampaignService) Update(campaign api.Campaign) error {
	sqlStatement :=
		fmt.Sprintf(`UPDATE public.campaigns
		SET name=:name, cost_per_month=:cost_per_month
		WHERE id=%v`, campaign.ID)

	_, err := s.DB.NamedExec(sqlStatement, campaign)

	return err
}

func (s CampaignService) Find(ownerID int) ([]api.Campaign, error) {
	sqlStatement :=
		`SELECT * FROM public.campaigns
		WHERE owner_id = $1`

	campaigns := []api.Campaign{}
	err := s.DB.Select(&campaigns, sqlStatement, ownerID)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}

	return campaigns, err
}

// FindByID finds a single campaign by id
func (s CampaignService) FindByID(id int, ownerID int) (api.Campaign, error) {
	sqlStatement :=
		`SELECT * FROM public.campaigns
		WHERE id = $1
		AND owner_id = $2`

	var campaign api.Campaign

	switch err := s.DB.Get(&campaign, sqlStatement, id, ownerID); err {
	case sql.ErrNoRows:
		return api.Campaign{}, errors.New("Not found")
	case nil:
	default:
		return api.Campaign{}, err
	}

	return campaign, nil
}

// ScanForNewCampaigns scans the tracks for any new campaigns and creates new
//  campaings if it detects a pattern that hasn't been matched.
//  TODO: optimize this to run on a schedule
func (s CampaignService) ScanForNewCampaigns(ownerID int) (int, error) {
	sqlStatement := `SELECT campaign_name FROM tracks
	WHERE campaign_name <> ''
	AND NOT EXISTS (
		SELECT 1
		FROM campaigns
		WHERE owner_id = $1
		AND column_name = 'campaign_name'
		AND column_value = campaign_name
	)
	GROUP BY 1`

	var newCampaignNames []string
	err := s.DB.Select(&newCampaignNames, sqlStatement, ownerID)
	if err != nil {
		return 0, err
	}

	// Loop over detected campaign names and create campaigns for them
	storedCount := 0
	for _, campaignName := range newCampaignNames {
		newCampaign := api.Campaign{
			OwnerID:     ownerID,
			Name:        campaignName,
			ColumnName:  "campaign_name",
			ColumnValue: campaignName,
			CreatedAt:   time.Now(),
		}
		_, err := s.Store(newCampaign)
		if err != nil {
			return storedCount, err
		}
		storedCount++
	}

	return storedCount, nil
}
