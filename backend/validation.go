package main

var missionStatuses = map[string]bool{"planned": true, "flying": true, "landed": true, "aborted": true}

func ValidateMissionStatus(status string) error {
	checker := defaultStatusChecker()
	if err := checker.Check(status); err != nil {
		return err
	}
	return nil
}
