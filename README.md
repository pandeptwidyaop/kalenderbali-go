# 🌴 Kalender Bali Go

Pustaka dan CLI Kalender Bali yang lengkap, ditulis dalam **pure Go** — tanpa database, tanpa API eksternal, tanpa dependency. Semua perhitungan menggunakan rumus matematika murni.

[![Go](https://img.shields.io/badge/Go-1.21+-blue)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-green)](LICENSE)

## ✨ Fitur

| Package | Deskripsi |
|---------|-----------|
| `pawukon` | Siklus 210 hari Pawukon — hitung posisi hari dari tanggal Masehi |
| `wewaran` | 10 sistem minggu bersamaan (Ekawara s/d Dasawara) + nilai urip |
| `lunar` | Fase bulan (Jean Meeus algorithm) — Purnama & Tilem + Sasih |
| `hariraya` | Hari raya Hindu Bali — Galungan, Kuningan, Nyepi, Saraswati, dll |
| `dewasa` | Dewasa Ayu (hari baik) & Dewasa Ala (hari buruk) — rule engine |
| `pararasan` | Laku/Pararasan — 35 ramalan harian dari Saptawara × Pancawara |
| `jodoh` | Ramalan kecocokan jodoh — 6 metode tradisional Bali |

### Unsur Wariga Tambahan
- **Ingkel** — pantangan hari berdasarkan wuku
- **Watek Madya & Alit** — karakter hari berdasarkan urip
- **Jejepan** — klasifikasi makhluk hidup pada hari tersebut
- **Sasih** — bulan lunar Bali (Kasa s/d Sadha)
- **Penanggal/Pangelong** — fase bulan waxing/waning (1–15)

## 📦 Instalasi

```bash
go install github.com/pandeptwidyaop/kalenderbali-go/cmd/kalenderbali@latest
```

Atau build dari source:

```bash
git clone https://github.com/pandeptwidyaop/kalenderbali-go.git
cd kalenderbali-go
go build -o kalenderbali ./cmd/kalenderbali/
```

## 🖥️ Penggunaan CLI

```bash
# Info lengkap hari ini
kalenderbali today

# Info tanggal tertentu (lengkap: wewaran, sasih, dewasa, laku, hari raya)
kalenderbali date 2026-03-23

# Hanya wuku untuk tanggal tertentu
kalenderbali wuku 2026-03-23

# Jadwal Purnama & Tilem setahun penuh
kalenderbali purnama 2026
kalenderbali tilem 2026

# Hari raya setahun penuh
kalenderbali hariraya 2026

# 10 hari raya berikutnya
kalenderbali next

# Cari hari raya spesifik berikutnya
kalenderbali next galungan
kalenderbali next purnama
kalenderbali next tilem

# Cari hari raya berdasarkan nama
kalenderbali search tumpek 2026

# Dewasa Ayu/Ala untuk tanggal tertentu
kalenderbali dewasa 2026-03-23

# Semua hari baik/buruk dalam setahun
kalenderbali dewasa-ayu 2026
kalenderbali dewasa-ala 2026

# Pararasan (Laku) untuk tanggal
kalenderbali pararasan 2026-03-23

# Tampilkan semua 35 kombinasi Saptawara × Pancawara
kalenderbali pararasan all

# Ramalan kelahiran (watak, laku, wuku)
kalenderbali lahir 1997-10-08

# Kalender bulanan
kalenderbali bulan 2026-03

# Ramalan jodoh (2 tanggal lahir)
kalenderbali jodoh 1997-10-08 2002-08-17

# Output JSON (semua perintah mendukung --json)
kalenderbali today --json
kalenderbali date 2026-03-23 --json
kalenderbali hariraya 2026 --json
kalenderbali purnama 2026 --json
```

## 📚 Penggunaan sebagai Library

```go
package main

import (
    "fmt"
    "time"

    "github.com/pandeptwidyaop/kalenderbali-go/pawukon"
    "github.com/pandeptwidyaop/kalenderbali-go/wewaran"
    "github.com/pandeptwidyaop/kalenderbali-go/lunar"
    "github.com/pandeptwidyaop/kalenderbali-go/hariraya"
    "github.com/pandeptwidyaop/kalenderbali-go/dewasa"
    "github.com/pandeptwidyaop/kalenderbali-go/pararasan"
)

func main() {
    date := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)

    // ── Pawukon ──────────────────────────────────────────────────────────────
    pd := pawukon.DayOfCycle(date)
    fmt.Printf("Pawukon hari ke-%d dari 210\n", pd)

    // ── Wewaran (semua 10 sistem minggu) ──────────────────────────────────
    w := wewaran.Calculate(date)
    fmt.Printf("Wuku       : %s (ke-%d)\n", w.WukuName, w.WukuIndex+1)
    fmt.Printf("Saptawara  : %s\n", w.SaptawaraName)
    fmt.Printf("Pancawara  : %s\n", w.PancawaraName)
    fmt.Printf("Triwara    : %s\n", w.TriwaraName)
    fmt.Printf("Caturwara  : %s\n", w.CaturwaraName)
    fmt.Printf("Dasawara   : %s (urip %d)\n", w.DasawaraName, w.Urip)

    // ── Lunar — Sasih & fase bulan ───────────────────────────────────────
    _, sasih := lunar.SasihForDate(date)
    fmt.Printf("Sasih      : %s\n", sasih)

    phases := lunar.PhasesInYear(2026)
    for _, p := range phases {
        _, s := lunar.SasihForDate(p.Date)
        fmt.Printf("%s %s — %s\n", p.Phase, p.Date.Format("2 Jan"), s)
    }

    // ── Hari Raya ────────────────────────────────────────────────────────
    holidays := hariraya.HolidaysInYear(2026)
    for _, h := range holidays {
        fmt.Printf("%s: %s [%s]\n", h.Date.Format("2 Jan"), h.Name, h.Category)
    }

    // Hari raya hari ini
    today := hariraya.ForDate(date)
    for _, h := range today {
        fmt.Printf("  → %s\n", h.Name)
    }

    // 10 hari raya berikutnya
    next10 := hariraya.NextN(date, 10)
    for _, h := range next10 {
        fmt.Printf("  → %s: %s\n", h.Date.Format("2 Jan"), h.Name)
    }

    // ── Dewasa ───────────────────────────────────────────────────────────
    d := dewasa.Calculate(date)
    fmt.Printf("Ingkel     : %s\n", d.Ingkel)
    fmt.Printf("Watek Madya: %s\n", d.WatekMadya)
    fmt.Printf("IsPurnama  : %v\n", d.IsPurnama)
    for _, dw := range d.DewasaList {
        fmt.Printf("  %s: %s — %s\n", dw.Type, dw.Name, dw.Description)
    }

    // ── Pararasan (Laku) ──────────────────────────────────────────────────
    p := pararasan.Calculate(date)
    fmt.Printf("Laku       : %s — %s\n", p.LakunName, p.LakunElement)
    fmt.Printf("            %s\n", p.LakunDesc)
}
```

## 📋 Daftar Perintah CLI

| Perintah | Deskripsi |
|----------|-----------|
| `today` | Info lengkap hari ini |
| `date <YYYY-MM-DD>` | Info lengkap tanggal tertentu |
| `wuku [YYYY-MM-DD]` | Wuku untuk tanggal (default: hari ini) |
| `purnama [tahun]` | Daftar Purnama setahun |
| `tilem [tahun]` | Daftar Tilem setahun |
| `hariraya [tahun]` | Semua hari raya setahun |
| `next [keyword]` | Hari raya berikutnya (default: 10 terdekat) |
| `search <keyword> [tahun]` | Cari hari raya berdasarkan nama |
| `dewasa [YYYY-MM-DD]` | Dewasa Ayu/Ala untuk tanggal |
| `dewasa-ayu [tahun]` | Semua hari baik setahun |
| `dewasa-ala [tahun]` | Semua hari buruk setahun |
| `pararasan [YYYY-MM-DD\|all]` | Laku untuk tanggal atau semua 35 kombinasi |
| `lahir <YYYY-MM-DD>` | Ramalan kelahiran (laku, watak, wuku) |
| `bulan [YYYY-MM]` | Kalender Bali bulanan |
| `jodoh <tgl1> <tgl2>` | Ramalan kecocokan jodoh |

Semua perintah mendukung flag `--json` untuk output machine-readable.

## 🧮 Bagaimana Perhitungannya?

### Pawukon (Siklus 210 Hari)

Semua hari dalam kalender Bali berputar dalam siklus 210 hari (30 wuku × 7 hari). Rumus inti:

```
epoch    = 21 Mei 2000 (Redite Sinta, hari ke-0)
diff     = tanggal − epoch (dalam hari)
pawukon  = ((diff % 210) + 210) % 210
```

### Wewaran (10 Minggu Bersamaan)

Dari posisi pawukon `d`, dihitung 10 sistem minggu sekaligus:

| Sistem | Siklus | Rumus |
|--------|--------|-------|
| Ekawara | 1 | `urip % 2` |
| Dwiwara | 2 | `urip % 2` |
| Triwara | 3 | `d % 3` |
| Caturwara | 4 | `d % 4` (dengan aturan khusus hari 71–72) |
| Pancawara | 5 | `d % 5` |
| Sadwara | 6 | `d % 6` |
| Saptawara | 7 | `d % 7` |
| Astawara | 8 | `d % 8` (dengan aturan khusus hari 71–72) |
| Sangawara | 9 | `(d-3) % 9` (3 hari pertama = Dangu) |
| Dasawara | 10 | `urip − 1` |

**Aturan khusus Dunggulan:** Hari ke-71 dan 72 dalam siklus (wuku Dunggulan) sama-sama "Jaya" (Caturwara) dan "Kala" (Astawara) — satu posisi dipakai dua hari. Ini menyebabkan offset +2 untuk semua hari setelahnya.

### Purnama & Tilem (Algoritma Jean Meeus)

Fase bulan dihitung dengan algoritma astronomi dari buku *Astronomical Algorithms* (2nd ed., 1998). Akurasi ±1 hari untuk tanggal modern.

### Dewasa Ayu/Ala

Rule engine berbasis kombinasi wewaran + wuku + fase bulan. Setiap dewasa adalah fungsi matcher deterministik — tidak ada tabel lookup eksternal.

### Pararasan (Laku)

Tabel 7×5 (Saptawara × Pancawara) yang memetakan setiap kombinasi ke satu dari 12 jenis laku. Diverifikasi terhadap kalenderbali.org.

## 🕉️ Konteks Budaya

Kalender Bali (Kalender Caka) adalah sistem penanggalan **lunisolar** yang digunakan masyarakat Hindu Bali untuk menentukan:

- **Hari Raya** — Galungan (kemenangan dharma), Kuningan, Nyepi (Tahun Baru Saka), Saraswati (Dewi Ilmu), dan puluhan hari raya lainnya
- **Dewasa Ayu** — hari baik untuk upacara, pernikahan, membangun rumah, memulai usaha
- **Dewasa Ala** — hari yang perlu dihindari untuk kegiatan penting
- **Palalintangan** — ramalan watak berdasarkan hari lahir
- **Perjodohan** — kecocokan pasangan berdasarkan urip dan neptu

### Hari Raya Pawukon (Siklus 210 Hari)

Berulang setiap 210 hari, posisinya tetap dalam siklus:

| Hari Raya | Posisi dalam Siklus |
|-----------|---------------------|
| Pagerwesi | Buda Kliwon Sinta (hari ke-3) |
| Tumpek Landep | Saniscara Kliwon Landep (hari ke-13) |
| Tumpek Uduh | Saniscara Kliwon Wariga (hari ke-48) |
| Galungan | Buda Kliwon Dunggulan (hari ke-73) |
| Kuningan | Saniscara Kliwon Kuningan (hari ke-83) |
| Tumpek Krulut | Saniscara Kliwon Krulut (hari ke-118) |
| Tumpek Kandang | Saniscara Kliwon Uye (hari ke-153) |
| Tumpek Wayang | Saniscara Kliwon Wayang (hari ke-188) |
| Saraswati | Saniscara Umanis Watugunung (hari ke-209) |

### Hari Raya Lunar

Bergantung pada siklus bulan, terjadi setiap tahun Masehi pada tanggal berbeda:

- **Purnama** & **Tilem** — tiap sasih (bulan lunar Bali)
- **Siwaratri** — Tilem Kapitu
- **Nyepi** — sehari setelah Tilem Kasanga (Tahun Baru Saka)
- **Ngembak Geni** — sehari setelah Nyepi

## 📖 Referensi

- [kalenderbali.org](https://kalenderbali.org) — I Wayan Nuarsa, Universitas Udayana (sumber verifikasi utama)
- [Pawukon Calendar — Wikipedia](https://en.wikipedia.org/wiki/Pawukon_calendar)
- Jean Meeus, *Astronomical Algorithms*, 2nd ed. (1998) — algoritma fase bulan
- Lontar Wariga — pustaka tradisional Bali

## 📄 Lisensi

MIT License — bebas digunakan, dimodifikasi, dan didistribusikan.

## 🙏 Kredit

Dibuat oleh [Pande Putu Widya Okta Pratama](https://github.com/pandeptwidyaop)

Terima kasih kepada para penyusun kalender Bali yang telah mewariskan ilmu wariga selama berabad-abad. 🙏
