package app

import (
	"errors"
	"log"
)

const (
	DefaultModelIDValue = "first-touch"
)

type Service struct {
	tracksDAO TracksDAO
	kpisDAO   KpisDAO
	usersDAO  UsersDAO
}

// NewService returns new service object
func NewService(tracksDAO TracksDAO, kpisDAO KpisDAO, usersDAO UsersDAO) Service {
	return Service{
		tracksDAO: tracksDAO,
		kpisDAO:   kpisDAO,
		usersDAO:  usersDAO,
	}
}

func (s Service) NewTrack(t Track, ownerSecret string) (int64, error) {
	users, err := s.usersDAO.FindBySecret(ownerSecret)
	if err != nil {
		return 0, err
	}

	// TODO: Make this print a 4xx error instead of flowing up to a 500
	// this will probably involve creating a custom error object
	if len(users) == 0 {
		return 0, errors.New("No user was found for that secret")
	}

	if len(users) > 1 {
		errStr := "Found multiple users for one secret key"
		// Note: This error is serious af... idk how this could happen
		log.Println(errStr)
		return 0, errors.New(errStr)
	}
	user := users[0]

	t.OwnerID = user.UUID

	return s.tracksDAO.Store(t)
}

func (s Service) NewKpi(kpi Kpi) (int64, error) {
	if kpi.ModelID == "" {
		kpi.ModelID = DefaultModelIDValue
	}
	return s.kpisDAO.Store(kpi)
}

func (s Service) UpdateKpi(kpi Kpi) error {
	return s.kpisDAO.Update(kpi)
}

func (s Service) DeleteKpi(kpi Kpi) (int64, error) {
	return s.kpisDAO.Delete(kpi.ID, kpi.OwnerID)
}

func (s Service) GetKpisForUser(ownerID string) ([]Kpi, error) {
	kpis, err := s.kpisDAO.FindByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}

	// Get aggregates for the kpi
	for i, kpi := range kpis {
		// Get aggregate data
		aggregate, err := s.tracksDAO.GetNormalizedJourneyAggregate(kpi.OwnerID, "campaign_name", kpi.PatternMatchColumnName, kpi.PatternMatchRowValue)
		if err != nil {
			return nil, err
		}
		if aggregate == nil {
			aggregate = []PosAggregate{}
		}
		kpis[i].CampaignNameJourneyAggregate = aggregate
	}

	// Format
	if kpis == nil {
		kpis = []Kpi{}
	}

	return kpis, nil
}
