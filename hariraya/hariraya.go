// Package hariraya calculates Hindu Balinese holy days (Hari Raya).
// Holy days are derived from Pawukon positions (fixed in 210-day cycle) and
// lunar calculations (Purnama/Tilem).
package hariraya

import (
	"sort"
	"strings"
	"time"

	"github.com/pandeptwidyaop/kalenderbali-go/lunar"
	"github.com/pandeptwidyaop/kalenderbali-go/pawukon"
	"github.com/pandeptwidyaop/kalenderbali-go/wewaran"
)

// HariRaya represents a Balinese holy day event.
type HariRaya struct {
	Date        time.Time
	Name        string
	Description string
	Category    string // "pawukon", "lunar", "conjunction"
}

// pawukon-based holy days: (0-indexed pawukon day, name, description)
// Pawukon days are 0-indexed (0–209). The original spec uses 1-indexed days,
// so all values below are spec_day - 1.
// Verified:
//   Pagerwesi   day  3: Buda(3%7=3)  Keliwon(3%5=3)  Sinta(3/7=0)     ✓
//   TumpekLandep day 13: Saniscara(13%7=6) Keliwon(13%5=3) Landep(13/7=1) ✓
//   TumpekUduh  day 48: Saniscara(48%7=6) Keliwon(48%5=3) Wariga(48/7=6)  ✓
//   Galungan    day 73: Buda(73%7=3)  Keliwon(73%5=3) Dunggulan(73/7=10) ✓
//   Kuningan    day 83: Saniscara(83%7=6) Keliwon(83%5=3) Kuningan(83/7=11) ✓
//   TumpekKrulut day118: Saniscara(118%7=6) Keliwon(118%5=3) Krulut(118/7=16) ✓
//   TumpekKandang day153: Saniscara(153%7=6) Keliwon(153%5=3) Uye(153/7=21) ✓
//   TumpekWayang day188: Saniscara(188%7=6) Keliwon(188%5=3) Wayang(188/7=26) ✓
//   Saraswati   day209: Saniscara(209%7=6) Umanis(209%5=4) Watugunung(209/7=29) ✓
var pawukonHolidays = []struct {
	Day         int
	Name        string
	Description string
}{
	{3, "Pagerwesi", "Buda Kliwon Sinta — Hari memagari diri dari pengaruh negatif"},
	{13, "Tumpek Landep", "Saniscara Kliwon Landep — Hari memberkati pusaka/benda tajam"},
	{48, "Tumpek Uduh", "Saniscara Kliwon Wariga — Hari memberkati pohon/tanaman"},
	{73, "Galungan", "Buda Keliwon Dunggulan — Hari kemenangan dharma atas adharma"},
	{83, "Kuningan", "Saniscara Kliwon Kuningan — Hari kemenangan, roh leluhur kembali ke surga"},
	{118, "Tumpek Krulut", "Saniscara Kliwon Krulut — Hari memberkati alat musik/kesenian"},
	{153, "Tumpek Kandang", "Saniscara Kliwon Uye — Hari memberkati hewan peliharaan/ternak"},
	{188, "Tumpek Wayang", "Saniscara Kliwon Wayang — Hari memberkati seni pewayangan"},
	{209, "Saraswati", "Saniscara Umanis Watugunung — Hari dewi ilmu pengetahuan"},
	// Special: Banyu Pinaruh is day after Saraswati (day 0 of next cycle)
	// Handled separately below
}

// HolidaysInYear returns all Balinese holy days in the given Gregorian year.
func HolidaysInYear(year int) []HariRaya {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	return HolidaysBetween(start, end)
}

// HolidaysBetween returns all holy days between start (inclusive) and end (exclusive).
func HolidaysBetween(start, end time.Time) []HariRaya {
	var results []HariRaya

	// ── 1. Pawukon-based holy days ──────────────────────────────────────────
	// Scan day by day from a bit before start to catch the first occurrence
	// Start scanning from far enough back to include the first holiday cycle
	scan := start.AddDate(0, 0, -210) // go back one full cycle

	for d := scan; d.Before(end); d = d.AddDate(0, 0, 1) {
		pd := pawukon.DayOfCycle(d)

		for _, h := range pawukonHolidays {
			if pd == h.Day {
				if !d.Before(start) {
					results = append(results, HariRaya{
						Date:        d,
						Name:        h.Name,
						Description: h.Description,
						Category:    "pawukon",
					})

					// Add derived holidays for Galungan
					if h.Name == "Galungan" {
						penampahan := d.AddDate(0, 0, -1)
						manis := d.AddDate(0, 0, 1)
						if !penampahan.Before(start) && penampahan.Before(end) {
							results = append(results, HariRaya{
								Date:        penampahan,
								Name:        "Penampahan Galungan",
								Description: "Sehari sebelum Galungan — persiapan upacara",
								Category:    "pawukon",
							})
						}
						if !manis.Before(start) && manis.Before(end) {
							results = append(results, HariRaya{
								Date:        manis,
								Name:        "Manis Galungan",
								Description: "Sehari setelah Galungan — hari bersilaturahmi",
								Category:    "pawukon",
							})
						}
					}

					// Add derived holidays for Kuningan
					if h.Name == "Kuningan" {
						manis := d.AddDate(0, 0, 1)
						if !manis.Before(start) && manis.Before(end) {
							results = append(results, HariRaya{
								Date:        manis,
								Name:        "Manis Kuningan",
								Description: "Sehari setelah Kuningan",
								Category:    "pawukon",
							})
						}
					}

					// Add Banyu Pinaruh after Saraswati
					if h.Name == "Saraswati" {
						banyuPinaruh := d.AddDate(0, 0, 1)
						if !banyuPinaruh.Before(start) && banyuPinaruh.Before(end) {
							results = append(results, HariRaya{
								Date:        banyuPinaruh,
								Name:        "Banyu Pinaruh",
								Description: "Sehari setelah Saraswati — mandi/membersihkan diri di sumber air",
								Category:    "pawukon",
							})
						}
					}
				}
			}
		}
	}

	// ── 2. Lunar-based holy days ────────────────────────────────────────────
	phases := lunar.PhasesBetween(start.AddDate(0, -2, 0), end.AddDate(0, 2, 0))

	for _, p := range phases {
		if p.Date.Before(start) || !p.Date.Before(end) {
			continue
		}

		sasihIdx, sasihName := lunar.SasihForDate(p.Date)

		if p.Phase == lunar.FullMoon {
			// Purnama for each Sasih
			results = append(results, HariRaya{
				Date:        p.Date,
				Name:        "Purnama " + sasihName,
				Description: "Bulan purnama Sasih " + sasihName,
				Category:    "lunar",
			})
		} else {
			// Tilem for each Sasih
			results = append(results, HariRaya{
				Date:        p.Date,
				Name:        "Tilem " + sasihName,
				Description: "Bulan mati Sasih " + sasihName,
				Category:    "lunar",
			})

			// Siwaratri = Tilem Kapitu (sasih index 6)
			if sasihIdx == 6 {
				results = append(results, HariRaya{
					Date:        p.Date,
					Name:        "Siwaratri",
					Description: "Tilem Kapitu — malam pemujaan Dewa Siwa, malam terpanjang",
					Category:    "lunar",
				})
			}

			// Nyepi = day after Tilem Kasanga (sasih index 8)
			// Ngembak Geni = day after Nyepi (2 days after Tilem Kasanga)
			if sasihIdx == 8 {
				nyepi := p.Date.AddDate(0, 0, 1)
				ngembak := p.Date.AddDate(0, 0, 2)
				if !nyepi.Before(start) && nyepi.Before(end) {
					results = append(results, HariRaya{
						Date:        nyepi,
						Name:        "Nyepi",
						Description: "Tahun Baru Saka — hari keheningan, catur brata penyepian",
						Category:    "lunar",
					})
				}
				if !ngembak.Before(start) && ngembak.Before(end) {
					results = append(results, HariRaya{
						Date:        ngembak,
						Name:        "Ngembak Geni",
						Description: "Sehari setelah Nyepi — hari silaturahmi dan menyalakan api kembali",
						Category:    "lunar",
					})
				}
			}
		}
	}

	// ── 3. Special conjunctions (daily scan) ────────────────────────────────
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		w := wewaran.Calculate(d)

		// Kajeng Keliwon: Triwara Kajeng (index 2) + Pancawara Keliwon (index 3)
		if w.TriwaraIndex == 2 && w.PancawaraIndex == 3 {
			results = append(results, HariRaya{
				Date:        d,
				Name:        "Kajeng Keliwon",
				Description: "Kajeng Keliwon — hari berdoa/penguatan spiritual, setiap 15 hari",
				Category:    "conjunction",
			})
		}

		// Buda Wage: Saptawara Buda (index 3) + Pancawara Wage (index 2)
		if w.SaptawaraIndex == 3 && w.PancawaraIndex == 2 {
			results = append(results, HariRaya{
				Date:        d,
				Name:        "Buda Wage",
				Description: "Buda Wage — hari keberuntungan dalam berdagang",
				Category:    "conjunction",
			})
		}

		// Anggara Kasih: Saptawara Anggara (index 2) + Pancawara Keliwon (index 3)
		if w.SaptawaraIndex == 2 && w.PancawaraIndex == 3 {
			results = append(results, HariRaya{
				Date:        d,
				Name:        "Anggara Kasih",
				Description: "Anggara Kasih — hari berbakti kepada leluhur",
				Category:    "conjunction",
			})
		}
	}

	// Remove duplicates and sort
	results = deduplicate(results)
	sort.Slice(results, func(i, j int) bool {
		if results[i].Date.Equal(results[j].Date) {
			return results[i].Name < results[j].Name
		}
		return results[i].Date.Before(results[j].Date)
	})

	return results
}

func deduplicate(days []HariRaya) []HariRaya {
	seen := make(map[string]bool)
	var result []HariRaya
	for _, d := range days {
		key := d.Date.Format("2006-01-02") + "|" + d.Name
		if !seen[key] {
			seen[key] = true
			result = append(result, d)
		}
	}
	return result
}

// ForDate returns all holy days on a specific date.
func ForDate(t time.Time) []HariRaya {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	all := HolidaysBetween(t, t.AddDate(0, 0, 1))
	return all
}

// NextHoliday returns the next occurrence of a named holy day at or after t.
// name is case-insensitive substring match.
func NextHoliday(t time.Time, name string) *HariRaya {
	name = strings.ToLower(name)
	// Search forward year by year
	for year := t.Year(); year <= t.Year()+3; year++ {
		days := HolidaysInYear(year)
		for _, d := range days {
			if !d.Date.Before(t) && strings.Contains(strings.ToLower(d.Name), name) {
				return &d
			}
		}
	}
	return nil
}

// NextN returns the next n holy days at or after t.
func NextN(t time.Time, n int) []HariRaya {
	var result []HariRaya
	end := t.AddDate(2, 0, 0) // look up to 2 years ahead
	all := HolidaysBetween(t, end)
	for _, d := range all {
		if !d.Date.Before(t) {
			result = append(result, d)
			if len(result) >= n {
				break
			}
		}
	}
	return result
}

// Search returns all holy days matching name (case-insensitive) in the given year.
func Search(name string, year int) []HariRaya {
	name = strings.ToLower(name)
	all := HolidaysInYear(year)
	var result []HariRaya
	for _, d := range all {
		if strings.Contains(strings.ToLower(d.Name), name) {
			result = append(result, d)
		}
	}
	return result
}
