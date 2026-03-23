// kalenderbali — Balinese Calendar CLI
// Pure Go, zero external dependencies.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pandeptwidyaop/kalenderbali-go/dewasa"
	"github.com/pandeptwidyaop/kalenderbali-go/hariraya"
	"github.com/pandeptwidyaop/kalenderbali-go/lunar"
	"github.com/pandeptwidyaop/kalenderbali-go/pararasan"
	"github.com/pandeptwidyaop/kalenderbali-go/wewaran"
)

// ── Helpers ──────────────────────────────────────────────────────────────────

var jsonOut bool

func parseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("format tanggal tidak valid: %q (gunakan YYYY-MM-DD)", s)
	}
	return t, nil
}

func today() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.UTC)
}

func currentYear() int { return time.Now().Year() }

func bail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "❌ "+format+"\n", a...)
	os.Exit(1)
}

func printJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func hr() { fmt.Println(strings.Repeat("─", 54)) }

// ── Wewaran Printer ───────────────────────────────────────────────────────────

func printWewaran(t time.Time) {
	w := wewaran.Calculate(t)
	d := dewasa.Calculate(t)
	p := pararasan.Calculate(t)

	// Lunar info
	sasihIdx, sasihName := lunar.SasihForDate(t)
	_ = sasihIdx

	fmt.Printf("📅  %s\n", t.Format("Monday, 2 January 2006"))
	hr()
	fmt.Printf("🌀  Pawukon    : Hari ke-%d dari 210\n", w.PawukonDay+1)
	fmt.Printf("🎋  Wuku       : %s (wuku ke-%d)\n", w.WukuName, w.WukuIndex+1)
	hr()
	fmt.Println("── Wewaran (10 sistem minggu) ──")
	fmt.Printf("   Ekawara    : %s\n", ifEmpty(w.EkawaraName, "–"))
	fmt.Printf("   Dwiwara    : %s\n", w.DwiwaraName)
	fmt.Printf("   Triwara    : %s\n", w.TriwaraName)
	fmt.Printf("   Caturwara  : %s\n", w.CaturwaraName)
	fmt.Printf("   Pancawara  : %s (urip %d)\n", w.PancawaraName, wewaran.PancawaraUrip[w.PancawaraIndex])
	fmt.Printf("   Sadwara    : %s\n", w.SadwaraName)
	fmt.Printf("   Saptawara  : %s (urip %d)\n", w.SaptawaraName, wewaran.SaptawaraUrip[w.SaptawaraIndex])
	fmt.Printf("   Astawara   : %s\n", w.AstawaraName)
	fmt.Printf("   Sangawara  : %s\n", w.SangawaraName)
	fmt.Printf("   Dasawara   : %s\n", w.DasawaraName)
	fmt.Printf("   Urip       : %d\n", w.Urip)
	hr()
	fmt.Println("── Unsur Tambahan ──")
	fmt.Printf("   Ingkel     : %s\n", d.Ingkel)
	fmt.Printf("   Watek Madya: %s\n", d.WatekMadya)
	fmt.Printf("   Watek Alit : %s\n", d.WatekAlit)
	fmt.Printf("   Jejepan    : %s\n", d.Jejepan)
	hr()
	fmt.Println("── Sasih (Bulan Bali) ──")
	lunarPhaseStr := ""
	if d.IsPurnama {
		lunarPhaseStr = " 🌕 PURNAMA"
	} else if d.IsTilem {
		lunarPhaseStr = " 🌑 TILEM"
	} else if d.Penanggal > 0 {
		lunarPhaseStr = fmt.Sprintf(" (penanggal ke-%d)", d.Penanggal)
	} else if d.Pangelong > 0 {
		lunarPhaseStr = fmt.Sprintf(" (pangelong ke-%d)", d.Pangelong)
	}
	fmt.Printf("   Sasih      : %s%s\n", sasihName, lunarPhaseStr)
	hr()
	fmt.Printf("🌊  Laku (Pararasan): %s — %s\n", p.LakunName, p.LakunElement)
	fmt.Printf("   %s\n", p.LakunDesc)
	hr()

	// Hari Raya on this date
	hr2 := hariraya.ForDate(t)
	if len(hr2) > 0 {
		fmt.Println("🎊  Hari Raya / Upacara:")
		for _, h := range hr2 {
			fmt.Printf("   • %s — %s\n", h.Name, h.Description)
		}
		hr()
	}

	// Dewasa
	if len(d.DewasaList) > 0 {
		fmt.Println("✨  Dewasa (Kualitas Hari):")
		for _, dw := range d.DewasaList {
			icon := dewasaIcon(dw.Type)
			fmt.Printf("   %s %s\n", icon, dw.Name)
			fmt.Printf("      %s\n", dw.Description)
		}
	} else {
		fmt.Println("✨  Dewasa: Tidak ada dewasa khusus hari ini.")
	}
}

func ifEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func dewasaIcon(t dewasa.DewasaType) string {
	switch t {
	case dewasa.DewasaAyu:
		return "🟢"
	case dewasa.DewasaAla:
		return "🔴"
	case dewasa.DewasaConjunction:
		return "🔵"
	}
	return "⚪"
}

// ── Commands ─────────────────────────────────────────────────────────────────

func cmdToday() {
	t := today()
	if jsonOut {
		printDateJSON(t)
		return
	}
	fmt.Println("🌴  KALENDER BALI — HARI INI")
	printWewaran(t)
}

func cmdDate(args []string) {
	if len(args) == 0 {
		bail("Butuh argumen tanggal. Contoh: kalenderbali date 2026-03-23")
	}
	t, err := parseDate(args[0])
	if err != nil {
		bail("%v", err)
	}
	if jsonOut {
		printDateJSON(t)
		return
	}
	fmt.Println("🌴  KALENDER BALI")
	printWewaran(t)
}

func printDateJSON(t time.Time) {
	w := wewaran.Calculate(t)
	d := dewasa.Calculate(t)
	p := pararasan.Calculate(t)
	sasihIdx, sasihName := lunar.SasihForDate(t)
	_ = sasihIdx
	hr2 := hariraya.ForDate(t)

	type dewasaJSON struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	type hariRayaJSON struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
	}

	out := map[string]any{
		"date":      t.Format("2006-01-02"),
		"pawukon_day": w.PawukonDay + 1,
		"wuku":      map[string]any{"index": w.WukuIndex + 1, "name": w.WukuName},
		"wewaran": map[string]any{
			"ekawara":   w.EkawaraName,
			"dwiwara":   w.DwiwaraName,
			"triwara":   w.TriwaraName,
			"caturwara": w.CaturwaraName,
			"pancawara": w.PancawaraName,
			"sadwara":   w.SadwaraName,
			"saptawara": w.SaptawaraName,
			"astawara":  w.AstawaraName,
			"sangawara": w.SangawaraName,
			"dasawara":  w.DasawaraName,
			"urip":      w.Urip,
		},
		"extras": map[string]any{
			"ingkel":      d.Ingkel,
			"watek_madya": d.WatekMadya,
			"watek_alit":  d.WatekAlit,
			"jejepan":     d.Jejepan,
		},
		"sasih":     sasihName,
		"penanggal": d.Penanggal,
		"pangelong":  d.Pangelong,
		"is_purnama": d.IsPurnama,
		"is_tilem":  d.IsTilem,
		"laku": map[string]any{
			"name":        string(p.LakunName),
			"element":     p.LakunElement,
			"description": p.LakunDesc,
		},
	}

	var dewasaList []dewasaJSON
	for _, dw := range d.DewasaList {
		dewasaList = append(dewasaList, dewasaJSON{
			Name:        dw.Name,
			Type:        dw.Type.String(),
			Description: dw.Description,
		})
	}
	out["dewasa"] = dewasaList

	var hrList []hariRayaJSON
	for _, h := range hr2 {
		hrList = append(hrList, hariRayaJSON{
			Name:        h.Name,
			Description: h.Description,
			Category:    h.Category,
		})
	}
	out["hari_raya"] = hrList

	printJSON(out)
}

func cmdWuku(args []string) {
	t := today()
	if len(args) > 0 {
		var err error
		t, err = parseDate(args[0])
		if err != nil {
			bail("%v", err)
		}
	}
	w := wewaran.Calculate(t)
	if jsonOut {
		printJSON(map[string]any{
			"date":       t.Format("2006-01-02"),
			"wuku_index": w.WukuIndex + 1,
			"wuku_name":  w.WukuName,
			"pawukon_day": w.PawukonDay + 1,
		})
		return
	}
	fmt.Printf("📅 %s\n", t.Format("2 January 2006"))
	fmt.Printf("🎋 Wuku: %s (wuku ke-%d, hari pawukon ke-%d)\n", w.WukuName, w.WukuIndex+1, w.PawukonDay+1)
}

func cmdPurnama(args []string) {
	year := currentYear()
	if len(args) > 0 {
		y, err := strconv.Atoi(args[0])
		if err != nil || y < 1 {
			bail("Tahun tidak valid: %q", args[0])
		}
		year = y
	}
	phases := lunar.PhasesInYear(year)
	var list []lunar.MoonPhase
	for _, p := range phases {
		if p.Phase == lunar.FullMoon {
			list = append(list, p)
		}
	}
	if jsonOut {
		type item struct {
			Date  string `json:"date"`
			Phase string `json:"phase"`
			Sasih string `json:"sasih"`
		}
		var out []item
		for _, p := range list {
			_, s := lunar.SasihForDate(p.Date)
			out = append(out, item{p.Date.Format("2006-01-02"), "Purnama", s})
		}
		printJSON(out)
		return
	}
	fmt.Printf("🌕  PURNAMA %d\n", year)
	hr()
	for i, p := range list {
		_, s := lunar.SasihForDate(p.Date)
		fmt.Printf("  %2d. %s  — Sasih %s\n", i+1, p.Date.Format("Mon, 2 Jan 2006"), s)
	}
}

func cmdTilem(args []string) {
	year := currentYear()
	if len(args) > 0 {
		y, err := strconv.Atoi(args[0])
		if err != nil || y < 1 {
			bail("Tahun tidak valid: %q", args[0])
		}
		year = y
	}
	phases := lunar.PhasesInYear(year)
	var list []lunar.MoonPhase
	for _, p := range phases {
		if p.Phase == lunar.NewMoon {
			list = append(list, p)
		}
	}
	if jsonOut {
		type item struct {
			Date  string `json:"date"`
			Phase string `json:"phase"`
			Sasih string `json:"sasih"`
		}
		var out []item
		for _, p := range list {
			_, s := lunar.SasihForDate(p.Date)
			out = append(out, item{p.Date.Format("2006-01-02"), "Tilem", s})
		}
		printJSON(out)
		return
	}
	fmt.Printf("🌑  TILEM %d\n", year)
	hr()
	for i, p := range list {
		_, s := lunar.SasihForDate(p.Date)
		fmt.Printf("  %2d. %s  — Sasih %s\n", i+1, p.Date.Format("Mon, 2 Jan 2006"), s)
	}
}

func cmdHariRaya(args []string) {
	year := currentYear()
	if len(args) > 0 {
		y, err := strconv.Atoi(args[0])
		if err != nil || y < 1 {
			bail("Tahun tidak valid: %q", args[0])
		}
		year = y
	}
	list := hariraya.HolidaysInYear(year)
	if jsonOut {
		printJSON(list)
		return
	}
	fmt.Printf("🎊  HARI RAYA BALI %d\n", year)
	hr()
	for _, h := range list {
		w := wewaran.Calculate(h.Date)
		fmt.Printf("  📅 %-28s %s, %s %s\n",
			h.Name,
			h.Date.Format("2 Jan 2006"),
			w.SaptawaraName,
			w.WukuName,
		)
	}
}

func cmdNext(args []string) {
	t := today()
	if len(args) == 0 {
		// Next 10 holy days
		list := hariraya.NextN(t, 10)
		if jsonOut {
			printJSON(list)
			return
		}
		fmt.Println("⏭️   10 HARI RAYA BERIKUTNYA")
		hr()
		for _, h := range list {
			fmt.Printf("  📅 %-28s %s\n", h.Name, h.Date.Format("Mon, 2 Jan 2006"))
		}
		return
	}

	keyword := strings.Join(args, " ")
	lk := strings.ToLower(keyword)

	// Special cases: purnama, tilem
	if lk == "purnama" {
		p := lunar.NextPhase(t, lunar.FullMoon)
		_, s := lunar.SasihForDate(p.Date)
		if jsonOut {
			printJSON(map[string]any{"date": p.Date.Format("2006-01-02"), "phase": "Purnama", "sasih": s})
			return
		}
		fmt.Printf("🌕  Purnama berikutnya: %s — Sasih %s\n", p.Date.Format("Mon, 2 Jan 2006"), s)
		return
	}
	if lk == "tilem" {
		p := lunar.NextPhase(t, lunar.NewMoon)
		_, s := lunar.SasihForDate(p.Date)
		if jsonOut {
			printJSON(map[string]any{"date": p.Date.Format("2006-01-02"), "phase": "Tilem", "sasih": s})
			return
		}
		fmt.Printf("🌑  Tilem berikutnya: %s — Sasih %s\n", p.Date.Format("Mon, 2 Jan 2006"), s)
		return
	}

	h := hariraya.NextHoliday(t, keyword)
	if h == nil {
		fmt.Printf("Tidak ditemukan hari raya dengan kata kunci %q\n", keyword)
		os.Exit(1)
	}
	if jsonOut {
		printJSON(h)
		return
	}
	fmt.Printf("⏭️   %s berikutnya: %s\n", h.Name, h.Date.Format("Mon, 2 Jan 2006"))
}

func cmdSearch(args []string) {
	if len(args) == 0 {
		bail("Butuh keyword. Contoh: kalenderbali search pagerwesi 2026")
	}
	year := currentYear()
	keyword := args[0]
	if len(args) > 1 {
		y, err := strconv.Atoi(args[1])
		if err == nil && y > 0 {
			year = y
		}
	}
	list := hariraya.Search(keyword, year)
	if jsonOut {
		printJSON(list)
		return
	}
	if len(list) == 0 {
		fmt.Printf("Tidak ditemukan %q di tahun %d\n", keyword, year)
		return
	}
	fmt.Printf("🔍  Pencarian %q tahun %d:\n", keyword, year)
	hr()
	for _, h := range list {
		fmt.Printf("  📅 %-28s %s  [%s]\n", h.Name, h.Date.Format("2 Jan 2006"), h.Category)
	}
}

// ── Dewasa Commands ───────────────────────────────────────────────────────────

func cmdDewasa(args []string) {
	t := today()
	if len(args) > 0 {
		var err error
		t, err = parseDate(args[0])
		if err != nil {
			bail("%v", err)
		}
	}
	d := dewasa.Calculate(t)
	w := wewaran.Calculate(t)
	p := pararasan.Calculate(t)

	if jsonOut {
		type dewasaJSON struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Description string `json:"description"`
		}
		var dlist []dewasaJSON
		for _, dw := range d.DewasaList {
			dlist = append(dlist, dewasaJSON{dw.Name, dw.Type.String(), dw.Description})
		}
		printJSON(map[string]any{
			"date":       t.Format("2006-01-02"),
			"saptawara":  w.SaptawaraName,
			"pancawara":  w.PancawaraName,
			"wuku":       w.WukuName,
			"ingkel":     d.Ingkel,
			"watek_madya": d.WatekMadya,
			"watek_alit": d.WatekAlit,
			"jejepan":    d.Jejepan,
			"penanggal":  d.Penanggal,
			"pangelong":  d.Pangelong,
			"is_purnama": d.IsPurnama,
			"is_tilem":   d.IsTilem,
			"sasih":      d.SasihName,
			"laku":       string(p.LakunName),
			"dewasa":     dlist,
		})
		return
	}

	fmt.Printf("🔮  DEWASA — %s\n", t.Format("2 January 2006"))
	hr()
	fmt.Printf("   %s %s, Wuku %s\n", w.SaptawaraName, w.PancawaraName, w.WukuName)
	fmt.Printf("   Ingkel: %-12s  Watek Madya: %-8s  Watek Alit: %-8s\n",
		d.Ingkel, d.WatekMadya, d.WatekAlit)
	fmt.Printf("   Jejepan: %s\n", d.Jejepan)
	lunarStr := ""
	if d.IsPurnama {
		lunarStr = "🌕 PURNAMA"
	} else if d.IsTilem {
		lunarStr = "🌑 TILEM"
	} else if d.Penanggal > 0 {
		lunarStr = fmt.Sprintf("Penanggal %d", d.Penanggal)
	} else if d.Pangelong > 0 {
		lunarStr = fmt.Sprintf("Pangelong %d", d.Pangelong)
	}
	fmt.Printf("   Sasih: %s  %s\n", d.SasihName, lunarStr)
	fmt.Printf("   Laku: %s — %s\n", p.LakunName, p.LakunElement)
	hr()

	if len(d.DewasaList) == 0 {
		fmt.Println("   Tidak ada dewasa khusus hari ini.")
		return
	}

	// Group by type
	var ayu, ala, conj []dewasa.Dewasa
	for _, dw := range d.DewasaList {
		switch dw.Type {
		case dewasa.DewasaAyu:
			ayu = append(ayu, dw)
		case dewasa.DewasaAla:
			ala = append(ala, dw)
		case dewasa.DewasaConjunction:
			conj = append(conj, dw)
		}
	}

	if len(ala) > 0 {
		fmt.Println("🔴  DEWASA ALA (Hari Buruk):")
		for _, dw := range ala {
			fmt.Printf("   ⚠️  %s\n", dw.Name)
			fmt.Printf("      %s\n", dw.Description)
		}
	}
	if len(ayu) > 0 {
		fmt.Println("🟢  DEWASA AYU (Hari Baik):")
		for _, dw := range ayu {
			fmt.Printf("   ✅ %s\n", dw.Name)
			fmt.Printf("      %s\n", dw.Description)
		}
	}
	if len(conj) > 0 {
		fmt.Println("🔵  KONJUNGSI KHUSUS:")
		for _, dw := range conj {
			fmt.Printf("   🌀 %s\n", dw.Name)
			fmt.Printf("      %s\n", dw.Description)
		}
	}
}

func cmdDewasaAyu(args []string) {
	year := currentYear()
	if len(args) > 0 {
		y, err := strconv.Atoi(args[0])
		if err == nil && y > 0 {
			year = y
		}
	}
	list := dewasa.AyuInYear(year)
	if jsonOut {
		printJSON(list)
		return
	}
	fmt.Printf("🟢  DEWASA AYU %d — Semua Hari Baik\n", year)
	hr()
	for _, d := range list {
		w := wewaran.Calculate(d.Date)
		var names []string
		for _, dw := range d.DewasaList {
			if dw.Type == dewasa.DewasaAyu {
				names = append(names, dw.Name)
			}
		}
		fmt.Printf("  📅 %-12s  %-12s %-10s  %s\n",
			d.Date.Format("2 Jan 2006"),
			w.SaptawaraName+" "+w.PancawaraName,
			w.WukuName,
			strings.Join(names, ", "),
		)
	}
	fmt.Printf("\nTotal: %d hari baik di tahun %d\n", len(list), year)
}

func cmdDewasaAla(args []string) {
	year := currentYear()
	if len(args) > 0 {
		y, err := strconv.Atoi(args[0])
		if err == nil && y > 0 {
			year = y
		}
	}
	list := dewasa.AlaInYear(year)
	if jsonOut {
		printJSON(list)
		return
	}
	fmt.Printf("🔴  DEWASA ALA %d — Semua Hari Buruk\n", year)
	hr()
	for _, d := range list {
		w := wewaran.Calculate(d.Date)
		var names []string
		for _, dw := range d.DewasaList {
			if dw.Type == dewasa.DewasaAla {
				names = append(names, dw.Name)
			}
		}
		fmt.Printf("  📅 %-12s  %-12s %-10s  %s\n",
			d.Date.Format("2 Jan 2006"),
			w.SaptawaraName+" "+w.PancawaraName,
			w.WukuName,
			strings.Join(names, ", "),
		)
	}
	fmt.Printf("\nTotal: %d hari buruk di tahun %d\n", len(list), year)
}

// ── Pararasan Command ─────────────────────────────────────────────────────────

func cmdPararasan(args []string) {
	if len(args) > 0 && args[0] == "all" {
		all := pararasan.All35()
		if jsonOut {
			printJSON(all)
			return
		}
		fmt.Println("🌊  PARARASAN — 35 Kombinasi Saptawara × Pancawara")
		hr()
		fmt.Printf("  %-12s %-10s %-16s %s\n", "Saptawara", "Pancawara", "Laku", "Elemen")
		hr()
		for _, r := range all {
			fmt.Printf("  %-12s %-10s %-16s %s\n",
				r.SaptawaraName, r.PancawaraName, r.LakunName, r.LakunElement)
		}
		return
	}

	t := today()
	if len(args) > 0 {
		var err error
		t, err = parseDate(args[0])
		if err != nil {
			bail("%v", err)
		}
	}
	r := pararasan.Calculate(t)
	if jsonOut {
		printJSON(r)
		return
	}
	fmt.Printf("🌊  PARARASAN (LAKU) — %s\n", t.Format("2 January 2006"))
	hr()
	fmt.Printf("   Saptawara  : %s\n", r.SaptawaraName)
	fmt.Printf("   Pancawara  : %s\n", r.PancawaraName)
	fmt.Printf("   Laku       : %s\n", r.LakunName)
	fmt.Printf("   Elemen     : %s\n", r.LakunElement)
	fmt.Printf("   Makna      : %s\n", r.LakunDesc)
}

// ── Lahir Command ─────────────────────────────────────────────────────────────

func cmdLahir(args []string) {
	if len(args) == 0 {
		bail("Butuh tanggal lahir. Contoh: kalenderbali lahir 1997-10-08")
	}
	t, err := parseDate(args[0])
	if err != nil {
		bail("%v", err)
	}

	w := wewaran.Calculate(t)
	p := pararasan.Calculate(t)
	d := dewasa.Calculate(t)

	if jsonOut {
		printJSON(map[string]any{
			"tanggal_lahir": t.Format("2006-01-02"),
			"saptawara":     w.SaptawaraName,
			"pancawara":     w.PancawaraName,
			"wuku":          w.WukuName,
			"ingkel":        d.Ingkel,
			"watek_madya":   d.WatekMadya,
			"watek_alit":    d.WatekAlit,
			"jejepan":       d.Jejepan,
			"dasawara":      w.DasawaraName,
			"urip":          w.Urip,
			"laku":          string(p.LakunName),
			"laku_elemen":   p.LakunElement,
			"laku_desc":     p.LakunDesc,
		})
		return
	}

	fmt.Printf("👶  RAMALAN KELAHIRAN — %s\n", t.Format("2 January 2006"))
	hr()
	fmt.Printf("   Saptawara  : %s\n", w.SaptawaraName)
	fmt.Printf("   Pancawara  : %s\n", w.PancawaraName)
	fmt.Printf("   Wuku       : %s\n", w.WukuName)
	fmt.Printf("   Dasawara   : %s\n", w.DasawaraName)
	fmt.Printf("   Urip       : %d\n", w.Urip)
	hr()
	fmt.Println("── Watak & Karakter ──")
	fmt.Printf("   Ingkel     : %s\n", d.Ingkel)
	fmt.Printf("   Watek Madya: %s\n", d.WatekMadya)
	fmt.Printf("   Watek Alit : %s\n", d.WatekAlit)
	fmt.Printf("   Jejepan    : %s\n", d.Jejepan)
	hr()
	fmt.Printf("🌊  Laku (Pararasan): %s — %s\n", p.LakunName, p.LakunElement)
	fmt.Printf("   %s\n", p.LakunDesc)
}

// ── Help & Main ───────────────────────────────────────────────────────────────

func usage() {
	fmt.Println(`🌴  kalenderbali — Kalender Bali (Pure Go)

PERINTAH:
  today                      Info kalender Bali hari ini (lengkap)
  date <YYYY-MM-DD>          Info kalender untuk tanggal tertentu
  wuku [YYYY-MM-DD]          Wuku untuk tanggal (default: hari ini)
  purnama [tahun]            Daftar semua Purnama di tahun ini/tahun tertentu
  tilem [tahun]              Daftar semua Tilem di tahun ini/tahun tertentu
  hariraya [tahun]           Daftar semua hari raya Bali
  next [keyword]             10 hari raya berikutnya, atau cari "galungan", "purnama", dll
  search <keyword> [tahun]   Cari hari raya berdasarkan nama
  dewasa [YYYY-MM-DD]        Dewasa Ayu/Ala untuk tanggal tertentu (default: hari ini)
  dewasa-ayu [tahun]         Semua hari baik di tahun ini/tahun tertentu
  dewasa-ala [tahun]         Semua hari buruk di tahun ini/tahun tertentu
  pararasan [YYYY-MM-DD|all] Laku (pararasan) untuk tanggal, atau tampilkan semua 35 kombinasi
  lahir <YYYY-MM-DD>         Ramalan kelahiran: laku, watak, wuku

FLAG:
  --json                     Output dalam format JSON

CONTOH:
  kalenderbali today
  kalenderbali date 2026-03-23
  kalenderbali next galungan
  kalenderbali dewasa 2026-03-23
  kalenderbali dewasa-ayu 2026
  kalenderbali pararasan all
  kalenderbali lahir 1997-10-08
  kalenderbali hariraya 2026 --json`)
}

func main() {
	args := os.Args[1:]

	// Strip --json flag
	var filtered []string
	for _, a := range args {
		if a == "--json" {
			jsonOut = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered

	if len(args) == 0 {
		usage()
		os.Exit(0)
	}

	cmd := strings.ToLower(args[0])
	rest := args[1:]

	switch cmd {
	case "today":
		cmdToday()
	case "date":
		cmdDate(rest)
	case "wuku":
		cmdWuku(rest)
	case "purnama":
		cmdPurnama(rest)
	case "tilem":
		cmdTilem(rest)
	case "hariraya", "hari-raya":
		cmdHariRaya(rest)
	case "next":
		cmdNext(rest)
	case "search":
		cmdSearch(rest)
	case "dewasa":
		cmdDewasa(rest)
	case "dewasa-ayu":
		cmdDewasaAyu(rest)
	case "dewasa-ala":
		cmdDewasaAla(rest)
	case "pararasan":
		cmdPararasan(rest)
	case "lahir":
		cmdLahir(rest)
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "Perintah tidak dikenal: %q\n", cmd)
		usage()
		os.Exit(1)
	}
}
