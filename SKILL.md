# Kalender Bali Go — Skill Reference

## Overview
Pure Go library + CLI for complete Balinese Hindu calendar calculations. Zero external dependencies, zero API calls, zero database. Everything is deterministic math.

## Quick Start

```bash
cd ~/projects/kalenderbali-go
go build -o kalenderbali ./cmd/kalenderbali/
./kalenderbali today
```

## CLI Commands

### Info Hari
```bash
# Hari ini (lengkap: wewaran, wuku, sasih, laku, dewasa)
kalenderbali today

# Tanggal spesifik
kalenderbali date 2026-03-23

# Kalender bulanan (wewaran + hari raya)
kalenderbali bulan 2026-03
```

### Hari Raya
```bash
# Semua hari raya setahun
kalenderbali hariraya 2026

# 10 hari raya berikutnya
kalenderbali next

# Cari hari raya spesifik
kalenderbali next galungan
kalenderbali next purnama
kalenderbali next nyepi
kalenderbali search "pagerwesi" 2026
```

### Purnama & Tilem
```bash
kalenderbali purnama 2026    # Semua purnama
kalenderbali tilem 2026      # Semua tilem
```

### Dewasa Ayu/Ala (Hari Baik/Buruk)
```bash
kalenderbali dewasa 2026-03-23    # Dewasa untuk tanggal
kalenderbali dewasa-ayu 2026      # Semua hari baik setahun
kalenderbali dewasa-ala 2026      # Semua hari buruk setahun
```

### Pararasan (Laku / Ramalan Harian)
```bash
kalenderbali pararasan 2026-03-23  # Laku untuk tanggal
kalenderbali pararasan all         # Semua 35 kombinasi
```

### Ramalan Jodoh
```bash
# Masukkan tanggal lahir pria & wanita
kalenderbali jodoh 1997-10-08 2002-08-17
```
Output: 6 metode kecocokan dengan indikator 💚/💛/❤️‍🩹

### JSON Output
Semua command mendukung flag `--json`:
```bash
kalenderbali today --json
kalenderbali hariraya 2026 --json
kalenderbali jodoh 1997-10-08 2002-08-17 --json
```

## Packages (Library API)

### `pawukon` — Siklus 210 Hari
```go
import "github.com/pandeptwidyaop/kalenderbali-go/pawukon"

day := pawukon.DayOf(date)     // Posisi dalam siklus (0-209)
wuku := pawukon.WukuOf(date)   // Nama wuku (Sinta..Watugunung)
```

**Epoch:** 21 Mei 2000 = Redite Sinta (hari ke-0)
**Rumus:** `pawukonDay = ((daysSinceEpoch % 210) + 210) % 210`

### `wewaran` — 10 Sistem Minggu
```go
import "github.com/pandeptwidyaop/kalenderbali-go/wewaran"

w := wewaran.Of(date)
// w.Saptawara  → Redite/Soma/Anggara/Buda/Wraspati/Sukra/Saniscara
// w.Pancawara  → Paing/Pon/Wage/Keliwon/Umanis
// w.Triwara    → Pasah/Beteng/Kajeng
// w.Caturwara  → Sri/Laba/Jaya/Menala
// w.Sadwara    → Tungleh/Aryang/Urukung/Paniron/Was/Maulu
// w.Astawara   → Sri/Indra/Guru/Yama/Ludra/Brahma/Kala/Uma
// w.Sangawara  → Dangu/Jangur/Gigis/Nohan/Ogan/Erangan/Urungan/Tulus/Dadi
// w.Dwiwara    → Menga/Pepet
// w.Ekawara    → Luang atau kosong
// w.Dasawara   → Pandita/Pati/Suka/Duka/Sri/Manuh/Manusa/Raja/Dewa/Raksasa
// w.Ingkel     → Wong/Sato/Mina/Manuk/Taru/Buku
// w.WatekMadya → Gajah/Watu/Bhuta/Suku/Wong
// w.WatekAlit  → Uler/Gajah/Lembu/Lintah
// w.Jejepan    → Mina/Taru/Sato/Patra/Wong/Paksi
```

**Urip (nilai ritual):**
- Pancawara: Paing=9, Pon=7, Wage=4, Keliwon=8, Umanis=5
- Saptawara: Redite=5, Soma=4, Anggara=3, Buda=7, Wraspati=8, Sukra=6, Saniscara=9

### `lunar` — Fase Bulan (Jean Meeus)
```go
import "github.com/pandeptwidyaop/kalenderbali-go/lunar"

// Semua purnama & tilem dalam setahun
phases := lunar.PhasesInYear(2026)

// Purnama/Tilem berikutnya
next := lunar.NextPhase(time.Now(), lunar.FullMoon)

// Sasih (bulan Bali) untuk tanggal
idx, name := lunar.SasihForDate(date)
// Nama: Kasa, Karo, Katiga, Kapat, Kalima, Kanem,
//       Kapitu, Kawolu, Kasanga, Kadasa, Jyesta, Sadha
```

### `hariraya` — Hari Raya Hindu Bali
```go
import "github.com/pandeptwidyaop/kalenderbali-go/hariraya"

holidays := hariraya.InYear(2026)
// Termasuk:
// - Galungan, Kuningan, Penampahan, Manis Galungan/Kuningan
// - Saraswati, Banyu Pinaruh, Pagerwesi
// - Nyepi (dari Tilem Kasanga), Ngembak Geni
// - Siwaratri (Tilem Kapitu)
// - 5 Tumpek (Landep, Uduh, Kandang, Krulut, Wayang)
// - Purnama & Tilem per sasih
// - Kajeng Keliwon (tiap 15 hari)
// - Buda Wage, Anggara Kasih
```

**Pawukon-based (posisi tetap dalam siklus 210):**
| Hari Raya | Pawukon Day |
|-----------|-------------|
| Galungan | 74 (Buda Keliwon Dunggulan) |
| Kuningan | 84 (10 hari setelah Galungan) |
| Saraswati | 210 (Saniscara Umanis Watugunung) |
| Pagerwesi | 4 (Buda Keliwon Sinta) |

**Lunar-based:**
| Hari Raya | Trigger |
|-----------|---------|
| Nyepi | Sehari setelah Tilem Kasanga |
| Siwaratri | Tilem Kapitu |

### `dewasa` — Hari Baik/Buruk
```go
import "github.com/pandeptwidyaop/kalenderbali-go/dewasa"

results := dewasa.ForDate(date, wewaranInfo, lunarInfo)
// result.Name        → "Dewasa Mentas", "Ala Dahat", dll
// result.Type        → "ayu" (baik) atau "ala" (buruk)
// result.Description → penjelasan Bahasa Indonesia
```

**Dewasa Ayu (Baik):** Dewasa Mentas, Budha Wanas, Dewa Setata, Subcara, Ayu Nulus, Mertha Buana, Sedana Yoga, Siwa Sampurna, Upadana Mertha, Mertha Yoga, dll.

**Dewasa Ala (Buruk):** Ala Dahat, Dagdig Krana, Sarik Agung, Pati Paten, Tali Wangke, dll.

### `pararasan` — Laku (Ramalan Harian)
```go
import "github.com/pandeptwidyaop/kalenderbali-go/pararasan"

laku := pararasan.Of(saptawara, pancawara)
// laku.Name    → "Surya", "Bulan", "Api", "Toya", dll
// laku.Element → "Sun", "Moon", "Fire", "Water", dll
// laku.Desc    → deskripsi lengkap
```

35 kombinasi Saptawara × Pancawara, masing-masing punya Laku tetap.

### `jodoh` — Ramalan Kecocokan
```go
import "github.com/pandeptwidyaop/kalenderbali-go/jodoh"

results := jodoh.Check(birthPria, birthWanita)
// 6 metode: Saptawara, Neptu Mod 5, Neptu Mod 4,
//           Pertemuan Mod 9, Tri Pramana, Ramalan 5 Tahun
```

## Arsitektur Kalkulasi

```
Tanggal Masehi
    │
    ├──→ Pawukon (mod 210)
    │       ├──→ Wuku (30 nama)
    │       ├──→ Wewaran (10 sistem minggu)
    │       │       ├──→ Urip values
    │       │       ├──→ Ingkel, Watek, Jejepan
    │       │       └──→ Pararasan (Laku)
    │       ├──→ Hari Raya Pawukon (Galungan, dll)
    │       └──→ Dewasa (Pawukon-based rules)
    │
    └──→ Lunar (Jean Meeus)
            ├──→ Purnama & Tilem
            ├──→ Sasih (12 bulan Bali)
            ├──→ Penanggal/Pangelong
            ├──→ Hari Raya Lunar (Nyepi, Siwaratri)
            └──→ Dewasa (Lunar-based rules)
```

## Validasi

Cross-checked dengan kalenderbali.org untuk:
- Wewaran harian (Saptawara, Pancawara, Wuku, dll)
- Hari raya 2026 (Galungan, Kuningan, Nyepi, Tumpek, dll)
- Purnama & Tilem dates (±1 hari akurasi)
- Sasih alignment (Kapitu, Kasanga, dll)
- Pararasan/Laku harian

## File Structure
```
kalenderbali-go/
├── cmd/kalenderbali/main.go  # CLI tool
├── pawukon/pawukon.go        # Siklus 210 hari
├── wewaran/wewaran.go        # 10 sistem minggu + urip
├── lunar/lunar.go            # Jean Meeus moon phases + sasih
├── hariraya/hariraya.go      # Hari raya Hindu Bali
├── dewasa/dewasa.go          # Dewasa Ayu/Ala (rule engine)
├── pararasan/pararasan.go    # Laku/Pararasan (35 kombinasi)
├── jodoh/jodoh.go            # Ramalan kecocokan jodoh
├── *_test.go                 # Tests per package
├── go.mod                    # Go module (zero deps)
├── README.md                 # Documentation
└── SKILL.md                  # This file
```

## Referensi
- Lontar Wariga (pustaka leluhur Bali)
- kalenderbali.org (I Wayan Nuarsa, Universitas Udayana)
- sastrabali.com/rumus-menghitung-wariga
- Wikipedia: Pawukon Calendar
- Jean Meeus, *Astronomical Algorithms* (2nd ed.)
- Dershowitz & Reingold, *Calendrical Calculations*
