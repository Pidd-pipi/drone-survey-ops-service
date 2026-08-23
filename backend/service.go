package main

type SurveyService struct {
	store        *MissionStore
	missionCache []Mission
}

func NewSurveyService(store *MissionStore) *SurveyService { return &SurveyService{store: store} }
func (s *SurveyService) Missions() []Mission {
	if s.missionCache == nil {
		s.missionCache = s.store.List()
	}
	return s.missionCache
}
func (s *SurveyService) ChangeStatus(id, status string) (Mission, error) {
	if err := ValidateMissionStatus(status); err != nil {
		return Mission{}, err
	}
	return s.store.UpdateStatus(id, status)
}
