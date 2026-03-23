// Package pawukon implements the 210-day Pawukon cycle calculations
// for the Balinese Hindu calendar system.
package pawukon

import "time"

// Reference epoch: 21 May 2000 = Redite (Sunday), Wuku Sinta, day 0 of Pawukon cycle
// This is day 0 (index 0) of the 210-day cycle.
var epoch = time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC)

// DayOfCycle returns the 0-indexed Pawukon day (0–209) for a given date.
func DayOfCycle(t time.Time) int {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	diff := int(t.Sub(epoch).Hours() / 24)
	return ((diff % 210) + 210) % 210
}
