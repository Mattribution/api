package app

type Service struct {
	tracksDAO TracksDAO
	kpisDAO   KpisDAO
}

// NewService returns new service object
func NewService(tracksDAO TracksDAO, kpisDAO KpisDAO) Service {
	return Service{
		tracksDAO: tracksDAO,
		kpisDAO:   kpisDAO,
	}
}

func (s Service) NewTrack(t Track) (int64, error) {
	return s.tracksDAO.Store(t)
}

func (s Service) NewKpi(kpi Kpi) (int64, error) {
	return s.kpisDAO.Store(kpi)
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
		kpis[i].CampaignNameJourneyAggregate = aggregate
	}

	// Format
	if kpis == nil {
		kpis = []Kpi{}
	}

	return kpis, nil
}
