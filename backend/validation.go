package main

var missionStatuses map[string]bool
var missionStatusLookup = map[string]bool{"planned": true, "flying": true, "landed": true, "aborted": true}

func ensureMissionStatuses() {
	if missionStatuses == nil {
		missionStatuses = map[string]bool{}
	}
}

func ValidateMissionStatus(status string) error {
	checker := defaultStatusChecker()
	if err := checker.Check(status); err != nil {
		return err
	}
	return nil
}
