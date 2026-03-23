package lunar

import (
	"testing"
	"time"
)

func TestNewMoonPhases_2026(t *testing.T) {
	phases := PhasesInYear(2026)
	var newMoons, fullMoons []MoonPhase
	for _, p := range phases {
		if p.Phase == NewMoon {
			newMoons = append(newMoons, p)
		} else {
			fullMoons = append(fullMoons, p)
		}
	}
	if len(newMoons) < 12 || len(newMoons) > 13 {
		t.Errorf("expected 12-13 new moons in 2026, got %d", len(newMoons))
	}
	if len(fullMoons) < 12 || len(fullMoons) > 13 {
		t.Errorf("expected 12-13 full moons in 2026, got %d", len(fullMoons))
	}
}

func TestPhasesInYear_AllInYear(t *testing.T) {
	phases := PhasesInYear(2026)
	for _, p := range phases {
		if p.Date.Year() != 2026 {
			t.Errorf("phase date %s not in 2026", p.Date.Format("2006-01-02"))
		}
	}
}

func TestPhasesInYear_Sorted(t *testing.T) {
	phases := PhasesInYear(2026)
	for i := 1; i < len(phases); i++ {
		if phases[i].Date.Before(phases[i-1].Date) {
			t.Errorf("phases not sorted: %s before %s",
				phases[i].Date.Format("2006-01-02"),
				phases[i-1].Date.Format("2006-01-02"))
		}
	}
}

// TestKnownFullMoons validates against astronomically confirmed dates.
// Source: US Naval Observatory / timeanddate.com
func TestKnownFullMoons(t *testing.T) {
	knownFullMoons := []string{
		"2026-01-03",
		"2026-02-01",
		"2026-03-03",
		"2026-04-02",
		"2026-05-01",
		"2026-05-31",
		"2026-06-30",
		"2026-07-29",
		"2026-08-28",
		"2026-09-26",
		"2026-10-26",
		"2026-11-24",
		"2026-12-24",
	}

	phases := PhasesInYear(2026)
	fullMoonDates := make(map[string]bool)
	for _, p := range phases {
		if p.Phase == FullMoon {
			fullMoonDates[p.Date.Format("2006-01-02")] = true
		}
	}

	for _, want := range knownFullMoons {
		if !fullMoonDates[want] {
			// Allow ±1 day tolerance for algorithm accuracy
			d, _ := time.Parse("2006-01-02", want)
			prev := d.AddDate(0, 0, -1).Format("2006-01-02")
			next := d.AddDate(0, 0, 1).Format("2006-01-02")
			if !fullMoonDates[prev] && !fullMoonDates[next] {
				t.Errorf("expected full moon near %s, not found (±1 day)", want)
			}
		}
	}
}

// TestKnownNewMoons validates new moon dates for 2026.
func TestKnownNewMoons(t *testing.T) {
	knownNewMoons := []string{
		"2026-01-18",
		"2026-02-17",
		"2026-03-19",
		"2026-04-17",
		"2026-05-16",
		"2026-06-15",
		"2026-07-14",
		"2026-08-12",
		"2026-09-11",
		"2026-10-10",
		"2026-11-09",
		"2026-12-09",
	}

	phases := PhasesInYear(2026)
	newMoonDates := make(map[string]bool)
	for _, p := range phases {
		if p.Phase == NewMoon {
			newMoonDates[p.Date.Format("2006-01-02")] = true
		}
	}

	for _, want := range knownNewMoons {
		if !newMoonDates[want] {
			d, _ := time.Parse("2006-01-02", want)
			prev := d.AddDate(0, 0, -1).Format("2006-01-02")
			next := d.AddDate(0, 0, 1).Format("2006-01-02")
			if !newMoonDates[prev] && !newMoonDates[next] {
				t.Errorf("expected new moon near %s, not found (±1 day)", want)
			}
		}
	}
}

func TestNextPhase_FullMoon(t *testing.T) {
	// From 2026-03-01, next full moon should be around 2026-03-03
	start, _ := time.Parse("2006-01-02", "2026-03-01")
	p := NextPhase(start, FullMoon)
	if p.Phase != FullMoon {
		t.Errorf("wrong phase type: got %v", p.Phase)
	}
	if p.Date.Before(start) {
		t.Errorf("next full moon %s is before start %s", p.Date.Format("2006-01-02"), start.Format("2006-01-02"))
	}
	// Should be within 30 days
	if p.Date.After(start.AddDate(0, 0, 30)) {
		t.Errorf("next full moon %s is too far ahead", p.Date.Format("2006-01-02"))
	}
}

func TestNextPhase_NewMoon(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2026-03-15")
	p := NextPhase(start, NewMoon)
	if p.Phase != NewMoon {
		t.Errorf("wrong phase type: got %v", p.Phase)
	}
	if p.Date.Before(start) {
		t.Errorf("next new moon %s is before start", p.Date.Format("2006-01-02"))
	}
}

func TestPhasesBetween_CrossYear(t *testing.T) {
	start, _ := time.Parse("2006-01-02", "2025-12-01")
	end, _ := time.Parse("2006-01-02", "2026-02-01")
	phases := PhasesBetween(start, end)
	for _, p := range phases {
		if p.Date.Before(start) || !p.Date.Before(end) {
			t.Errorf("phase %s outside [%s, %s)",
				p.Date.Format("2006-01-02"),
				start.Format("2006-01-02"),
				end.Format("2006-01-02"))
		}
	}
	if len(phases) < 3 {
		t.Errorf("expected at least 3 phases in 2-month window, got %d", len(phases))
	}
}

func TestPhaseType_String(t *testing.T) {
	if NewMoon.String() != "Tilem" {
		t.Errorf("NewMoon string: want Tilem, got %s", NewMoon.String())
	}
	if FullMoon.String() != "Purnama" {
		t.Errorf("FullMoon string: want Purnama, got %s", FullMoon.String())
	}
}
