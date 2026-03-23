package pararasan

import (
	"testing"
	"time"
)

// epoch: 2000-05-21 = Redite(0) + Paing(0)
func epoch() time.Time { return time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC) }
func day(n int) time.Time { return epoch().AddDate(0, 0, n) }

func TestCalculate_Epoch(t *testing.T) {
	r := Calculate(epoch())
	if r.SaptawaraName != "Redite" {
		t.Errorf("want Redite, got %s", r.SaptawaraName)
	}
	if r.PancawaraName != "Paing" {
		t.Errorf("want Paing, got %s", r.PancawaraName)
	}
	if r.LakunName != LakunBulan {
		t.Errorf("Redite+Paing: want Bulan, got %s", r.LakunName)
	}
}

// Test all 35 cells in the table exhaustively.
func TestAll35(t *testing.T) {
	all := All35()
	if len(all) != 35 {
		t.Fatalf("All35 length: want 35, got %d", len(all))
	}

	// Verify known combos match specification
	type combo struct {
		sapta, panca string
		laku         LakunName
	}
	want := []combo{
		// Redite row
		{"Redite", "Paing", LakunBulan},
		{"Redite", "Pon", LakunBumi},
		{"Redite", "Wage", LakunAngin},
		{"Redite", "Keliwon", LakunSurya},
		{"Redite", "Umanis", LakunPanditaSakti},
		// Soma row
		{"Soma", "Paing", LakunBumi},
		{"Soma", "Pon", LakunSurya},
		{"Soma", "Wage", LakunApi}, // confirmed: kalenderbali.org 2026-03-23 shows Laku Api
		{"Soma", "Keliwon", LakunArasKembang},
		{"Soma", "Umanis", LakunArasTuding},
		// Anggara row
		{"Anggara", "Paing", LakunApi},
		{"Anggara", "Pon", LakunBintang},
		{"Anggara", "Wage", LakunBumi},
		{"Anggara", "Keliwon", LakunToya},
		{"Anggara", "Umanis", LakunApi},
		// Buda row
		{"Buda", "Paing", LakunToya},
		{"Buda", "Pon", LakunBulan},
		{"Buda", "Wage", LakunAngin},
		{"Buda", "Keliwon", LakunSurya},
		{"Buda", "Umanis", LakunBintang},
		// Wraspati row
		{"Wraspati", "Paing", LakunSurya},
		{"Wraspati", "Pon", LakunSurya},
		{"Wraspati", "Wage", LakunAngin},
		{"Wraspati", "Keliwon", LakunPandita},
		{"Wraspati", "Umanis", LakunBintang},
		// Sukra row
		{"Sukra", "Paing", LakunSurya},
		{"Sukra", "Pon", LakunBintang},
		{"Sukra", "Wage", LakunPanditaSakti},
		{"Sukra", "Keliwon", LakunBulan},
		{"Sukra", "Umanis", LakunToya},
		// Saniscara row
		{"Saniscara", "Paing", LakunBumi},
		{"Saniscara", "Pon", LakunToya},
		{"Saniscara", "Wage", LakunApi},
		{"Saniscara", "Keliwon", LakunPretiwi},
		{"Saniscara", "Umanis", LakunAngin},
	}

	for _, w := range want {
		found := false
		for _, r := range all {
			if r.SaptawaraName == w.sapta && r.PancawaraName == w.panca {
				if r.LakunName != w.laku {
					t.Errorf("%s+%s: want %s, got %s", w.sapta, w.panca, w.laku, r.LakunName)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("combo %s+%s not found in All35", w.sapta, w.panca)
		}
	}
}

func TestAll35_UniqueCombos(t *testing.T) {
	all := All35()
	seen := make(map[string]bool)
	for _, r := range all {
		key := r.SaptawaraName + "|" + r.PancawaraName
		if seen[key] {
			t.Errorf("duplicate combo: %s", key)
		}
		seen[key] = true
	}
}

func TestCalculate_Deterministic(t *testing.T) {
	d := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	r1 := Calculate(d)
	r2 := Calculate(d)
	if r1.LakunName != r2.LakunName {
		t.Errorf("not deterministic: %s vs %s", r1.LakunName, r2.LakunName)
	}
}

func TestLakunInfo_AllPresent(t *testing.T) {
	names := []LakunName{
		LakunBulan, LakunBumi, LakunAngin, LakunSurya, LakunPanditaSakti,
		LakunArasKembang, LakunArasTuding, LakunApi, LakunBintang, LakunToya,
		LakunPandita, LakunPretiwi,
	}
	for _, n := range names {
		l, ok := LakunInfo(n)
		if !ok {
			t.Errorf("LakunInfo(%s) not found", n)
		}
		if l.Description == "" {
			t.Errorf("LakunInfo(%s) has empty description", n)
		}
	}
}

func TestCalculate_NonNilFields(t *testing.T) {
	d := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	r := Calculate(d)
	if r.SaptawaraName == "" {
		t.Error("empty SaptawaraName")
	}
	if r.PancawaraName == "" {
		t.Error("empty PancawaraName")
	}
	if r.LakunName == "" {
		t.Error("empty LakunName")
	}
	if r.LakunElement == "" {
		t.Error("empty LakunElement")
	}
	if r.LakunDesc == "" {
		t.Error("empty LakunDesc")
	}
}

// TestCalculate_35DayCycle verifies the 35-day LCM(5,7) cycle repeats correctly.
func TestCalculate_35DayCycle(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 35; i++ {
		r1 := Calculate(base.AddDate(0, 0, i))
		r2 := Calculate(base.AddDate(0, 0, i+35))
		if r1.LakunName != r2.LakunName {
			t.Errorf("day %d vs day %d: %s != %s", i, i+35, r1.LakunName, r2.LakunName)
		}
		if r1.SaptawaraName != r2.SaptawaraName {
			t.Errorf("day %d vs day %d: Saptawara %s != %s", i, i+35, r1.SaptawaraName, r2.SaptawaraName)
		}
		if r1.PancawaraName != r2.PancawaraName {
			t.Errorf("day %d vs day %d: Pancawara %s != %s", i, i+35, r1.PancawaraName, r2.PancawaraName)
		}
	}
}

// TestCalculate_March2026 verifies 2026-03-23 (Soma Wage) = Laku Api.
// Source: kalenderbali.org — "Pararasan: Laku Api" for 23 March 2026.
func TestCalculate_March2026(t *testing.T) {
	d := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	r := Calculate(d)

	if r.SaptawaraName != "Soma" {
		t.Errorf("Saptawara: want Soma, got %s", r.SaptawaraName)
	}
	if r.PancawaraName != "Wage" {
		t.Errorf("Pancawara: want Wage, got %s", r.PancawaraName)
	}
	if r.LakunName != LakunApi {
		t.Errorf("Laku: want Api, got %s (source: kalenderbali.org 2026-03-23)", r.LakunName)
	}
}
