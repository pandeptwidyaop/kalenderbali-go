// Package pararasan implements Laku (daily fortune/tendency) calculation
// based on the 35 fixed combinations of Saptawara × Pancawara.
//
// Each combination maps deterministically to one of 9 Laku types,
// each associated with a natural element and description.
package pararasan

import (
	"time"

	"github.com/pandeptwidyaop/kalenderbali-go/pawukon"
	"github.com/pandeptwidyaop/kalenderbali-go/wewaran"
)

// Lakuname is the name of a Laku.
type LakunName string

const (
	LakunBulan       LakunName = "Bulan"
	LakunBumi        LakunName = "Bumi"
	LakunAngin       LakunName = "Angin"
	LakunSurya       LakunName = "Surya"
	LakunPanditaSakti LakunName = "Pandita Sakti"
	LakunArasKembang LakunName = "Aras Kembang"
	LakunArasTuding  LakunName = "Aras Tuding"
	LakunApi         LakunName = "Api"
	LakunBintang     LakunName = "Bintang"
	LakunToya        LakunName = "Toya"
	LakunPandita     LakunName = "Pandita"
	LakunPretiwi     LakunName = "Pretiwi"
)

// Laku holds the full information for a Laku type.
type Laku struct {
	Name        LakunName
	Element     string // Natural element
	Description string
}

// lakunDescriptions maps each Laku name to its metadata.
var lakunDescriptions = map[LakunName]Laku{
	LakunBulan: {
		Name:        LakunBulan,
		Element:     "Bulan (Moon)",
		Description: "Hari dengan laku Bulan: ketenangan, keindahan, perasaan halus, cocok untuk seni dan meditasi. Energi lembut namun berpengaruh kuat.",
	},
	LakunBumi: {
		Name:        LakunBumi,
		Element:     "Bumi (Earth)",
		Description: "Hari dengan laku Bumi: kestabilan, kesabaran, kesuburan. Baik untuk pekerjaan yang membutuhkan ketekunan dan fondasi kuat.",
	},
	LakunAngin: {
		Name:        LakunAngin,
		Element:     "Angin (Wind)",
		Description: "Hari dengan laku Angin: pergerakan, komunikasi, penyebaran. Cocok untuk perjalanan, negosiasi, dan menyebarkan informasi.",
	},
	LakunSurya: {
		Name:        LakunSurya,
		Element:     "Surya (Sun)",
		Description: "Hari dengan laku Surya: kepemimpinan, kejayaan, kekuatan. Sangat baik untuk memulai usaha besar dan tampil di depan publik.",
	},
	LakunPanditaSakti: {
		Name:        LakunPanditaSakti,
		Element:     "Pandita Sakti (Holy Sage)",
		Description: "Hari dengan laku Pandita Sakti: kebijaksanaan spiritual yang agung, kekuatan batin luar biasa. Baik untuk ritual tinggi dan keputusan penting.",
	},
	LakunArasKembang: {
		Name:        LakunArasKembang,
		Element:     "Aras Kembang (Flower Touch)",
		Description: "Hari dengan laku Aras Kembang: keindahan, kelembutan, kasih sayang. Cocok untuk kesenian, cinta, dan acara yang membutuhkan keanggunan.",
	},
	LakunArasTuding: {
		Name:        LakunArasTuding,
		Element:     "Aras Tuding (Pointing Touch)",
		Description: "Hari dengan laku Aras Tuding: ketegasan, kejelasan arah. Baik untuk pengambilan keputusan dan menentukan tujuan.",
	},
	LakunApi: {
		Name:        LakunApi,
		Element:     "Api (Fire)",
		Description: "Hari dengan laku Api: semangat membara, transformasi, keberanian. Cocok untuk aksi berani dan perubahan, namun hati-hati terhadap konflik.",
	},
	LakunBintang: {
		Name:        LakunBintang,
		Element:     "Bintang (Star)",
		Description: "Hari dengan laku Bintang: harapan, panduan, inspirasi. Baik untuk perencanaan jangka panjang dan menemukan arah baru.",
	},
	LakunToya: {
		Name:        LakunToya,
		Element:     "Toya (Water)",
		Description: "Hari dengan laku Toya: keluwesan, pemurnian, adaptasi. Cocok untuk penyembuhan, pembersihan spiritual, dan penyesuaian rencana.",
	},
	LakunPandita: {
		Name:        LakunPandita,
		Element:     "Pandita (Sage)",
		Description: "Hari dengan laku Pandita: kebijaksanaan, pengetahuan, pengajaran. Baik untuk belajar, mengajar, dan konsultasi.",
	},
	LakunPretiwi: {
		Name:        LakunPretiwi,
		Element:     "Pretiwi (Earth Goddess)",
		Description: "Hari dengan laku Pretiwi: kemakmuran bumi, kesuburan, rezeki dari tanah. Baik untuk pertanian, properti, dan memperkuat fondasi kehidupan.",
	},
}

// lakunTable is the 7×5 lookup table indexed by [saptawaraIdx][pancawaraIdx].
// Saptawara: 0=Redite, 1=Soma, 2=Anggara, 3=Buda, 4=Wraspati, 5=Sukra, 6=Saniscara
// Pancawara: 0=Paing, 1=Pon, 2=Wage, 3=Keliwon, 4=Umanis
var lakunTable = [7][5]LakunName{
	// Redite
	{LakunBulan, LakunBumi, LakunAngin, LakunSurya, LakunPanditaSakti},
	// Soma
	{LakunBumi, LakunSurya, LakunApi, LakunArasKembang, LakunArasTuding},
	// Anggara
	{LakunApi, LakunBintang, LakunBumi, LakunToya, LakunApi},
	// Buda
	{LakunToya, LakunBulan, LakunAngin, LakunSurya, LakunBintang},
	// Wraspati
	{LakunSurya, LakunSurya, LakunAngin, LakunPandita, LakunBintang},
	// Sukra
	{LakunSurya, LakunBintang, LakunPanditaSakti, LakunBulan, LakunToya},
	// Saniscara
	{LakunBumi, LakunToya, LakunApi, LakunPretiwi, LakunAngin},
}

// Result holds the Pararasan (Laku) result for a date.
type Result struct {
	Date           time.Time
	SaptawaraName  string
	PancawaraName  string
	LakunName      LakunName
	LakunElement   string
	LakunDesc      string
}

// Calculate returns the Pararasan (Laku) for a given date.
func Calculate(t time.Time) Result {
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	d := pawukon.DayOfCycle(t)

	saptaIdx := d % 7
	pancaIdx := d % 5

	lkName := lakunTable[saptaIdx][pancaIdx]
	lk := lakunDescriptions[lkName]

	return Result{
		Date:          t,
		SaptawaraName: wewaran.SaptawaraNames[saptaIdx],
		PancawaraName: wewaran.PancawaraNames[pancaIdx],
		LakunName:     lkName,
		LakunElement:  lk.Element,
		LakunDesc:     lk.Description,
	}
}

// All35 returns all 35 possible Saptawara×Pancawara combinations with their Laku,
// in order (Redite-Paing to Saniscara-Umanis).
func All35() []Result {
	results := make([]Result, 0, 35)
	// Use a base date and iterate through 35 days starting from a point
	// where saptawara=0(Redite) and pancawara=0(Paing).
	// Epoch 2000-05-21 is Redite + Paing (day 0: sapta=0%7=0, panca=0%5=0)
	base := time.Date(2000, 5, 21, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 35; i++ {
		d := base.AddDate(0, 0, i)
		results = append(results, Calculate(d))
	}
	return results
}

// LakunInfo returns the full Laku metadata for a given LakunName.
func LakunInfo(name LakunName) (Laku, bool) {
	l, ok := lakunDescriptions[name]
	return l, ok
}
