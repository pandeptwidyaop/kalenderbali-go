// Package lunar implements moon phase calculations using the Jean Meeus algorithm
// (from "Astronomical Algorithms", 2nd edition) for determining Purnama (full moon)
// and Tilem (new moon) dates in the Balinese Hindu calendar.
package lunar

import (
	"math"
	"time"
)

// MoonPhase represents a calculated moon phase event.
type MoonPhase struct {
	Date  time.Time
	Phase PhaseType
	JDE   float64 // Julian Day in dynamical time
}

// PhaseType indicates the type of moon phase.
type PhaseType int

const (
	NewMoon  PhaseType = 0 // Tilem
	FullMoon PhaseType = 2 // Purnama
)

func (p PhaseType) String() string {
	switch p {
	case NewMoon:
		return "Tilem"
	case FullMoon:
		return "Purnama"
	}
	return "Unknown"
}

// jdeToTime converts a Julian Ephemeris Day to a UTC time.Time.
// We use the approximation that Dynamical Time ≈ UTC for modern dates
// (ΔT is small, a few minutes, well within our 1-day accuracy requirement).
func jdeToTime(jde float64) time.Time {
	// JDE 2451545.0 = J2000.0 = 2000-01-01 12:00 TT
	// Difference in days from J2000.0
	jd := jde
	// Julian Day to calendar date (Meeus Ch. 7)
	jd += 0.5
	z := math.Floor(jd)
	f := jd - z

	var a float64
	if z < 2299161 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}
	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	day := b - d - math.Floor(30.6001*e) + f

	var month float64
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}
	var year float64
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}

	dayInt := int(math.Floor(day))
	fracDay := day - float64(dayInt)
	totalSec := int(math.Round(fracDay * 86400))
	hours := totalSec / 3600
	minutes := (totalSec % 3600) / 60
	seconds := totalSec % 60

	return time.Date(int(year), time.Month(int(month)), dayInt, hours, minutes, seconds, 0, time.UTC)
}

// rad converts degrees to radians.
func rad(deg float64) float64 { return deg * math.Pi / 180 }

// normalize360 brings an angle into [0, 360).
func normalize360(a float64) float64 {
	a = math.Mod(a, 360)
	if a < 0 {
		a += 360
	}
	return a
}

// meanPhaseJDE returns the approximate JDE for the k-th occurrence of the given phase.
// phase: 0 = new moon, 0.5 = first quarter, 1 = full moon, etc.
// (We use 0 for new moon and 0.5 for full moon per Meeus convention for k+phase.)
func meanPhaseJDE(k float64) float64 {
	// Meeus Ch. 49, Eq. 49.1
	T := k / 1236.85
	JDE := 2451550.09766 +
		29.530588861*k +
		0.00015437*T*T -
		0.000000150*T*T*T +
		0.00000000073*T*T*T*T
	return JDE
}

// phaseCorrection computes the full corrected JDE for a new or full moon.
// phase: 0 for new moon, 0.5 for full moon (added to k).
func phaseCorrection(k, phase float64) float64 {
	kp := k + phase
	T := kp / 1236.85

	JDE := meanPhaseJDE(kp)

	// Sun's mean anomaly (Meeus Eq. 49.4)
	M := normalize360(2.5534 + 29.10535670*kp - 0.0000014*T*T - 0.00000011*T*T*T)
	// Moon's mean anomaly
	Mprime := normalize360(201.5643 + 385.81693528*kp + 0.0107582*T*T + 0.00001238*T*T*T - 0.000000058*T*T*T*T)
	// Moon's argument of latitude
	F := normalize360(160.7108 + 390.67050284*kp - 0.0016118*T*T - 0.00000227*T*T*T + 0.000000011*T*T*T*T)
	// Longitude of ascending node
	Omega := normalize360(124.7746 - 1.56375588*kp + 0.0020672*T*T + 0.00000215*T*T*T)

	// Correction factor E (for Sun's anomaly)
	E := 1 - 0.002516*T - 0.0000074*T*T

	if phase == 0 {
		// New Moon corrections
		corr := -0.40720*math.Sin(rad(Mprime)) +
			0.17241*E*math.Sin(rad(M)) +
			0.01608*math.Sin(rad(2*Mprime)) +
			0.01039*math.Sin(rad(2*F)) +
			0.00739*E*math.Sin(rad(Mprime-M)) -
			0.00514*E*math.Sin(rad(Mprime+M)) +
			0.00208*E*E*math.Sin(rad(2*M)) -
			0.00111*math.Sin(rad(Mprime-2*F)) -
			0.00057*math.Sin(rad(Mprime+2*F)) +
			0.00056*E*math.Sin(rad(2*Mprime+M)) -
			0.00042*math.Sin(rad(3*Mprime)) +
			0.00042*E*math.Sin(rad(M+2*F)) +
			0.00038*E*math.Sin(rad(M-2*F)) -
			0.00024*E*math.Sin(rad(2*Mprime-M)) -
			0.00017*math.Sin(rad(Omega)) -
			0.00007*math.Sin(rad(Mprime+2*M)) +
			0.00004*math.Sin(rad(2*Mprime-2*F)) +
			0.00004*math.Sin(rad(3*M)) +
			0.00003*math.Sin(rad(Mprime+M-2*F)) +
			0.00003*math.Sin(rad(2*Mprime+2*F)) -
			0.00003*math.Sin(rad(Mprime+M+2*F)) +
			0.00003*math.Sin(rad(Mprime-M+2*F)) -
			0.00002*math.Sin(rad(Mprime-M-2*F)) -
			0.00002*math.Sin(rad(3*Mprime+M)) +
			0.00002*math.Sin(rad(4*Mprime))
		JDE += corr
	} else {
		// Full Moon corrections
		corr := -0.40614*math.Sin(rad(Mprime)) +
			0.17302*E*math.Sin(rad(M)) +
			0.01614*math.Sin(rad(2*Mprime)) +
			0.01043*math.Sin(rad(2*F)) +
			0.00734*E*math.Sin(rad(Mprime-M)) -
			0.00515*E*math.Sin(rad(Mprime+M)) +
			0.00209*E*E*math.Sin(rad(2*M)) -
			0.00111*math.Sin(rad(Mprime-2*F)) -
			0.00057*math.Sin(rad(Mprime+2*F)) +
			0.00056*E*math.Sin(rad(2*Mprime+M)) -
			0.00042*math.Sin(rad(3*Mprime)) +
			0.00042*E*math.Sin(rad(M+2*F)) +
			0.00038*E*math.Sin(rad(M-2*F)) -
			0.00024*E*math.Sin(rad(2*Mprime-M)) -
			0.00017*math.Sin(rad(Omega)) -
			0.00007*math.Sin(rad(Mprime+2*M)) +
			0.00004*math.Sin(rad(2*Mprime-2*F)) +
			0.00004*math.Sin(rad(3*M)) +
			0.00003*math.Sin(rad(Mprime+M-2*F)) +
			0.00003*math.Sin(rad(2*Mprime+2*F)) -
			0.00003*math.Sin(rad(Mprime+M+2*F)) +
			0.00003*math.Sin(rad(Mprime-M+2*F)) -
			0.00002*math.Sin(rad(Mprime-M-2*F)) -
			0.00002*math.Sin(rad(3*Mprime+M)) +
			0.00002*math.Sin(rad(4*Mprime))
		JDE += corr
	}

	// Additional planetary corrections (Meeus Table 49.d)
	add := 0.000325*math.Sin(rad(299.77+0.107408*kp-0.009173*T*T)) +
		0.000165*math.Sin(rad(251.88+0.016321*kp)) +
		0.000164*math.Sin(rad(251.83+26.651886*kp)) +
		0.000126*math.Sin(rad(349.42+36.412478*kp)) +
		0.000110*math.Sin(rad(84.66+18.206239*kp)) +
		0.000062*math.Sin(rad(141.74+53.303771*kp)) +
		0.000060*math.Sin(rad(207.14+2.453732*kp)) +
		0.000056*math.Sin(rad(154.84+7.306860*kp)) +
		0.000047*math.Sin(rad(34.52+27.261239*kp)) +
		0.000042*math.Sin(rad(207.19+0.121824*kp)) +
		0.000040*math.Sin(rad(291.34+1.844379*kp)) +
		0.000037*math.Sin(rad(161.72+24.198154*kp)) +
		0.000035*math.Sin(rad(239.56+25.513099*kp)) +
		0.000023*math.Sin(rad(331.55+3.592518*kp))
	JDE += add

	return JDE
}

// kForDate returns the approximate k value for a date and phase offset.
func kForDate(t time.Time, phase float64) float64 {
	year := float64(t.Year()) + (float64(t.YearDay())-1)/365.25
	return math.Floor((year-2000)*12.3685) - phase
}

// PhasesInYear returns all new moon and full moon events in the given year,
// sorted chronologically.
func PhasesInYear(year int) []MoonPhase {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)
	return PhasesBetween(start, end)
}

// phaseDate converts JDE to a Balinese date (UTC+8 approx).
// We subtract 0.5 days (12 hours) to better align with Balinese ritual day
// boundaries where events in the early AM are often counted as the previous day.
func phaseDate(jde float64) time.Time {
	t := jdeToTime(jde - 0.5)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// PhasesBetween returns all new moon and full moon events between start (inclusive)
// and end (exclusive).
func PhasesBetween(start, end time.Time) []MoonPhase {
	var results []MoonPhase

	// Start scanning from ~1 month before start to catch nearby phases
	k0 := kForDate(start.AddDate(0, -1, 0), 0)

	for _, phase := range []float64{0, 0.5} { // 0=new, 0.5=full
		k := k0
		for {
			jde := phaseCorrection(k, phase)
			d := phaseDate(jde)
			if d.Before(start) {
				k++
				continue
			}
			if !d.Before(end) {
				break
			}
			pt := NewMoon
			if phase == 0.5 {
				pt = FullMoon
			}
			results = append(results, MoonPhase{
				Date:  d,
				Phase: pt,
				JDE:   jde,
			})
			k++
		}
	}

	// Sort by date
	sortPhases(results)
	return results
}

func sortPhases(phases []MoonPhase) {
	n := len(phases)
	for i := 1; i < n; i++ {
		for j := i; j > 0 && phases[j].Date.Before(phases[j-1].Date); j-- {
			phases[j], phases[j-1] = phases[j-1], phases[j]
		}
	}
}

// NextPhase returns the next occurrence of the given phase at or after t.
func NextPhase(t time.Time, phase PhaseType) MoonPhase {
	var phaseOffset float64
	if phase == FullMoon {
		phaseOffset = 0.5
	}
	k := kForDate(t.AddDate(0, -1, 0), phaseOffset)
	for {
		jde := phaseCorrection(k, phaseOffset)
		d := phaseDate(jde)
		if !d.Before(t) {
			return MoonPhase{Date: d, Phase: phase, JDE: jde}
		}
		k++
	}
}

// ── Sasih (Balinese lunar month) ────────────────────────────────────────────

// SasihNames lists the 12 Balinese lunar months starting from Kasa.
var SasihNames = [12]string{
	"Kasa", "Karo", "Katiga", "Kapat", "Kalima", "Kanem",
	"Kapitu", "Kawolu", "Kasanga", "Kadasa", "Jyesta", "Sadha",
}

// SasihForDate returns the Sasih (Balinese lunar month) name for a given date.
// Sasih is based on the count of new moons since a reference point.
// Reference: Tilem Kapitu (the 7th new moon of the Balinese lunar year)
// corresponds roughly to January new moon near the start of the Gregorian year.
//
// We use the convention: the Balinese New Year (Nyepi) follows Tilem Kasanga.
// The Balinese lunar year begins with Kasa after the Tilem of the Kasanga month.
// We count new moons relative to the Tilem Kadasa reference.
// SasihForDate returns the Sasih (Balinese lunar month) index and name for a date.
//
// Balinese lunar month convention:
// - Sasih runs from the day AFTER Tilem to the next Tilem (inclusive).
// - The Tilem day itself is the LAST day of that sasih.
// - The day after Tilem begins the NEXT sasih (penanggal 1).
//
// So if date == Tilem Kasanga → sasih = Kasanga.
// If date == day after Tilem Kasanga → sasih = Kadasa.
func SasihForDate(t time.Time) (int, string) {
	day := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	// Find the most recent Tilem (new moon) at or before this date.
	prevTilem := NextPhase(day.AddDate(0, -2, 0), NewMoon)
	for prevTilem.Date.After(day) {
		prevTilem = NextPhase(prevTilem.Date.AddDate(0, -2, 0), NewMoon)
	}
	for {
		next := NextPhase(prevTilem.Date.AddDate(0, 0, 1), NewMoon)
		if next.Date.After(day) {
			break
		}
		prevTilem = next
	}

	// k value for this Tilem
	k := math.Round((float64(prevTilem.Date.Year()) +
		float64(prevTilem.Date.YearDay()-1)/365.25 - 2000) * 12.3685)

	// Base sasih for this Tilem: offset +8 aligns with known references.
	// At k=0 (2000-01-06) → Kapitu(6). Verified: Jan 2026 Tilem = Kapitu.
	tilemSasih := int(math.Mod(k+8, 12))
	if tilemSasih < 0 {
		tilemSasih += 12
	}

	// If the date IS the Tilem day → sasih = tilemSasih (last day of that month).
	// If the date is AFTER the Tilem → we're in penanggal of the NEXT sasih.
	sasihIndex := tilemSasih
	if day.After(prevTilem.Date) {
		sasihIndex = (tilemSasih + 1) % 12
	}

	return sasihIndex, SasihNames[sasihIndex]
}

// TilemSasih returns the sasih index for a given Tilem date (the Tilem itself).
// This is used by hariraya to check e.g. Tilem Kasanga for Nyepi.
func TilemSasih(tilemDate time.Time) (int, string) {
	k := math.Round((float64(tilemDate.Year()) +
		float64(tilemDate.YearDay()-1)/365.25 - 2000) * 12.3685)
	idx := int(math.Mod(k+8, 12))
	if idx < 0 {
		idx += 12
	}
	return idx, SasihNames[idx]
}
