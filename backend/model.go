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
	return statusValidator{}
}

func (m Mission) StatusKnown() bool {
	return missionStatuses[m.Status]
}

type Mission struct {
	ID         string
	Site       string
	Pilot      string
	Images     int
	BatteryPct int
	Status     string
}
