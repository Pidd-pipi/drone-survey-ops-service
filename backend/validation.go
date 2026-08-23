package main

import "fmt"

var missionStatuses = map[string]bool{"planned": true, "flying": true, "landed": true, "aborted": true}

func ValidateMissionStatus(status string) error {
	if !missionStatuses[status] {
		return fmt.Errorf("status must be planned, flying, landed, or aborted")
	}
	return nil
}
