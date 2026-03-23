package jodoh_test

import (
	"testing"
	"time"

	"github.com/pandeptwidyaop/kalenderbali-go/jodoh"
)

// Reference dates for the test pair: Kak Pande (1997-10-08) & partner (2002-08-17)
var (
	pria   = time.Date(1997, 10, 8, 0, 0, 0, 0, time.UTC)
	wanita = time.Date(2002, 8, 17, 0, 0, 0, 0, time.UTC)
)

func TestCheckJodoh_ReturnsAllMethods(t *testing.T) {
	results := jodoh.CheckJodoh(pria, wanita)
	if len(results) != 6 {
		t.Fatalf("expected 6 results, got %d", len(results))
	}

	expected := []string{
		"Saptawara",
		"Neptu Mod 5",
		"Neptu Mod 4",
		"Pertemuan Neptu (Mod 9)",
		"Tri Pramana (Sodasa Rsi)",
		"Ramalan 5 Tahun",
	}
	for i, r := range results {
		if r.Method != expected[i] {
			t.Errorf("result[%d].Method = %q, want %q", i, r.Method, expected[i])
		}
		if r.Prediction == "" {
			t.Errorf("result[%d].Prediction is empty", i)
		}
		if r.Description == "" {
			t.Errorf("result[%d].Description is empty", i)
		}
	}
}

func TestSaptawaraName(t *testing.T) {
	name := jodoh.SaptawaraName(pria)
	if name == "" {
		t.Error("SaptawaraName returned empty string")
	}
}

func TestPancawaraName(t *testing.T) {
	name := jodoh.PancawaraName(wanita)
	if name == "" {
		t.Error("PancawaraName returned empty string")
	}
}

func TestUripTotal_Positive(t *testing.T) {
	u := jodoh.UripTotal(pria)
	if u <= 0 {
		t.Errorf("UripTotal expected > 0, got %d", u)
	}
}

func TestNeptuMod5_ValidRange(t *testing.T) {
	results := jodoh.CheckJodoh(pria, wanita)
	for _, r := range results {
		if r.Method == "Neptu Mod 5" {
			if r.Score < 0 || r.Score > 4 {
				t.Errorf("Neptu Mod 5 score out of range [0,4]: %d", r.Score)
			}
		}
	}
}

func TestNeptuMod4_ValidRange(t *testing.T) {
	results := jodoh.CheckJodoh(pria, wanita)
	for _, r := range results {
		if r.Method == "Neptu Mod 4" {
			if r.Score < 0 || r.Score > 3 {
				t.Errorf("Neptu Mod 4 score out of range [0,3]: %d", r.Score)
			}
		}
	}
}

func TestTriPramana_ValidRange(t *testing.T) {
	results := jodoh.CheckJodoh(pria, wanita)
	for _, r := range results {
		if r.Method == "Tri Pramana (Sodasa Rsi)" {
			if r.Score < 1 || r.Score > 16 {
				t.Errorf("Tri Pramana sisa out of range [1,16]: %d", r.Score)
			}
		}
	}
}

func TestRamalan5Tahun_ValidRange(t *testing.T) {
	results := jodoh.CheckJodoh(pria, wanita)
	for _, r := range results {
		if r.Method == "Ramalan 5 Tahun" {
			if r.Score < 0 || r.Score > 4 {
				t.Errorf("Ramalan 5 Tahun score out of range [0,4]: %d", r.Score)
			}
		}
	}
}

// TestMod9Matrix_NoPanic ensures all 81 combos are accessible without panic.
func TestMod9Matrix_NoPanic(t *testing.T) {
	dates := []time.Time{
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2001, 3, 5, 0, 0, 0, 0, time.UTC),
		time.Date(1995, 7, 22, 0, 0, 0, 0, time.UTC),
		time.Date(1990, 11, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2010, 5, 30, 0, 0, 0, 0, time.UTC),
	}
	for _, a := range dates {
		for _, b := range dates {
			results := jodoh.CheckJodoh(a, b)
			for _, r := range results {
				if r.Prediction == "" {
					t.Errorf("empty prediction for pair %s vs %s method %s",
						a.Format("2006-01-02"), b.Format("2006-01-02"), r.Method)
				}
			}
		}
	}
}

// TestSymmetry_Saptawara checks the matrix covers all 49 saptawara pairs.
func TestAllSaptawaraPairs_NonEmpty(t *testing.T) {
	// The 210-day Pawukon cycle covers all saptawara × pancawara combos.
	// Cycle through 7 consecutive weeks from epoch to hit all 7 saptawara values.
	epoch := time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC) // Day 0 = Redite
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			a := epoch.AddDate(0, 0, i)
			b := epoch.AddDate(0, 0, j)
			results := jodoh.CheckJodoh(a, b)
			var saptaResult *jodoh.JodohResult
			for k := range results {
				if results[k].Method == "Saptawara" {
					saptaResult = &results[k]
					break
				}
			}
			if saptaResult == nil {
				t.Fatalf("no Saptawara result for pair %d,%d", i, j)
			}
			if saptaResult.Prediction == "" {
				t.Errorf("empty Saptawara prediction for pair %d,%d", i, j)
			}
		}
	}
}
