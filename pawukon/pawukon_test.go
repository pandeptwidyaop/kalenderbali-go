package pawukon

import (
	"testing"
	"time"
)

func TestDayOfCycle_Epoch(t *testing.T) {
	// Epoch day: 21 May 2000 should be day 0
	epoch := time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC)
	got := DayOfCycle(epoch)
	if got != 0 {
		t.Errorf("epoch: want 0, got %d", got)
	}
}

func TestDayOfCycle_CycleBoundary(t *testing.T) {
	// Day 209 (last day)
	d209 := time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 209)
	got := DayOfCycle(d209)
	if got != 209 {
		t.Errorf("day 209: want 209, got %d", got)
	}

	// Day 210 wraps to 0
	d210 := time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 210)
	got = DayOfCycle(d210)
	if got != 0 {
		t.Errorf("day 210 (wrap): want 0, got %d", got)
	}
}

func TestDayOfCycle_BeforeEpoch(t *testing.T) {
	// Day before epoch should be 209
	before := time.Date(2000, 5, 20, 0, 0, 0, 0, time.UTC)
	got := DayOfCycle(before)
	if got != 209 {
		t.Errorf("day before epoch: want 209, got %d", got)
	}
}

func TestDayOfCycle_KnownDate_March2026(t *testing.T) {
	// 2026-03-23 = pawukon day 197 (Soma Wage Dukut)
	// Verified against kalenderbali.org
	d := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	got := DayOfCycle(d)
	if got != 197 {
		t.Errorf("2026-03-23: want pawukon day 197, got %d", got)
	}
}

func TestDayOfCycle_WukuDukut(t *testing.T) {
	// Wuku Dukut starts at day 196 (28*7=196), ends at 202.
	// 2026-03-23 = pd 197 → in Wuku Dukut
	d := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	pd := DayOfCycle(d)
	wuku := pd / 7
	if wuku != 28 { // 0-indexed: Dukut is index 28
		t.Errorf("2026-03-23: want wuku index 28 (Dukut), got %d", wuku)
	}
}

func TestDayOfCycle_Deterministic(t *testing.T) {
	d := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	if DayOfCycle(d) != DayOfCycle(d) {
		t.Error("DayOfCycle is not deterministic")
	}
}

func TestDayOfCycle_Monotonic(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	prev := DayOfCycle(start)
	for i := 1; i < 210; i++ {
		cur := DayOfCycle(start.AddDate(0, 0, i))
		expected := (prev + 1) % 210
		if cur != expected {
			t.Errorf("day %d: want %d, got %d", i, expected, cur)
		}
		prev = cur
	}
}
