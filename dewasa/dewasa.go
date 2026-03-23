// Package dewasa implements Dewasa Ayu (auspicious days) and Dewasa Ala (inauspicious days)
// calculation for the Balinese Hindu calendar. All rules are pure deterministic
// functions of wewaran, wuku, and lunar positions — no external data needed.
//
// Source: https://sastrabali.com/rumus-menghitung-wariga/
package dewasa

import (
	"time"

	"github.com/pandeptwidyaop/kalenderbali-go/lunar"
	"github.com/pandeptwidyaop/kalenderbali-go/pawukon"
	"github.com/pandeptwidyaop/kalenderbali-go/wewaran"
)

// ── Supplementary Calculations ──────────────────────────────────────────────

// InternalDay holds enriched day data used by all rule matchers.
type InternalDay struct {
	Date    time.Time
	W       wewaran.WewaranResult
	Penanggal int // 1-15 (waxing) — lunar day within Sasih (1=first day after Tilem, 15=Purnama)
	Pangelong int // 1-15 (waning) — 0 means we're in penanggal phase
	IsPurnama bool
	IsTilem   bool
	SasihIdx  int
	SasihName string

	// Derived elements
	IngkelIdx  int // 0-5
	IngkelName string

	WatekMadyaIdx  int // 0-4
	WatekMadyaName string

	WatekAlitIdx  int // 0-3
	WatekAlitName string

	JejepanIdx  int // 0-5
	JejepanName string
}

var ingkelNames = [6]string{"Wong", "Sato", "Mina", "Manuk", "Taru", "Buku"}
var watekMadyaNames = [5]string{"Gajah", "Watu", "Bhuta", "Suku", "Wong"}
var watekAlitNames = [4]string{"Uler", "Gajah", "Lembu", "Lintah"}
var jejepanNames = [6]string{"Mina", "Taru", "Sato", "Patra", "Wong", "Paksi"}

// buildDay enriches a date with all derived values needed by rules.
func buildDay(t time.Time, phases []lunar.MoonPhase) InternalDay {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	w := wewaran.Calculate(t)
	pd := pawukon.DayOfCycle(t)

	// Ingkel: wukuNum % 6, 1-indexed (wukuNum is 1-indexed)
	wukuNum := w.WukuIndex + 1 // 1-30
	ingkelIdx := (wukuNum % 6) // 0 means remainder 0 → Buku (index 5)
	if ingkelIdx == 0 {
		ingkelIdx = 6
	}
	ingkelIdx-- // 0-indexed

	// Watek Madya: (uripPanca + uripSapta) % 5, 1-indexed
	uripSum := wewaran.PancawaraUrip[pd%5] + wewaran.SaptawaraUrip[pd%7]
	wmIdx := uripSum % 5
	if wmIdx == 0 {
		wmIdx = 5
	}
	wmIdx-- // 0-indexed

	// Watek Alit: (uripPanca + uripSapta) % 4, 1-indexed
	waIdx := uripSum % 4
	if waIdx == 0 {
		waIdx = 4
	}
	waIdx-- // 0-indexed

	// Jejepan: (wukuNum*7 + saptaNum) % 6, 1-indexed
	// saptaNum is bilangan saptawara: Redite=0, Soma=1, ..., Saniscara=6
	saptaNum := pd % 7
	jejRaw := (wukuNum*7 + saptaNum) % 6
	if jejRaw == 0 {
		jejRaw = 6
	}
	jejRaw-- // 0-indexed

	// Lunar phase context
	var isPurnama, isTilem bool
	var penanggal, pangelong int
	sasihIdx, sasihName := lunar.SasihForDate(t)

	// Find most recent Tilem at or before t, and most recent Purnama at or before t
	var recentTilem, recentPurnama time.Time
	for _, p := range phases {
		if !p.Date.After(t) {
			if p.Phase == lunar.NewMoon {
				if recentTilem.IsZero() || p.Date.After(recentTilem) {
					recentTilem = p.Date
				}
			} else {
				if recentPurnama.IsZero() || p.Date.After(recentPurnama) {
					recentPurnama = p.Date
				}
			}
		}
	}

	if !recentTilem.IsZero() && recentTilem.Equal(t) {
		isTilem = true
		penanggal = 0
		pangelong = 15
	} else if !recentPurnama.IsZero() && recentPurnama.Equal(t) {
		isPurnama = true
		penanggal = 15
		pangelong = 0
	} else if !recentTilem.IsZero() && !recentPurnama.IsZero() {
		daysSinceTilem := int(t.Sub(recentTilem).Hours() / 24)
		daysSincePurnama := int(t.Sub(recentPurnama).Hours() / 24)

		if recentTilem.After(recentPurnama) {
			// In penanggal phase (waxing, after Tilem)
			penanggal = daysSinceTilem
			if penanggal < 1 {
				penanggal = 1
			}
			if penanggal > 15 {
				penanggal = 15
			}
		} else {
			// In pangelong phase (waning, after Purnama)
			pangelong = daysSincePurnama
			if pangelong < 1 {
				pangelong = 1
			}
			if pangelong > 15 {
				pangelong = 15
			}
		}
	}

	_ = recentTilem
	_ = recentPurnama

	return InternalDay{
		Date:           t,
		W:              w,
		Penanggal:      penanggal,
		Pangelong:      pangelong,
		IsPurnama:      isPurnama,
		IsTilem:        isTilem,
		SasihIdx:       sasihIdx,
		SasihName:      sasihName,
		IngkelIdx:      ingkelIdx,
		IngkelName:     ingkelNames[ingkelIdx],
		WatekMadyaIdx:  wmIdx,
		WatekMadyaName: watekMadyaNames[wmIdx],
		WatekAlitIdx:   waIdx,
		WatekAlitName:  watekAlitNames[waIdx],
		JejepanIdx:     jejRaw,
		JejepanName:    jejepanNames[jejRaw],
	}
}

// ── Dewasa Result ────────────────────────────────────────────────────────────

// DewasaType indicates whether a dewasa is auspicious or inauspicious.
type DewasaType int

const (
	DewasaAyu DewasaType = iota // Hari Baik
	DewasaAla                   // Hari Buruk
	DewasaConjunction           // Conjunction / Special
)

func (d DewasaType) String() string {
	switch d {
	case DewasaAyu:
		return "Ayu"
	case DewasaAla:
		return "Ala"
	case DewasaConjunction:
		return "Konjungsi"
	}
	return "?"
}

// Dewasa represents one matching holy-day quality for a date.
type Dewasa struct {
	Name        string
	Type        DewasaType
	Description string
}

// DewasaResult is the full result for a single date.
type DewasaResult struct {
	Date        time.Time
	Ingkel      string
	WatekMadya  string
	WatekAlit   string
	Jejepan     string
	Penanggal   int
	Pangelong   int
	IsPurnama   bool
	IsTilem     bool
	SasihName   string
	DewasaList  []Dewasa
}

// ── Rule Engine ──────────────────────────────────────────────────────────────

// rule is a single dewasa rule.
type rule struct {
	dewasa  Dewasa
	matches func(d InternalDay) bool
}

// allRules is the master rule list, ordered: Ala first (higher priority), then Ayu, then Conjunction.
var allRules []rule

func init() {
	// ── HARI BURUK (Dewasa Ala) ──────────────────────────────────────────────

	// Ala Dahat: Saniscara nemu Purnama/Tilem
	registerRule("Ala Dahat", DewasaAla,
		"Hari teramat buruk — Saniscara bertemu Purnama atau Tilem. Hindari memulai kegiatan apapun.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 6 && (d.IsPurnama || d.IsTilem)
		})

	// Dagdig Krana: specific saptawara + penanggal/pangelong day
	// Redite 2, Soma 1, Anggara 10, Budha 7, Wraspati 3, Sukra 6
	dagdigRules := [][2]int{{0, 2}, {1, 1}, {2, 10}, {3, 7}, {4, 3}, {5, 6}}
	registerRule("Dagdig Krana", DewasaAla,
		"Hindari memulai (nuwasen/ngawit) karya atau pekerjaan penting jangka panjang.",
		func(d InternalDay) bool {
			for _, r := range dagdigRules {
				if d.W.SaptawaraIndex == r[0] {
					day := d.Penanggal
					if day == 0 {
						day = d.Pangelong
					}
					if day == r[1] {
						return true
					}
				}
			}
			return false
		})

	// Geneng Menyenget: specific saptawara + lunar day
	// Redite penanggal 4; Soma penanggal 1, pangelong 7; Anggara penanggal 2,10;
	// Budha pangelong 10; Wraspati pangelong 5; Sukra penanggal 14; Saniscara penanggal 1,9
	registerRule("Geneng Menyenget", DewasaAla,
		"Hari penuh godaan dan rintangan. Akan berakibat buruk bagi pemilik dan pelaku kegiatan.",
		func(d InternalDay) bool {
			s := d.W.SaptawaraIndex
			p := d.Penanggal
			pg := d.Pangelong
			switch s {
			case 0: // Redite
				return p == 4
			case 1: // Soma
				return p == 1 || pg == 7
			case 2: // Anggara
				return p == 2 || p == 10
			case 3: // Budha
				return pg == 10
			case 4: // Wraspati
				return pg == 5
			case 5: // Sukra
				return p == 14
			case 6: // Saniscara
				return p == 1 || p == 9
			}
			return false
		})

	// Sarik Agung: Budha + wuku Bala(24), Kulantir(3), Dunggulan(10), Merakih(17)
	sarikWuku := map[int]bool{3: true, 10: true, 17: true, 24: true}
	registerRule("Sarik Agung", DewasaAla,
		"Budha jatuh pada wuku Bala, Kulantir, Dunggulan, atau Merakih — hindari memulai usaha.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 3 && sarikWuku[d.W.WukuIndex]
		})

	// Pati Paten: Sukra Tilem, atau Sukra penanggal/pangelong 10
	registerRule("Pati Paten", DewasaAla,
		"Sukra bertemu Tilem, atau Sukra hari ke-10. Apapun kegiatan/usaha akan berakibat buruk.",
		func(d InternalDay) bool {
			if d.W.SaptawaraIndex != 5 { // bukan Sukra
				return false
			}
			if d.IsTilem {
				return true
			}
			day := d.Penanggal
			if day == 0 {
				day = d.Pangelong
			}
			return day == 10
		})

	// Tali Wangke: specific saptawara + wuku combos
	taliWangkeMap := map[int][]int{
		1: {21, 22, 23, 24, 25},    // Soma: Uye,Menail,Prangbakat,Bala,Ugu (0-indexed: 21-25 → wukuIdx 21..25 but names differ, let's map by wukuIdx 0-indexed)
		2: {26, 27, 28, 29, 0},     // Anggara: Wayang,Klawu,Dukut,Watugunung,Sinta
		3: {1, 2, 3, 4, 5},         // Budha: Landep,Ukir,Kulantir,Tolu→Taulu,Gumbreg
		4: {6, 7, 8, 9, 10, 16, 17, 18, 19, 13, 20}, // Wraspati: Wariga,Warigadean,Julungwangi,Sungsang,Dunggulan,Pahang→Krulut,Merakih,Tambir,Medangsia→(idx 13),Matal
		5: {11, 12, 13, 14, 15},    // Sukra: Kuningan,Langkir,Medangsia,Pujut,Pahang
	}
	registerRule("Tali Wangke", DewasaAla,
		"Hari buruk untuk memulai pekerjaan/usaha — saptawara bertemu wuku pantangan.",
		func(d InternalDay) bool {
			wuku := taliWangkeMap[d.W.SaptawaraIndex]
			for _, w := range wuku {
				if d.W.WukuIndex == w {
					return true
				}
			}
			return false
		})

	// Titi Buwuk: per-wuku specific saptawara+pancawara combo
	// (wukuIdx → [saptaIdx, pancaIdx])
	titiBuwukMap := map[int][2]int{
		0:  {2, 2}, // Sinta: Anggara Wage
		1:  {3, 0}, // Landep: Budha Paing (Paing=idx0)
		2:  {5, 4}, // Ukir: Sukra Umanis
		3:  {3, 4}, // Kulantir: Budha Umanis
		4:  {3, 1}, // Taulu: Budha Pon
		5:  {4, 4}, // Gumbreg: Wraspati Umanis
		6:  {1, 3}, // Wariga: Soma Keliwon
		7:  {1, 0}, // Warigadian: Soma Paing
		8:  {1, 2}, // Julungwangi: Soma Wage
		9:  {3, 1}, // Sungsang: Budha Pon
		10: {5, 0}, // Dunggulan: Sukra Paing
		11: {5, 2}, // Kuningan: Sukra Wage
		12: {4, 3}, // Langkir: Wraspati Keliwon
		13: {4, 0}, // Medangsia: Wraspati Paing
		14: {5, 3}, // Pujut: Sukra Keliwon
		15: {5, 0}, // Pahang: Sukra Paing
		16: {4, 1}, // Krulut: Wraspati Pon
		17: {0, 4}, // Merakih: Redite Umanis
		18: {3, 4}, // Tambir: Budha Umanis
		19: {1, 4}, // Medangkungan: Soma Umanis
		20: {2, 2}, // Matal: Anggara Wage
		21: {4, 1}, // Uye: Wraspati Pon
		22: {6, 0}, // Menail: Saniscara Paing
		23: {4, 0}, // Parangbakat: Wraspati Paing
		24: {3, 1}, // Bala: Budha Pon
		25: {0, 0}, // Ugu: Redite Paing
		26: {0, 2}, // Wayang: Redite Wage
		27: {0, 4}, // Kelawu: Redite Umanis
		28: {0, 1}, // Dukut: Redite Pon
		29: {0, 3}, // Watugunung: Redite Keliwon
	}
	registerRule("Titi Buwuk", DewasaAla,
		"Hari buruk untuk melakukan perjalanan jauh — saptawara+pancawara pantangan sesuai wuku.",
		func(d InternalDay) bool {
			combo, ok := titiBuwukMap[d.W.WukuIndex]
			if !ok {
				return false
			}
			return d.W.SaptawaraIndex == combo[0] && d.W.PancawaraIndex == combo[1]
		})

	// Tanpa Guru: wuku tanpa Guru (Asthawara) dalam satu wuku
	// Wuku: Kuningan(11), Medangkungan(19), Kelawu(27), Gumbreg(5)
	tanpaGuruWuku := map[int]bool{5: true, 11: true, 19: true, 27: true}
	registerRule("Tanpa Guru", DewasaAla,
		"Wuku ini tidak mengandung Guru (Asthawara) — tidak baik memulai usaha atau belajar.",
		func(d InternalDay) bool {
			return tanpaGuruWuku[d.W.WukuIndex]
		})

	// Prabu Pendah: Sukra penanggal 4
	registerRule("Prabu Pendah", DewasaAla,
		"Sukra penanggal ke-4 — hindari penobatan/pelantikan.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 5 && d.Penanggal == 4
		})

	// Prenggewa: Anggara penanggal 1
	registerRule("Prenggewa", DewasaAla,
		"Anggara penanggal pertama — hindari memulai usaha, rawan pertengkaran dan permusuhan.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 2 && d.Penanggal == 1
		})

	// Candung Watang: fires if date is not easily computable from wewaran alone
	// (needs piodalan context) — we implement as: day before kajeng keliwon after Purnama
	// Simplified: flag if Pangelong day 2 (2 days before next tilem cycle start) for common use
	// Note: true Candung Watang requires piodalan schedule, which varies per temple.
	// We skip this one as it requires external piodalan data — document it.

	// ── HARI BAIK (Dewasa Ayu) ───────────────────────────────────────────────

	// Dewasa Mentas: Wraspati pangelong 7 dan 15 (tilem)
	registerRule("Dewasa Mentas", DewasaAyu,
		"Wraspati pangelong ke-7 atau ke-15 (Tilem) — hari baik melakukan segala kegiatan usaha.",
		func(d InternalDay) bool {
			if d.W.SaptawaraIndex != 4 { // bukan Wraspati
				return false
			}
			return d.Pangelong == 7 || d.IsTilem
		})

	// Budha Wanas: Budha Wage nemu Purnama
	registerRule("Budha Wanas", DewasaAyu,
		"Budha Wage bertemu Purnama — hari baik semua kegiatan usaha, hasil pasti baik.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 3 && d.W.PancawaraIndex == 2 && d.IsPurnama
		})

	// Purwa Suka: Keliwon nemu Purnama
	registerRule("Purwa Suka", DewasaAyu,
		"Keliwon (Pancawara) bertemu Purnama — hari baik tapa semadi atau puasa.",
		func(d InternalDay) bool {
			return d.W.PancawaraIndex == 3 && d.IsPurnama
		})

	// Dewa Setata: Redite penanggal 10, Soma 9, Anggara 6, Wraspati 1
	dewaSetataRules := [][2]int{{0, 10}, {1, 9}, {2, 6}, {4, 1}}
	registerRule("Dewa Setata", DewasaAyu,
		"Hari baik melakukan segala kegiatan.",
		func(d InternalDay) bool {
			for _, r := range dewaSetataRules {
				if d.W.SaptawaraIndex == r[0] && d.Penanggal == r[1] {
					return true
				}
			}
			return false
		})

	// Subcara: specific saptawara + penanggal combos
	subcaraRules := map[int][]int{
		0: {3, 15},        // Redite: 3, 15(Purnama)
		1: {3},            // Soma: 3
		2: {2, 7, 8},      // Anggara: 2,7,8
		3: {2, 3},         // Budha: 2,3
		4: {5},            // Wraspati: 5
		5: {1, 2, 3},      // Sukra: 1,2,3
		6: {4, 5},         // Saniscara: 4,5
	}
	registerRule("Subcara", DewasaAyu,
		"Hari baik melakukan segala kegiatan.",
		func(d InternalDay) bool {
			days, ok := subcaraRules[d.W.SaptawaraIndex]
			if !ok {
				return false
			}
			p := d.Penanggal
			for _, dd := range days {
				if dd == 15 && d.IsPurnama {
					return true
				}
				if p == dd {
					return true
				}
			}
			return false
		})

	// Ayu Nulus: Redite 6, Soma 3, Anggara 7, Budha 12/13, Wraspati 5, Sukra 5
	ayuNulusRules := map[int][]int{
		0: {6},     // Redite
		1: {3},     // Soma
		2: {7},     // Anggara
		3: {12, 13},// Budha
		4: {5},     // Wraspati
		5: {5},     // Sukra
	}
	registerRule("Ayu Nulus", DewasaAyu,
		"Hari baik melakukan segala kegiatan.",
		func(d InternalDay) bool {
			days, ok := ayuNulusRules[d.W.SaptawaraIndex]
			if !ok {
				return false
			}
			for _, dd := range days {
				if d.Penanggal == dd {
					return true
				}
			}
			return false
		})

	// Mertha Buana: Redite/Soma/Anggara nemu Purnama
	registerRule("Mertha Buana", DewasaAyu,
		"Redite, Soma, atau Anggara bertemu Purnama — hari baik upacara yadnya dan widhiwedana.",
		func(d InternalDay) bool {
			s := d.W.SaptawaraIndex
			return (s == 0 || s == 1 || s == 2) && d.IsPurnama
		})

	// Sedana Yoga: specific saptawara + penanggal/pangelong days (including Purnama/Tilem)
	sedanaYogaRules := map[int][]int{
		0: {8},     // Redite: 8 (dan 15/Purnama/Tilem)
		1: {3},     // Soma: 3
		2: {7},     // Anggara: 7
		3: {2, 3},  // Budha: 2,3
		4: {4, 5},  // Wraspati: 4,5 (dan 15/Purnama/Tilem)
		5: {1, 6},  // Sukra: 1,6
		6: {5},     // Saniscara: 5 (dan 15/Purnama/Tilem)
	}
	sedanaYogaPurnamaTimlem := map[int]bool{0: true, 4: true, 6: true}
	registerRule("Sedana Yoga", DewasaAyu,
		"Hari baik memulai bisnis dan segala usaha yang mendatangkan rezeki.",
		func(d InternalDay) bool {
			days := sedanaYogaRules[d.W.SaptawaraIndex]
			p := d.Penanggal
			if p == 0 {
				p = d.Pangelong
			}
			for _, dd := range days {
				if p == dd {
					return true
				}
			}
			if sedanaYogaPurnamaTimlem[d.W.SaptawaraIndex] && (d.IsPurnama || d.IsTilem) {
				return true
			}
			return false
		})

	// Dina Mandi: Anggara Purnama, Wraspati penanggal 2, Sukra penanggal 14, Saniscara penanggal 3
	registerRule("Dina Mandi", DewasaAyu,
		"Hari kemanjuran — baik membuat jimat, melukat.",
		func(d InternalDay) bool {
			s := d.W.SaptawaraIndex
			if s == 2 && d.IsPurnama {
				return true
			}
			if s == 4 && d.Penanggal == 2 {
				return true
			}
			if s == 5 && d.Penanggal == 14 {
				return true
			}
			if s == 6 && d.Penanggal == 3 {
				return true
			}
			return false
		})

	// Siwa Sampurna: Wraspati penanggal 4, 5, 10
	registerRule("Siwa Sampurna", DewasaAyu,
		"Wraspati penanggal ke-4, 5, atau 10 — hari baik memulai/nangiang karya upacara yadnya.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 4 && (d.Penanggal == 4 || d.Penanggal == 5 || d.Penanggal == 10)
		})

	// Upadana Mertha: Redite penanggal 1, 3, 6, 10
	registerRule("Upadana Mertha", DewasaAyu,
		"Redite penanggal ke-1, 3, 6, atau 10 — hari baik memulai bisnis dan usaha.",
		func(d InternalDay) bool {
			p := d.Penanggal
			return d.W.SaptawaraIndex == 0 && (p == 1 || p == 3 || p == 6 || p == 10)
		})

	// Mertha Yoga: Soma wuku Ugu/Dukut/Landep/Krulut, Wraspati penanggal 4, Saniscara penanggal 5
	// Also: Sasih Kadasa(9) penanggal 5, Sasih Sadha(11) penanggal 1
	merthaYogaSomaWuku := map[int]bool{25: true, 28: true, 1: true, 16: true} // Ugu,Dukut,Landep,Krulut (0-indexed)
	registerRule("Mertha Yoga", DewasaAyu,
		"Hari baik membangun — Soma wuku Ugu/Dukut/Landep/Krulut; Wraspati penanggal 4; Saniscara penanggal 5.",
		func(d InternalDay) bool {
			if d.W.SaptawaraIndex == 1 && merthaYogaSomaWuku[d.W.WukuIndex] {
				return true
			}
			if d.W.SaptawaraIndex == 4 && d.Penanggal == 4 {
				return true
			}
			if d.W.SaptawaraIndex == 6 && d.Penanggal == 5 {
				return true
			}
			if d.SasihIdx == 9 && d.Penanggal == 5 {
				return true
			}
			if d.SasihIdx == 11 && d.Penanggal == 1 {
				return true
			}
			return false
		})

	// Sampi Gumarang Turun: specific saptawara + wuku
	sampiMap := map[int]int{
		0: 18, // Redite: Tambir(18)
		1: 5,  // Soma: Gumbreg(5)
		2: 22, // Anggara: Menail(22)
		3: 9,  // Budha: Sungsang(9)
		4: 26, // Wraspati: Wayang(26)
		5: 13, // Sukra: Medangsia(13)
	}
	registerRule("Sampi Gumarang Turun", DewasaAyu,
		"Hari baik membangun/nangiang rumah — saptawara bertemu wuku pasangannya.",
		func(d InternalDay) bool {
			wIdx, ok := sampiMap[d.W.SaptawaraIndex]
			return ok && d.W.WukuIndex == wIdx
		})

	// Dewasa Ngelayang: specific saptawara + penanggal (bercocok tanam, membuat perahu)
	ngelayangRules := map[int][]int{
		0: {1, 8},  // Redite
		2: {7},     // Anggara
		3: {2, 13}, // Budha
		4: {4},     // Wraspati
		5: {6},     // Sukra
		6: {5},     // Saniscara
	}
	registerRule("Dewasa Ngelayang", DewasaAyu,
		"Hari baik bercocok tanam dan membuat perahu/jukung.",
		func(d InternalDay) bool {
			days, ok := ngelayangRules[d.W.SaptawaraIndex]
			if !ok {
				return false
			}
			for _, dd := range days {
				if d.Penanggal == dd {
					return true
				}
			}
			return false
		})

	// Tutur Mandi: Wraspati + wuku Ukir/Julungwangi/Pujut/Medangkungan/Matal; Redite + Ugu
	tuturMandiWraspatiWuku := map[int]bool{2: true, 8: true, 14: true, 19: true, 20: true}
	registerRule("Tutur Mandi", DewasaAyu,
		"Hari baik memberikan petuah, mengajar, dan memulai semadi.",
		func(d InternalDay) bool {
			if d.W.SaptawaraIndex == 4 && tuturMandiWraspatiWuku[d.W.WukuIndex] {
				return true
			}
			if d.W.SaptawaraIndex == 0 && d.W.WukuIndex == 25 { // Redite+Ugu
				return true
			}
			return false
		})

	// Sari Ketah: Saniscara penanggal 4 dan 5
	registerRule("Sari Ketah", DewasaAyu,
		"Saniscara penanggal ke-4 atau ke-5 — hari baik membuat tembok batas rumah/panyengker.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 6 && (d.Penanggal == 4 || d.Penanggal == 5)
		})

	// ── CONJUNCTIONS ─────────────────────────────────────────────────────────

	// Kajeng Keliwon: Triwara Kajeng (idx 2) + Pancawara Keliwon (idx 3), every 15 days
	registerRule("Kajeng Keliwon", DewasaConjunction,
		"Kajeng Keliwon — hari penguatan spiritual, doa kepada Bhuta Kala. Setiap 15 hari.",
		func(d InternalDay) bool {
			return d.W.TriwaraIndex == 2 && d.W.PancawaraIndex == 3
		})

	// Buda Wage
	registerRule("Buda Wage", DewasaConjunction,
		"Buda Wage — hari keberuntungan dalam berdagang.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 3 && d.W.PancawaraIndex == 2
		})

	// Anggara Kasih: Anggara + Keliwon
	registerRule("Anggara Kasih", DewasaConjunction,
		"Anggara Kasih — hari berbakti kepada leluhur.",
		func(d InternalDay) bool {
			return d.W.SaptawaraIndex == 2 && d.W.PancawaraIndex == 3
		})
}

func registerRule(name string, t DewasaType, desc string, fn func(InternalDay) bool) {
	allRules = append(allRules, rule{
		dewasa:  Dewasa{Name: name, Type: t, Description: desc},
		matches: fn,
	})
}

// ── Public API ───────────────────────────────────────────────────────────────

// Calculate returns the full DewasaResult for a given date.
func Calculate(t time.Time) DewasaResult {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	phases := lunar.PhasesBetween(t.AddDate(0, -2, 0), t.AddDate(0, 2, 0))
	day := buildDay(t, phases)
	return applyRules(day)
}

func applyRules(day InternalDay) DewasaResult {
	var matched []Dewasa
	for _, r := range allRules {
		if r.matches(day) {
			matched = append(matched, r.dewasa)
		}
	}
	return DewasaResult{
		Date:       day.Date,
		Ingkel:     day.IngkelName,
		WatekMadya: day.WatekMadyaName,
		WatekAlit:  day.WatekAlitName,
		Jejepan:    day.JejepanName,
		Penanggal:  day.Penanggal,
		Pangelong:  day.Pangelong,
		IsPurnama:  day.IsPurnama,
		IsTilem:    day.IsTilem,
		SasihName:  day.SasihName,
		DewasaList: matched,
	}
}

// AyuInYear returns all Dewasa Ayu dates in a given year.
func AyuInYear(year int) []DewasaResult {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	phases := lunar.PhasesBetween(start.AddDate(0, -2, 0), end.AddDate(0, 2, 0))

	var results []DewasaResult
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		day := buildDay(d, phases)
		res := applyRules(day)
		hasAyu := false
		for _, dw := range res.DewasaList {
			if dw.Type == DewasaAyu {
				hasAyu = true
				break
			}
		}
		if hasAyu {
			results = append(results, res)
		}
	}
	return results
}

// AlaInYear returns all Dewasa Ala dates in a given year.
func AlaInYear(year int) []DewasaResult {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	phases := lunar.PhasesBetween(start.AddDate(0, -2, 0), end.AddDate(0, 2, 0))

	var results []DewasaResult
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		day := buildDay(d, phases)
		res := applyRules(day)
		hasAla := false
		for _, dw := range res.DewasaList {
			if dw.Type == DewasaAla {
				hasAla = true
				break
			}
		}
		if hasAla {
			results = append(results, res)
		}
	}
	return results
}

// ── Supplementary export: element names ─────────────────────────────────────

// IngkelNames returns all Ingkel names.
func IngkelNames() [6]string { return ingkelNames }

// WatekMadyaNames returns all Watek Madya names.
func WatekMadyaNames() [5]string { return watekMadyaNames }

// WatekAlitNames returns all Watek Alit names.
func WatekAlitNames() [4]string { return watekAlitNames }

// JejepanNames returns all Jejepan names.
func JejepanNames() [6]string { return jejepanNames }


