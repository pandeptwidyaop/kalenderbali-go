package hariraya

import (
	"testing"
	"time"
)

func TestHolidaysInYear_HasGalungan(t *testing.T) {
	list := HolidaysInYear(2026)
	found := false
	for _, h := range list {
		if h.Name == "Galungan" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Galungan not found in 2026")
	}
}

func TestHolidaysInYear_HasKuningan(t *testing.T) {
	list := HolidaysInYear(2026)
	found := false
	for _, h := range list {
		if h.Name == "Kuningan" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Kuningan not found in 2026")
	}
}

func TestHolidaysInYear_GalunganKuningan10Days(t *testing.T) {
	list := HolidaysInYear(2026)
	var galunganDates, kuninganDates []time.Time
	for _, h := range list {
		if h.Name == "Galungan" {
			galunganDates = append(galunganDates, h.Date)
		}
		if h.Name == "Kuningan" {
			kuninganDates = append(kuninganDates, h.Date)
		}
	}
	if len(galunganDates) == 0 {
		t.Fatal("no Galungan found")
	}
	if len(galunganDates) != len(kuninganDates) {
		t.Errorf("Galungan count %d != Kuningan count %d", len(galunganDates), len(kuninganDates))
	}
	for i, g := range galunganDates {
		if i >= len(kuninganDates) {
			break
		}
		k := kuninganDates[i]
		diff := int(k.Sub(g).Hours() / 24)
		if diff != 10 {
			t.Errorf("Kuningan %s is not 10 days after Galungan %s (diff=%d)",
				k.Format("2006-01-02"), g.Format("2006-01-02"), diff)
		}
	}
}

func TestHolidaysInYear_HasSaraswati(t *testing.T) {
	list := HolidaysInYear(2026)
	found := false
	for _, h := range list {
		if h.Name == "Saraswati" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Saraswati not found in 2026")
	}
}

func TestHolidaysInYear_HasPagerdwesi(t *testing.T) {
	list := HolidaysInYear(2026)
	found := false
	for _, h := range list {
		if h.Name == "Pagerwesi" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Pagerwesi not found in 2026")
	}
}

func TestHolidaysInYear_HasPurnama(t *testing.T) {
	list := HolidaysInYear(2026)
	count := 0
	for _, h := range list {
		if len(h.Name) >= 7 && h.Name[:7] == "Purnama" {
			count++
		}
	}
	if count < 12 || count > 13 {
		t.Errorf("expected 12-13 Purnama in 2026, got %d", count)
	}
}

func TestHolidaysInYear_HasTilem(t *testing.T) {
	list := HolidaysInYear(2026)
	count := 0
	for _, h := range list {
		if len(h.Name) >= 6 && h.Name[:6] == "Tilem " {
			count++
		}
	}
	if count < 12 || count > 13 {
		t.Errorf("expected 12-13 Tilem in 2026, got %d", count)
	}
}

func TestHolidaysInYear_Sorted(t *testing.T) {
	list := HolidaysInYear(2026)
	for i := 1; i < len(list); i++ {
		if list[i].Date.Before(list[i-1].Date) {
			t.Errorf("not sorted: %s before %s",
				list[i].Date.Format("2006-01-02"),
				list[i-1].Date.Format("2006-01-02"))
		}
	}
}

func TestHolidaysInYear_AllInYear(t *testing.T) {
	list := HolidaysInYear(2026)
	for _, h := range list {
		if h.Date.Year() != 2026 {
			t.Errorf("holiday %s date %s not in 2026", h.Name, h.Date.Format("2006-01-02"))
		}
	}
}

func TestForDate_Today(t *testing.T) {
	// Just shouldn't panic
	today := time.Now()
	_ = ForDate(today)
}

func TestNextHoliday_Galungan(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	h := NextHoliday(start, "galungan")
	if h == nil {
		t.Fatal("NextHoliday galungan returned nil")
	}
	if h.Date.Before(start) {
		t.Errorf("Galungan %s is before start", h.Date.Format("2006-01-02"))
	}
}

func TestNextN_ReturnsN(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	list := NextN(start, 10)
	if len(list) != 10 {
		t.Errorf("NextN(10): want 10, got %d", len(list))
	}
	for _, h := range list {
		if h.Date.Before(start) {
			t.Errorf("NextN: holiday %s is before start", h.Date.Format("2006-01-02"))
		}
	}
}

func TestSearch_Pagerwesi(t *testing.T) {
	list := Search("pagerwesi", 2026)
	if len(list) == 0 {
		t.Error("Search pagerwesi 2026 returned empty")
	}
}

func TestHolidaysInYear_TumpekCount(t *testing.T) {
	list := HolidaysInYear(2026)
	count := 0
	for _, h := range list {
		if len(h.Name) >= 6 && h.Name[:6] == "Tumpek" {
			count++
		}
	}
	// Each 210-day cycle has 5 Tumpek; ~1.7 cycles per year → expect ~7-9
	if count < 5 || count > 12 {
		t.Errorf("unexpected Tumpek count: %d (want 5-12)", count)
	}
}

func TestHolidaysInYear_HasNyepi(t *testing.T) {
	// Nyepi occurs once per year (Tahun Baru Saka)
	list := HolidaysInYear(2026)
	count := 0
	for _, h := range list {
		if h.Name == "Nyepi" {
			count++
		}
	}
	if count < 1 || count > 2 {
		t.Errorf("Nyepi count: want 1, got %d", count)
	}
}

// TestMarch2026_Holidays validates specific March 2026 dates against kalenderbali.org.
func TestMarch2026_Holidays(t *testing.T) {
	list := HolidaysBetween(
		time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
	)

	// Build date→names map
	byDate := make(map[string][]string)
	for _, h := range list {
		key := h.Date.Format("2006-01-02")
		byDate[key] = append(byDate[key], h.Name)
	}

	assertHas := func(date, name string) {
		t.Helper()
		for _, n := range byDate[date] {
			if n == name {
				return
			}
		}
		t.Errorf("%s: want %q, got %v", date, name, byDate[date])
	}

	// 3 Mar: Purnama
	assertHas("2026-03-02", "Purnama Kasanga")
	// 14 Mar: Kajeng Keliwon + Tumpek Wayang
	assertHas("2026-03-14", "Kajeng Keliwon")
	assertHas("2026-03-14", "Tumpek Wayang")
	// 18 Mar: Buda Wage
	assertHas("2026-03-18", "Buda Wage")
	// 19 Mar: Tilem
	assertHas("2026-03-18", "Tilem Kasanga")
	// 20 Mar: Nyepi (day after Tilem Kasanga)
	assertHas("2026-03-19", "Nyepi")
	// 21 Mar: Ngembak Geni (day after Nyepi)
	assertHas("2026-03-20", "Ngembak Geni")
	// 24 Mar: Anggara Kasih
	assertHas("2026-03-24", "Anggara Kasih")
	// 29 Mar: Kajeng Keliwon
	assertHas("2026-03-29", "Kajeng Keliwon")
}
