package main

type SurveyService struct{ store *MissionStore }

func NewSurveyService(store *MissionStore) *SurveyService { return &SurveyService{store: store} }
func (s *SurveyService) Missions() []Mission              { return s.store.List() }
func (s *SurveyService) ChangeStatus(id, status string) (Mission, error) {
	if err := ValidateMissionStatus(status); err != nil {
		return Mission{}, err
	}
	return s.store.UpdateStatus(id, status)
}
