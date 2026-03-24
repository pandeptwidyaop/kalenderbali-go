package wewaran

import (
	"fmt"
	"testing"
	"time"
)

// Epoch: 2000-05-21 = day 0 = Redite(0), Paing(0), Wuku Sinta(0)
func epoch() time.Time { return time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC) }
func day(n int) time.Time { return epoch().AddDate(0, 0, n) }

func TestSaptawara_Epoch(t *testing.T) {
	idx, name := Saptawara(epoch())
	if idx != 0 || name != "Redite" {
		t.Errorf("epoch Saptawara: want Redite(0), got %s(%d)", name, idx)
	}
}

func TestSaptawara_Cycle(t *testing.T) {
	names := SaptawaraNames
	for i := 0; i < 21; i++ {
		idx, name := Saptawara(day(i))
		if idx != i%7 || name != names[i%7] {
			t.Errorf("day %d: want %s(%d), got %s(%d)", i, names[i%7], i%7, name, idx)
		}
	}
}

func TestPancawara_Epoch(t *testing.T) {
	idx, name := Pancawara(epoch())
	// Day 0 % 5 = 0 = Paing
	if idx != 0 || name != "Paing" {
		t.Errorf("epoch Pancawara: want Paing(0), got %s(%d)", name, idx)
	}
}

func TestTriwara_Epoch(t *testing.T) {
	idx, name := Triwara(epoch())
	if idx != 0 || name != "Pasah" {
		t.Errorf("epoch Triwara: want Pasah(0), got %s(%d)", name, idx)
	}
}

func TestWuku_Epoch(t *testing.T) {
	idx, name := Wuku(epoch())
	if idx != 0 || name != "Sinta" {
		t.Errorf("epoch Wuku: want Sinta(0), got %s(%d)", name, idx)
	}
}

func TestWuku_Cycle(t *testing.T) {
	for i := 0; i < 30; i++ {
		d := day(i * 7)
		idx, name := Wuku(d)
		if idx != i || name != WukuNames[i] {
			t.Errorf("wuku %d: want %s(%d), got %s(%d)", i, WukuNames[i], i, name, idx)
		}
	}
}

func TestWuku_Watugunung(t *testing.T) {
	// Last wuku starts at day 203 (29*7=203)
	d := day(203)
	idx, name := Wuku(d)
	if idx != 29 || name != "Watugunung" {
		t.Errorf("Watugunung: want idx 29, got %s(%d)", name, idx)
	}
}

func TestCaturwara_SpecialDays(t *testing.T) {
	// Days 71 and 72 should both be Jaya (index 2)
	for _, n := range []int{71, 72} {
		idx, name := Caturwara(day(n))
		if idx != 2 || name != "Jaya" {
			t.Errorf("Caturwara day %d: want Jaya(2), got %s(%d)", n, name, idx)
		}
	}
	// Day 73 should NOT be Jaya (continues after the double)
	idx, _ := Caturwara(day(73))
	if idx == 2 {
		t.Errorf("Caturwara day 73 should not be Jaya again")
	}
}

func TestAstawara_SpecialDays(t *testing.T) {
	// Days 71 and 72 should both be Kala (index 6)
	for _, n := range []int{71, 72} {
		idx, name := Astawara(day(n))
		if idx != 6 || name != "Kala" {
			t.Errorf("Astawara day %d: want Kala(6), got %s(%d)", n, name, idx)
		}
	}
}

func TestSangawara_FirstThreeDays(t *testing.T) {
	// Days 0, 1, 2 should all be Dangu
	for _, n := range []int{0, 1, 2} {
		idx, name := Sangawara(day(n))
		if idx != 0 || name != "Dangu" {
			t.Errorf("Sangawara day %d: want Dangu(0), got %s(%d)", n, name, idx)
		}
	}
	// Day 3 should be Dangu still (3-3=0 → Dangu)
	idx, name := Sangawara(day(3))
	if idx != 0 || name != "Dangu" {
		t.Errorf("Sangawara day 3: want Dangu(0), got %s(%d)", name, idx)
	}
	// Day 4 should be Jangur ((4-3)%9 = 1 = Jangur)
	idx, name = Sangawara(day(4))
	if idx != 1 || name != "Jangur" {
		t.Errorf("Sangawara day 4: want Jangur(1), got %s(%d)", name, idx)
	}
}

func TestUrip_Range(t *testing.T) {
	// Urip should always be 1-10
	for i := 0; i < 210; i++ {
		u := Urip(day(i))
		if u < 1 || u > 18 {
			t.Errorf("day %d: urip %d out of valid range", i, u)
		}
	}
}

func TestDasawara_Range(t *testing.T) {
	for i := 0; i < 210; i++ {
		idx, _ := Dasawara(day(i))
		if idx < 0 || idx > 9 {
			t.Errorf("day %d: Dasawara idx %d out of range", i, idx)
		}
	}
}

func TestCalculate_Completeness(t *testing.T) {
	// All 210 days should produce valid results
	for i := 0; i < 210; i++ {
		w := Calculate(day(i))
		if w.WukuName == "" {
			t.Errorf("day %d: empty WukuName", i)
		}
		if w.SaptawaraName == "" {
			t.Errorf("day %d: empty SaptawaraName", i)
		}
		if w.PancawaraName == "" {
			t.Errorf("day %d: empty PancawaraName", i)
		}
		if w.Urip < 1 || w.Urip > 18 {
			t.Errorf("day %d: invalid Urip %d", i, w.Urip)
		}
	}
}

// TestKnownDates validates against published Balinese calendar references.
// Reference: 2026-03-23 = known to be Redite(Sunday), which aligns with Gregorian weekday.
func TestKnownDates_Saptawara(t *testing.T) {
	tests := []struct {
		date      string
		wantSapta string
	}{
		{"2026-03-23", "Redite"},    // Monday? Let's check: epoch=Sun May 21 2000
		// Actually 2000-05-21 is a Sunday = Redite (index 0). Gregorian Sunday = Redite.
		// 2026-03-23 is a Monday = Soma (index 1). Let's verify.
		{"2026-03-22", "Redite"},    // Sunday
		{"2026-03-23", "Soma"},      // Monday
		{"2026-03-24", "Anggara"},   // Tuesday
		{"2026-03-25", "Buda"},      // Wednesday
		{"2026-03-26", "Wraspati"},  // Thursday
		{"2026-03-27", "Sukra"},     // Friday
		{"2026-03-28", "Saniscara"}, // Saturday
	}
	// Remove the duplicate entry
	tests = []struct {
		date      string
		wantSapta string
	}{
		{"2026-03-22", "Redite"},
		{"2026-03-23", "Soma"},
		{"2026-03-24", "Anggara"},
		{"2026-03-25", "Buda"},
		{"2026-03-26", "Wraspati"},
		{"2026-03-27", "Sukra"},
		{"2026-03-28", "Saniscara"},
	}
	for _, tt := range tests {
		d, _ := time.Parse("2006-01-02", tt.date)
		_, name := Saptawara(d)
		if name != tt.wantSapta {
			t.Errorf("%s: want %s, got %s", tt.date, tt.wantSapta, name)
		}
	}
}

// TestGalungan validates Galungan always falls on Buda(3) Keliwon(3) Dunggulan(10).
func TestGalungan_Wewaran(t *testing.T) {
	// Galungan is pawukon day 73 (0-indexed; spec says day 74 which is 1-indexed)
	d := day(73)
	w := Calculate(d)
	if w.SaptawaraName != "Buda" {
		t.Errorf("Galungan Saptawara: want Buda, got %s", w.SaptawaraName)
	}
	if w.PancawaraName != "Keliwon" {
		t.Errorf("Galungan Pancawara: want Keliwon, got %s", w.PancawaraName)
	}
	if w.WukuName != "Dunggulan" {
		t.Errorf("Galungan Wuku: want Dunggulan, got %s", w.WukuName)
	}
}

// TestKuningan validates Kuningan always falls on Saniscara(6) Keliwon(3) Kuningan(11).
func TestKuningan_Wewaran(t *testing.T) {
	// Kuningan is 10 days after Galungan = day 83 (0-indexed)
	d := day(83)
	w := Calculate(d)
	if w.SaptawaraName != "Saniscara" {
		t.Errorf("Kuningan Saptawara: want Saniscara, got %s", w.SaptawaraName)
	}
	if w.PancawaraName != "Keliwon" {
		t.Errorf("Kuningan Pancawara: want Keliwon, got %s", w.PancawaraName)
	}
	if w.WukuName != "Kuningan" {
		t.Errorf("Kuningan Wuku: want Kuningan, got %s", w.WukuName)
	}
}

// TestSaraswati validates Saraswati is the last day (day 209) = Saniscara Umanis Watugunung.
func TestSaraswati_Wewaran(t *testing.T) {
	d := day(209)
	w := Calculate(d)
	if w.SaptawaraName != "Saniscara" {
		t.Errorf("Saraswati Saptawara: want Saniscara, got %s", w.SaptawaraName)
	}
	if w.PancawaraName != "Umanis" {
		t.Errorf("Saraswati Pancawara: want Umanis, got %s", w.PancawaraName)
	}
	if w.WukuName != "Watugunung" {
		t.Errorf("Saraswati Wuku: want Watugunung, got %s", w.WukuName)
	}
}

// TestKajengKeliwon: Every 15 days (Triwara Kajeng + Pancawara Keliwon)
func TestKajengKeliwon_Period(t *testing.T) {
	// Find first Kajeng Keliwon from epoch
	first := -1
	for i := 0; i < 30; i++ {
		w := Calculate(day(i))
		if w.TriwaraName == "Kajeng" && w.PancawaraName == "Keliwon" {
			first = i
			break
		}
	}
	if first < 0 {
		t.Fatal("no Kajeng Keliwon found in first 30 days")
	}
	// Should repeat every 15 days
	for k := 1; k <= 10; k++ {
		n := first + k*15
		if n >= 210 {
			break
		}
		w := Calculate(day(n))
		if w.TriwaraName != "Kajeng" || w.PancawaraName != "Keliwon" {
			t.Errorf("Kajeng Keliwon expected at day %d (k=%d from %d), got %s %s",
				n, k, first, w.TriwaraName, w.PancawaraName)
		}
	}
}

// TestAllWewaran_March2026 validates all 10 wewaran + extras against
// kalenderbali.org for 2026-03-23 (Soma Wage Dukut).
// Reference: https://m.kalenderbali.org/?bl=3&tg=23&th=2026
// KBD output: "-, Menga, Kajeng, Menala, Wage, Maulu, Soma, Yama, Erangan, Dewa"
func TestAllWewaran_March2026(t *testing.T) {
	d, _ := time.Parse("2006-01-02", "2026-03-23")
	w := Calculate(d)

	checks := []struct {
		name string
		got  string
		want string
	}{
		{"PawukonDay", fmt.Sprintf("%d", w.PawukonDay), "197"},
		{"WukuName", w.WukuName, "Dukut"},
		{"Ekawara", func() string {
			if w.EkawaraName == "" {
				return "–"
			}
			return w.EkawaraName
		}(), "–"},
		{"Dwiwara", w.DwiwaraName, "Menga"},
		{"Triwara", w.TriwaraName, "Kajeng"},
		{"Caturwara", w.CaturwaraName, "Menala"},
		{"Pancawara", w.PancawaraName, "Wage"},
		{"Sadwara", w.SadwaraName, "Maulu"},
		{"Saptawara", w.SaptawaraName, "Soma"},
		{"Astawara", w.AstawaraName, "Yama"},
		{"Sangawara", w.SangawaraName, "Erangan"},
		{"Dasawara", w.DasawaraName, "Dewa"},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("2026-03-23 %s: want %q, got %q", c.name, c.want, c.got)
		}
	}
}

// TestEpochWewaran validates day 0 (first day, 1-indexed "day 1") of the Pawukon cycle.
// Reference: Wikipedia Pawukon table — first day = epoch 2000-05-21.
// Full sequence (Ekawara→Dasawara): -, Menga, Pasah, Sri, Paing, Tungleh, Redite, Sri, Dangu, Sri
func TestEpochWewaran(t *testing.T) {
	w := Calculate(epoch()) // day 0 = first day of cycle

	checks := []struct{ name, got, want string }{
		{"Ekawara", func() string {
			if w.EkawaraName == "" { return "-" }
			return w.EkawaraName
		}(), "-"},
		{"Dwiwara", w.DwiwaraName, "Menga"},
		{"Triwara", w.TriwaraName, "Pasah"},
		{"Caturwara", w.CaturwaraName, "Sri"},
		{"Pancawara", w.PancawaraName, "Paing"},
		{"Sadwara", w.SadwaraName, "Tungleh"},
		{"Saptawara", w.SaptawaraName, "Redite"},
		{"Astawara", w.AstawaraName, "Sri"},
		{"Sangawara", w.SangawaraName, "Dangu"},
		{"Dasawara", w.DasawaraName, "Sri"},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("epoch (day0) %s: want %q, got %q", c.name, c.want, c.got)
		}
	}
}
