// Package wewaran implements all 10 concurrent Balinese week systems (wewaran).
package wewaran

import (
	"github.com/pandeptwidyaop/kalenderbali-go/pawukon"
	"time"
)

// ── Wuku ────────────────────────────────────────────────────────────────────

var WukuNames = [30]string{
	"Sinta", "Landep", "Ukir", "Kulantir", "Taulu", "Gumbreg",
	"Wariga", "Warigadian", "Julungwangi", "Sungsang",
	"Dunggulan", "Kuningan", "Langkir", "Medangsia", "Pujut",
	"Pahang", "Krulut", "Merakih", "Tambir", "Medangkungan",
	"Matal", "Uye", "Menail", "Parangbakat", "Bala",
	"Ugu", "Wayang", "Kelawu", "Dukut", "Watugunung",
}

// Wuku returns the 0-indexed wuku number and name.
func Wuku(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := d / 7
	return idx, WukuNames[idx]
}

// ── Saptawara (7-day) ───────────────────────────────────────────────────────

var SaptawaraNames = [7]string{
	"Redite", "Soma", "Anggara", "Buda", "Wraspati", "Sukra", "Saniscara",
}

var SaptawaraUrip = [7]int{5, 4, 3, 7, 8, 6, 9}

func Saptawara(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := d % 7
	return idx, SaptawaraNames[idx]
}

// ── Pancawara (5-day) ───────────────────────────────────────────────────────

var PancawaraNames = [5]string{"Paing", "Pon", "Wage", "Keliwon", "Umanis"}

var PancawaraUrip = [5]int{9, 7, 4, 8, 5}

func Pancawara(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := d % 5
	return idx, PancawaraNames[idx]
}

// ── Triwara (3-day) ─────────────────────────────────────────────────────────

var TriwaraNames = [3]string{"Pasah", "Beteng", "Kajeng"}

func Triwara(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := d % 3
	return idx, TriwaraNames[idx]
}

// ── Sadwara (6-day) ─────────────────────────────────────────────────────────

var SadwaraNames = [6]string{"Tungleh", "Aryang", "Urukung", "Paniron", "Was", "Maulu"}

func Sadwara(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := d % 6
	return idx, SadwaraNames[idx]
}

// ── Caturwara (4-day) ───────────────────────────────────────────────────────
// Special rule: days 71 and 72 are both "Jaya" (penultimate day repeats).

var CaturwaraNames = [4]string{"Sri", "Laba", "Jaya", "Menala"}

func CaturwaraIndex(d int) int {
	if d == 71 || d == 72 {
		return 2 // Jaya
	}
	// Days 71 and 72 are BOTH "Jaya" — two days are consumed for one position.
	// That means days 73+ are shifted back by 2 (not 1) relative to a plain modulo.
	if d > 72 {
		d -= 2
	}
	return d % 4
}

func Caturwara(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := CaturwaraIndex(d)
	return idx, CaturwaraNames[idx]
}

// ── Astawara (8-day) ────────────────────────────────────────────────────────
// Special rule: days 71 and 72 are both "Kala".

var AstawaraNames = [8]string{"Sri", "Indra", "Guru", "Yama", "Ludra", "Brahma", "Kala", "Uma"}

func AstawaraIndex(d int) int {
	if d == 71 || d == 72 {
		return 6 // Kala
	}
	// Same double-day rule as Caturwara: days 73+ shift back by 2.
	if d > 72 {
		d -= 2
	}
	return d % 8
}

func Astawara(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := AstawaraIndex(d)
	return idx, AstawaraNames[idx]
}

// ── Sangawara (9-day) ───────────────────────────────────────────────────────
// Special rule: first 3 days of Pawukon are all "Dangu".

var SangawaraNames = [9]string{
	"Dangu", "Jangur", "Gigis", "Nohan", "Ogan", "Erangan", "Urungan", "Tulus", "Dadi",
}

func SangawaraIndex(d int) int {
	if d < 3 {
		return 0 // Dangu
	}
	return (d - 3) % 9
}

func Sangawara(t time.Time) (int, string) {
	d := pawukon.DayOfCycle(t)
	idx := SangawaraIndex(d)
	return idx, SangawaraNames[idx]
}

// ── Urip calculation ─────────────────────────────────────────────────────────

// Urip returns the combined urip value used for Ekawara, Dwiwara, and Dasawara.
func Urip(t time.Time) int {
	d := pawukon.DayOfCycle(t)
	pancaIdx := d % 5
	saptaIdx := d % 7
	val := PancawaraUrip[pancaIdx] + SaptawaraUrip[saptaIdx] + 1
	if val > 10 {
		val -= 10
	}
	return val
}

// ── Ekawara (1-day) ─────────────────────────────────────────────────────────

func Ekawara(t time.Time) string {
	if Urip(t)%2 == 0 {
		return "Luang"
	}
	return ""
}

// ── Dwiwara (2-day) ─────────────────────────────────────────────────────────

var DwiwaraNames = [2]string{"Menga", "Pepet"}

func Dwiwara(t time.Time) (int, string) {
	u := Urip(t)
	if u%2 == 0 {
		return 1, "Pepet"
	}
	return 0, "Menga"
}

// ── Dasawara (10-day) ───────────────────────────────────────────────────────

var DasawaraNames = [10]string{
	"Pandita", "Pati", "Suka", "Duka", "Sri",
	"Manuh", "Manusa", "Raja", "Dewa", "Raksasa",
}

func Dasawara(t time.Time) (int, string) {
	u := Urip(t)
	idx := u - 1 // urip is 1-10, index is 0-9
	return idx, DasawaraNames[idx]
}

// ── Full Wewaran result ──────────────────────────────────────────────────────

// Wewaran holds all 10 week system values for a given day.
type WewaranResult struct {
	PawukonDay int
	WukuIndex  int
	WukuName   string

	// 10 week systems
	EkawaraName string // empty if not Luang

	DwiwaraIndex int
	DwiwaraName  string

	TriwaraIndex int
	TriwaraName  string

	CaturwaraIndex int
	CaturwaraName  string

	PancawaraIndex int
	PancawaraName  string

	SadwaraIndex int
	SadwaraName  string

	SaptawaraIndex int
	SaptawaraName  string

	AstawaraIndex int
	AstawaraName  string

	SangawaraIndex int
	SangawaraName  string

	DasawaraIndex int
	DasawaraName  string

	Urip int
}

// Calculate returns all wewaran info for a given date.
func Calculate(t time.Time) WewaranResult {
	d := pawukon.DayOfCycle(t)

	wukuIdx, wukuName := Wuku(t)
	triIdx, triName := Triwara(t)
	pancaIdx, pancaName := Pancawara(t)
	sadIdx, sadName := Sadwara(t)
	saptaIdx, saptaName := Saptawara(t)
	caturIdx, caturName := Caturwara(t)
	astaIdx, astaName := Astawara(t)
	sangaIdx, sangaName := Sangawara(t)
	dwiIdx, dwiName := Dwiwara(t)
	dasaIdx, dasaName := Dasawara(t)
	ekaName := Ekawara(t)
	u := Urip(t)

	return WewaranResult{
		PawukonDay:     d,
		WukuIndex:      wukuIdx,
		WukuName:       wukuName,
		EkawaraName:    ekaName,
		DwiwaraIndex:   dwiIdx,
		DwiwaraName:    dwiName,
		TriwaraIndex:   triIdx,
		TriwaraName:    triName,
		CaturwaraIndex: caturIdx,
		CaturwaraName:  caturName,
		PancawaraIndex: pancaIdx,
		PancawaraName:  pancaName,
		SadwaraIndex:   sadIdx,
		SadwaraName:    sadName,
		SaptawaraIndex: saptaIdx,
		SaptawaraName:  saptaName,
		AstawaraIndex:  astaIdx,
		AstawaraName:   astaName,
		SangawaraIndex: sangaIdx,
		SangawaraName:  sangaName,
		DasawaraIndex:  dasaIdx,
		DasawaraName:   dasaName,
		Urip:           u,
	}
}
