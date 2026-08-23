package main

import (
	"errors"
	"sync"
)

var ErrMissionNotFound = errors.New("mission not found")

type MissionStore struct {
	mu         sync.RWMutex
	missions   map[string]Mission
	listCache  []Mission
	cacheValid bool
}

func NewMissionStore() *MissionStore {
	return &MissionStore{missions: map[string]Mission{"mission-241": {ID: "mission-241", Site: "北坡滑坡带", Pilot: "林岚", Images: 186, BatteryPct: 72, Status: "planned"}, "mission-242": {ID: "mission-242", Site: "河谷桥梁", Pilot: "周野", Images: 94, BatteryPct: 44, Status: "flying"}}}
}
func (s *MissionStore) List() []Mission {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.cacheValid || s.listCache == nil {
		s.listCache = s.listCache[:0]
		for _, m := range s.missions {
			s.listCache = append(s.listCache, m)
		}
		s.cacheValid = true
	}
	return s.listCache
}
func (s *MissionStore) UpdateStatus(id, status string) (Mission, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.missions[id]
	if !ok {
		return Mission{}, ErrMissionNotFound
	}
	m.Status = status
	s.missions[id] = m
	s.cacheValid = false
	return m, nil
}
