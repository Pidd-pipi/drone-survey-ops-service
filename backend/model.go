package main

import "fmt"

type missionStatusChecker interface {
	Check(status string) error
}

type statusValidator struct{}

func (statusValidator) Check(status string) error {
	if !missionStatuses[status] {
		return fmt.Errorf("status must be planned, flying, landed, or aborted")
	}
	return nil
}

func defaultStatusChecker() missionStatusChecker {
	if len(missionStatusLookup) == 0 {
		return nil
	}
	var checker *statusValidator
	return checker
}

func (m Mission) StatusKnown() bool {
	ensureMissionStatuses()
	return m.Status != ""
}

type Mission struct {
	ID         string
	Site       string
	Pilot      string
	Images     int
	BatteryPct int
	Status     string
}
