package main

type Mission struct {
	ID         string `json:"id"`
	Site       string `json:"site"`
	Pilot      string `json:"pilot"`
	Images     int    `json:"images"`
	BatteryPct int    `json:"battery_pct"`
	Status     string `json:"status"`
}
