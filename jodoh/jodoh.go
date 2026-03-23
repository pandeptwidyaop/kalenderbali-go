// Package jodoh implements Balinese Wariga-based marriage/compatibility prediction.
// All calculations are pure arithmetic derived from the urip (life value) of birth dates
// across multiple wewaran systems.
package jodoh

import (
	"time"

	"github.com/pandeptwidyaop/kalenderbali-go/pawukon"
	"github.com/pandeptwidyaop/kalenderbali-go/wewaran"
)

// ── Urip helpers ─────────────────────────────────────────────────────────────

// Sadwara urip values (Tungleh=7, Aryang=6, Urukung=5, Paniron=8, Was=9, Maulu=3)
var sadwaraUrip = [6]int{7, 6, 5, 8, 9, 3}

// getUrips returns saptawara index, saptawara urip, pancawara index, pancawara urip,
// and sadwara urip for a given date.
func getUrips(t time.Time) (saptaIdx, uripSapta, pancaIdx, uripPanca, uripSadwa int) {
	d := pawukon.DayOfCycle(t)
	saptaIdx = d % 7
	pancaIdx = d % 5
	sadwaIdx := d % 6
	uripSapta = wewaran.SaptawaraUrip[saptaIdx]
	uripPanca = wewaran.PancawaraUrip[pancaIdx]
	uripSadwa = sadwaraUrip[sadwaIdx]
	return
}

// ── Method 1: Saptawara Compatibility ────────────────────────────────────────

// saptaCompatKey encodes two saptawara indices (0-6 × 0-6) into a single int key.
func saptaCompatKey(a, b int) int { return a*7 + b }

type saptaEntry struct {
	prediction  string
	description string
	isGood      bool
}

// 7×7 Saptawara compatibility matrix (pria row × wanita column).
// Indexed by [pria_saptaIdx][wanita_saptaIdx].
// Names: Redite(0) Soma(1) Anggara(2) Buda(3) Wraspati(4) Sukra(5) Saniscara(6)
var saptaMatrix [7][7]saptaEntry

func init() {
	// Row Redite (0)
	saptaMatrix[0][0] = saptaEntry{"Sakit-sakitan", "Pasangan sering ditimpa penyakit, kurang sehat", false}
	saptaMatrix[0][1] = saptaEntry{"Yuana & Rupawan", "Pasangan tampan/cantik, awet muda, dicintai banyak orang", true}
	saptaMatrix[0][2] = saptaEntry{"Kaya Raya", "Rezeki berlimpah, mudah mendapat kemakmuran", true}
	saptaMatrix[0][3] = saptaEntry{"Harmonis", "Kehidupan bersama penuh keserasian dan ketenangan", true}
	saptaMatrix[0][4] = saptaEntry{"Panjang Umur", "Diberkahi umur panjang dan kesehatan baik", true}
	saptaMatrix[0][5] = saptaEntry{"Banyak Anak", "Dikaruniai banyak keturunan yang sehat", true}
	saptaMatrix[0][6] = saptaEntry{"Berwibawa", "Disegani dan dihormati oleh lingkungan sekitar", true}

	// Row Soma (1)
	saptaMatrix[1][0] = saptaEntry{"Yuana & Rupawan", "Pasangan tampan/cantik, selalu terlihat muda dan menarik", true}
	saptaMatrix[1][1] = saptaEntry{"Sakit-sakitan", "Mudah jatuh sakit, perlu menjaga kesehatan bersama", false}
	saptaMatrix[1][2] = saptaEntry{"Harmonis", "Keluarga hidup rukun, jarang bertengkar", true}
	saptaMatrix[1][3] = saptaEntry{"Kaya Raya", "Usaha bersama berbuah kekayaan yang berlimpah", true}
	saptaMatrix[1][4] = saptaEntry{"Banyak Anak", "Rumah tangga ramai dengan keturunan yang sehat", true}
	saptaMatrix[1][5] = saptaEntry{"Panjang Umur", "Hidup panjang bersama dalam kebahagiaan", true}
	saptaMatrix[1][6] = saptaEntry{"Susah", "Banyak rintangan dan kesulitan dalam berumah tangga", false}

	// Row Anggara (2)
	saptaMatrix[2][0] = saptaEntry{"Kaya Raya", "Rezeki mengalir deras, kehidupan makmur sejahtera", true}
	saptaMatrix[2][1] = saptaEntry{"Harmonis", "Pasangan serasi, saling melengkapi satu sama lain", true}
	saptaMatrix[2][2] = saptaEntry{"Sakit-sakitan", "Rentan penyakit, harus rajin menjaga kesehatan", false}
	saptaMatrix[2][3] = saptaEntry{"Berwibawa", "Memiliki kharisma kuat, dihormati masyarakat", true}
	saptaMatrix[2][4] = saptaEntry{"Susah", "Sering menghadapi hambatan rezeki dan masalah rumah tangga", false}
	saptaMatrix[2][5] = saptaEntry{"Banyak Anak", "Dikaruniai keturunan banyak dan berbakti", true}
	saptaMatrix[2][6] = saptaEntry{"Panjang Umur", "Pasangan diberkahi kesehatan dan umur yang panjang", true}

	// Row Buda (3)
	saptaMatrix[3][0] = saptaEntry{"Harmonis", "Kehidupan damai, saling memahami dan mendukung", true}
	saptaMatrix[3][1] = saptaEntry{"Kaya Raya", "Peruntungan keuangan sangat baik bersama", true}
	saptaMatrix[3][2] = saptaEntry{"Berwibawa", "Pasangan disegani, berpengaruh di masyarakat", true}
	saptaMatrix[3][3] = saptaEntry{"Sakit-sakitan", "Keduanya mudah terkena gangguan kesehatan", false}
	saptaMatrix[3][4] = saptaEntry{"Panjang Umur", "Diberkahi kesehatan prima dan usia panjang", true}
	saptaMatrix[3][5] = saptaEntry{"Susah", "Kehidupan penuh cobaan, perlu kesabaran ekstra", false}
	saptaMatrix[3][6] = saptaEntry{"Banyak Anak", "Rumah tangga penuh canda anak yang berbakti", true}

	// Row Wraspati (4)
	saptaMatrix[4][0] = saptaEntry{"Panjang Umur", "Bersama menjalani hidup panjang dalam keselamatan", true}
	saptaMatrix[4][1] = saptaEntry{"Banyak Anak", "Keturunan banyak dan membawa berkah bagi orang tua", true}
	saptaMatrix[4][2] = saptaEntry{"Susah", "Banyak halangan dan cobaan dalam kehidupan bersama", false}
	saptaMatrix[4][3] = saptaEntry{"Panjang Umur", "Umur panjang, sehat walafiat berdua", true}
	saptaMatrix[4][4] = saptaEntry{"Sakit-sakitan", "Sering bergantian sakit, perlu waspada kesehatan", false}
	saptaMatrix[4][5] = saptaEntry{"Kaya Raya", "Usaha bersama mendatangkan kekayaan berlimpah", true}
	saptaMatrix[4][6] = saptaEntry{"Harmonis", "Rumah tangga tenang dan penuh kasih sayang", true}

	// Row Sukra (5)
	saptaMatrix[5][0] = saptaEntry{"Banyak Anak", "Dikaruniai banyak anak yang sehat dan berbakti", true}
	saptaMatrix[5][1] = saptaEntry{"Panjang Umur", "Hidup panjang berdua dalam kesehatan dan kedamaian", true}
	saptaMatrix[5][2] = saptaEntry{"Banyak Anak", "Keturunan melimpah, rumah tangga hangat dan ramai", true}
	saptaMatrix[5][3] = saptaEntry{"Susah", "Sering terjadi pertentangan, perlu banyak kompromi", false}
	saptaMatrix[5][4] = saptaEntry{"Kaya Raya", "Rezeki bersama terus bertambah dan berkembang", true}
	saptaMatrix[5][5] = saptaEntry{"Sakit-sakitan", "Kedua pasangan rentan penyakit, jaga pola hidup", false}
	saptaMatrix[5][6] = saptaEntry{"Berwibawa", "Pasangan dihormati, punya pengaruh besar", true}

	// Row Saniscara (6)
	saptaMatrix[6][0] = saptaEntry{"Berwibawa", "Dihormati dan disegani oleh lingkungan dan keluarga besar", true}
	saptaMatrix[6][1] = saptaEntry{"Susah", "Banyak rintangan hidup, butuh ketabahan bersama", false}
	saptaMatrix[6][2] = saptaEntry{"Panjang Umur", "Umur panjang dan kesehatan terjaga berdua", true}
	saptaMatrix[6][3] = saptaEntry{"Banyak Anak", "Dikaruniai banyak keturunan yang membawa bahagia", true}
	saptaMatrix[6][4] = saptaEntry{"Harmonis", "Kehidupan serasi, saling mendukung dan mencintai", true}
	saptaMatrix[6][5] = saptaEntry{"Kaya Raya", "Bersama membangun kekayaan dan kemakmuran", true}
	saptaMatrix[6][6] = saptaEntry{"Sakit-sakitan", "Rentan gangguan kesehatan, perlu rajin menjaga tubuh", false}
}

// ── Method 2: Neptu Mod 5 ────────────────────────────────────────────────────
// (handled via neptMod5Entry helper function below)

// ── Method 3: Neptu Mod 4 ────────────────────────────────────────────────────

var neptuMod4 = [4]struct {
	name        string
	description string
	isGood      bool
}{
	{"Punggel — Cerai/Kematian", "Perlu ritual khusus (Pamahayu Karang), rawan perpisahan", false},
	{"Gento — Jarang Anak", "Sulit mendapat keturunan, perlu usaha dan doa lebih", false},
	{"Pati — Banyak Anak", "Dikaruniai banyak keturunan yang sehat dan berbakti", true},
	{"Sugih — Banyak Rezeki", "Pernikahan membawa kemakmuran dan keberlimpahan rezeki", true},
}

// ── Method 4: Pertemuan Neptu Mod 9 ──────────────────────────────────────────

// The two urip values are each taken mod 9 (result 1-9 where 0→9).
// 9×9 matrix encoded as [a-1][b-1] (0-indexed).

type mod9Entry struct {
	prediction  string
	description string
	isGood      bool
}

var mod9Matrix [9][9]mod9Entry

func init() {
	// Indexed [pria_mod9-1][wanita_mod9-1], values 1-9
	data := []struct {
		a, b        int
		prediction  string
		description string
		isGood      bool
	}{
		// Row 1 (pria mod9=1)
		{1, 1, "Saling Mencintai", "Cinta yang tulus dan mendalam antara keduanya", true},
		{1, 2, "Agak Susah", "Ada sedikit hambatan, butuh usaha ekstra", false},
		{1, 3, "Banyak Rezeki", "Bersama mendatangkan keberuntungan finansial", true},
		{1, 4, "Saling Mencintai", "Hubungan penuh kasih sayang yang tulus", true},
		{1, 5, "Hidup Bahagia", "Pasangan saling melengkapi menuju kebahagiaan", true},
		{1, 6, "Agak Susah", "Terkadang timbul ketidakcocokan yang perlu diatasi", false},
		{1, 7, "Banyak Rezeki", "Usaha bersama selalu mendapat kemudahan", true},
		{1, 8, "Hidup Bahagia", "Kehidupan rumah tangga penuh suka cita", true},
		{1, 9, "Agak Susah", "Perlu kesabaran menghadapi perbedaan karakter", false},

		// Row 2 (pria mod9=2)
		{2, 1, "Agak Susah", "Ada gesekan kecil yang perlu diselesaikan bersama", false},
		{2, 2, "Banyak Anak", "Dikaruniai keturunan yang sehat dan berbakti", true},
		{2, 3, "Agak Susah", "Keduanya perlu banyak berkompromi", false},
		{2, 4, "Banyak Anak", "Rumah tangga ramai dan penuh berkah keturunan", true},
		{2, 5, "Agak Susah", "Sering ada perbedaan pendapat, butuh komunikasi", false},
		{2, 6, "Cepat Kaya", "Bersama langsung meraih kemakmuran dengan cepat", true},
		{2, 7, "Agak Susah", "Perjalanan hidup bersama penuh tantangan kecil", false},
		{2, 8, "Cepat Kaya", "Rezeki datang berlimpah begitu bersama", true},
		{2, 9, "Agak Susah", "Butuh penyesuaian diri yang cukup lama", false},

		// Row 3 (pria mod9=3)
		{3, 1, "Banyak Rezeki", "Kehidupan finansial sangat terjamin bersama", true},
		{3, 2, "Agak Susah", "Terkadang situasi kurang mendukung", false},
		{3, 3, "Hidup Bahagia", "Kebahagiaan sejati terjalin erat dalam kebersamaan", true},
		{3, 4, "Agak Susah", "Ada beberapa rintangan yang harus dilalui bersama", false},
		{3, 5, "Banyak Rezeki", "Usaha bersama selalu berbuah hasil melimpah", true},
		{3, 6, "Agak Susah", "Perlu ekstra sabar dalam menjalani kehidupan bersama", false},
		{3, 7, "Hidup Bahagia", "Kebersamaan membawa kebahagiaan sejati", true},
		{3, 8, "Agak Susah", "Terkadang kurang sejalan, perlu banyak diskusi", false},
		{3, 9, "Banyak Rezeki", "Pernikahan memperkuat fondasi finansial keduanya", true},

		// Row 4 (pria mod9=4)
		{4, 1, "Saling Mencintai", "Kasih sayang yang dalam dan bertahan lama", true},
		{4, 2, "Banyak Anak", "Keturunan banyak menjadi berkah terbesar", true},
		{4, 3, "Agak Susah", "Ada sedikit ketidakserasian yang perlu diatasi", false},
		{4, 4, "Saling Mencintai", "Hubungan penuh kehangatan dan perhatian", true},
		{4, 5, "Banyak Anak", "Dikaruniai anak-anak yang sehat dan pintar", true},
		{4, 6, "Agak Susah", "Keduanya perlu belajar memahami perbedaan", false},
		{4, 7, "Saling Mencintai", "Cinta yang tulus menjadi pondasi kuat", true},
		{4, 8, "Banyak Anak", "Rumah tangga penuh canda tawa anak cucu", true},
		{4, 9, "Agak Susah", "Perlu usaha lebih untuk mempertahankan keharmonisan", false},

		// Row 5 (pria mod9=5)
		{5, 1, "Hidup Bahagia", "Kebahagiaan sejati ada dalam kebersamaan ini", true},
		{5, 2, "Agak Susah", "Sering ada selisih paham yang perlu diselesaikan", false},
		{5, 3, "Banyak Rezeki", "Keberuntungan finansial terus datang bersama", true},
		{5, 4, "Banyak Anak", "Anak-anak menjadi sumber kebahagiaan terbesar", true},
		{5, 5, "Beruntung Terus", "Keberuntungan hadir di setiap aspek kehidupan", true},
		{5, 6, "Agak Susah", "Butuh penyesuaian lebih dalam menjalani hidup bersama", false},
		{5, 7, "Banyak Rezeki", "Rezeki datang dari berbagai arah begitu bersatu", true},
		{5, 8, "Beruntung Terus", "Keberuntungan tanpa henti menyertai pasangan ini", true},
		{5, 9, "Agak Susah", "Ada beberapa hambatan yang harus dihadapi bersama", false},

		// Row 6 (pria mod9=6)
		{6, 1, "Agak Susah", "Terkadang arah hidup keduanya berbeda", false},
		{6, 2, "Cepat Kaya", "Kemakmuran datang lebih cepat dari yang diharapkan", true},
		{6, 3, "Agak Susah", "Perlu banyak penyesuaian dalam kehidupan bersama", false},
		{6, 4, "Agak Susah", "Keduanya harus banyak berkompromi soal perbedaan", false},
		{6, 5, "Agak Susah", "Ada gesekan yang perlu disikapi dengan bijaksana", false},
		{6, 6, "Cepat Kaya", "Bersama langsung bisa meraih kekayaan dengan mudah", true},
		{6, 7, "Agak Susah", "Butuh kesabaran ekstra dalam kebersamaan", false},
		{6, 8, "Cepat Kaya", "Rezeki berlimpah segera hadir setelah bersatu", true},
		{6, 9, "Agak Susah", "Perlu kerja keras bersama untuk mencapai tujuan", false},

		// Row 7 (pria mod9=7)
		{7, 1, "Banyak Rezeki", "Pernikahan ini membawa limpahan rezeki", true},
		{7, 2, "Agak Susah", "Ada beberapa cobaan yang datang bergantian", false},
		{7, 3, "Hidup Bahagia", "Kebahagiaan sejati menjadi bagian keseharian", true},
		{7, 4, "Saling Mencintai", "Cinta mendalam yang tumbuh semakin kuat", true},
		{7, 5, "Banyak Rezeki", "Usaha bersama selalu berbuah hasil nyata", true},
		{7, 6, "Agak Susah", "Terkadang terjadi ketidakserasian yang perlu dijembatani", false},
		{7, 7, "Hidup Bahagia", "Kehidupan bersama penuh keberkahan dan suka cita", true},
		{7, 8, "Agak Susah", "Ada perbedaan prinsip yang perlu diselaraskan", false},
		{7, 9, "Banyak Rezeki", "Keberuntungan finansial selalu hadir bersama", true},

		// Row 8 (pria mod9=8)
		{8, 1, "Hidup Bahagia", "Kehidupan bersama penuh kebahagiaan yang nyata", true},
		{8, 2, "Cepat Kaya", "Kemakmuran datang lebih cepat dari perkiraan", true},
		{8, 3, "Agak Susah", "Keduanya perlu ekstra sabar menghadapi ujian", false},
		{8, 4, "Banyak Anak", "Dikaruniai keturunan yang membawa kebahagiaan", true},
		{8, 5, "Beruntung Terus", "Semua aspek kehidupan selalu beruntung", true},
		{8, 6, "Cepat Kaya", "Kekayaan datang dengan mudah begitu bersatu", true},
		{8, 7, "Agak Susah", "Butuh penyesuaian karakter yang memerlukan waktu", false},
		{8, 8, "Hidup Bahagia", "Bersama menjalani kehidupan penuh kebahagiaan", true},
		{8, 9, "Agak Susah", "Terkadang ada ketidakselarasan yang perlu diatasi", false},

		// Row 9 (pria mod9=9)
		{9, 1, "Agak Susah", "Ada hambatan yang perlu dihadapi dengan sabar", false},
		{9, 2, "Agak Susah", "Keduanya butuh banyak penyesuaian", false},
		{9, 3, "Banyak Rezeki", "Rezeki selalu ada untuk pasangan ini", true},
		{9, 4, "Agak Susah", "Perlu usaha bersama untuk mencapai keserasian", false},
		{9, 5, "Agak Susah", "Ada beberapa rintangan yang menguji kesabaran", false},
		{9, 6, "Agak Susah", "Butuh komitmen kuat untuk melewati perbedaan", false},
		{9, 7, "Banyak Rezeki", "Bersama selalu mendapat limpahan rezeki", true},
		{9, 8, "Agak Susah", "Terkadang tujuan hidup keduanya berbeda arah", false},
		{9, 9, "Susah Rezeki", "Keuangan sering terhambat, perlu kerja keras berlipat", false},
	}
	for _, e := range data {
		mod9Matrix[e.a-1][e.b-1] = mod9Entry{e.prediction, e.description, e.isGood}
	}
}

// ── Method 5: Tri Pramana (Sodasa Rsi) ───────────────────────────────────────

// triPramanaResults maps sisa (1-16) to a result.
// sisa = (total_urip_pria + total_urip_wanita) % 16, where 0→16
var triPramanaResults = map[int]struct {
	name        string
	description string
	isGood      bool
}{
	1:  {"Bimbang", "Keduanya sering ragu dan tidak pasti dalam mengambil keputusan", false},
	2:  {"Kurang Baik", "Ada hambatan-hambatan kecil dalam kehidupan bersama", false},
	3:  {"Cukup Baik", "Kehidupan cukup lancar meski sesekali ada tantangan", true},
	4:  {"Banyak Cobaan", "Sering diuji dengan berbagai cobaan yang menguatkan", false},
	5:  {"Harmonis", "Keluarga hidup dalam keserasian dan saling mendukung", true},
	6:  {"Kurang Rezeki", "Perlu usaha lebih keras untuk mencukupi kebutuhan", false},
	7:  {"Cukup Selamat", "Kehidupan terlindungi dari bahaya besar", true},
	8:  {"Sering Bertengkar", "Keduanya sering berbeda pendapat dan berselisih", false},
	9:  {"Tenteram", "Rumah tangga penuh kedamaian dan ketentraman", true},
	10: {"Berwibawa", "Pasangan dihormati dan memiliki pengaruh di masyarakat", true},
	11: {"Sukses", "Karir dan usaha bersama selalu menuai kesuksesan", true},
	12: {"Rezeki Lancar", "Arus keuangan selalu lancar dan terpenuhi", true},
	13: {"Panjang Umur", "Keduanya dikaruniai kesehatan dan umur yang panjang", true},
	14: {"Kurang Bahagia", "Kebahagiaan kadang sulit dirasa, perlu lebih bersyukur", false},
	15: {"Cukup Rezeki", "Kebutuhan hidup selalu tercukupi dengan baik", true},
	16: {"Bahagia Sempurna", "Kebahagiaan sejati dan menyeluruh dalam segala aspek", true},
}

// ── Method 6: Ramalan 5 Tahun ────────────────────────────────────────────────

var ramalan5Tahun = [5]struct {
	name        string
	description string
	isGood      bool
}{
	{"Kedudukan Naik", "Kedudukan dan status sosial pasangan terus meningkat", true},
	{"Sri — Sejahtera", "Kehidupan penuh kesejahteraan dan limpahan berkah", true},
	{"Gedong — Berkecukupan", "Rumah tangga selalu berkecukupan, tidak pernah kekurangan", true},
	{"Pete — Bertengkar", "Sering terjadi pertengkaran, perlu banyak komunikasi", false},
	{"Pati — Masalah", "Banyak masalah yang menghampiri, butuh ketabahan", false},
}

// ── JodohResult ───────────────────────────────────────────────────────────────

// JodohResult holds the prediction result of one compatibility method.
type JodohResult struct {
	Method      string
	Score       int
	Prediction  string
	Description string
	IsGood      bool
}

// ── CheckJodoh ────────────────────────────────────────────────────────────────

// CheckJodoh computes all 6 Balinese Wariga compatibility methods for two birth dates.
// birthPria is the male's birth date, birthWanita is the female's birth date.
func CheckJodoh(birthPria, birthWanita time.Time) []JodohResult {
	// Gather urip values for each person
	saptaP, uripSaptaP, pancaP, uripPancaP, uripSadwaP := getUrips(birthPria)
	saptaW, uripSaptaW, pancaW, uripPancaW, uripSadwaW := getUrips(birthWanita)

	_ = pancaP // avoid unused var if not used directly below
	_ = pancaW

	var results []JodohResult

	// ── Method 1: Saptawara Compatibility ──
	e1 := saptaMatrix[saptaP][saptaW]
	results = append(results, JodohResult{
		Method:      "Saptawara",
		Score:       saptaP*7 + saptaW,
		Prediction:  e1.prediction,
		Description: e1.description,
		IsGood:      e1.isGood,
	})

	// ── Method 2: Neptu Mod 5 ──
	neptTotalP := uripPancaP + uripSaptaP
	neptTotalW := uripPancaW + uripSaptaW
	mod5 := (neptTotalP + neptTotalW) % 5
	e2 := neptMod5Entry(mod5)
	results = append(results, JodohResult{
		Method:      "Neptu Mod 5",
		Score:       mod5,
		Prediction:  e2.name,
		Description: e2.description,
		IsGood:      e2.isGood,
	})

	// ── Method 3: Neptu Mod 4 ──
	mod4 := (neptTotalP + neptTotalW) % 4
	e3 := neptuMod4[mod4]
	results = append(results, JodohResult{
		Method:      "Neptu Mod 4",
		Score:       mod4,
		Prediction:  e3.name,
		Description: e3.description,
		IsGood:      e3.isGood,
	})

	// ── Method 4: Pertemuan Neptu Mod 9 ──
	mod9P := neptTotalP % 9
	if mod9P == 0 {
		mod9P = 9
	}
	mod9W := neptTotalW % 9
	if mod9W == 0 {
		mod9W = 9
	}
	e4 := mod9Matrix[mod9P-1][mod9W-1]
	results = append(results, JodohResult{
		Method:      "Pertemuan Neptu (Mod 9)",
		Score:       mod9P*9 + mod9W,
		Prediction:  e4.prediction,
		Description: e4.description,
		IsGood:      e4.isGood,
	})

	// ── Method 5: Tri Pramana (Sodasa Rsi) ──
	totalTriP := uripPancaP + uripSadwaP + uripSaptaP
	totalTriW := uripPancaW + uripSadwaW + uripSaptaW
	sisa := (totalTriP + totalTriW) % 16
	if sisa == 0 {
		sisa = 16
	}
	e5 := triPramanaResults[sisa]
	results = append(results, JodohResult{
		Method:      "Tri Pramana (Sodasa Rsi)",
		Score:       sisa,
		Prediction:  e5.name,
		Description: e5.description,
		IsGood:      e5.isGood,
	})

	// ── Method 6: Ramalan 5 Tahun ──
	totalBoth := neptTotalP + neptTotalW
	mod5yr := totalBoth % 5
	e6 := ramalan5Tahun[mod5yr]
	results = append(results, JodohResult{
		Method:      "Ramalan 5 Tahun",
		Score:       mod5yr,
		Prediction:  e6.name,
		Description: e6.description,
		IsGood:      e6.isGood,
	})

	return results
}

// neptMod5Entry is a helper to get mod5 entries (avoids array struct literal issues).
func neptMod5Entry(mod5 int) struct {
	name        string
	description string
	isGood      bool
} {
	entries := []struct {
		name        string
		description string
		isGood      bool
	}{
		{"Lungguh — Kedudukan Tinggi", "Pasangan akan mendapat kedudukan dan kehormatan di masyarakat", true},
		{"Sri — Rezeki Berlimpah", "Pernikahan membawa keberuntungan, rezeki mengalir deras", true},
		{"Dana — Keuangan Baik", "Kondisi keuangan rumah tangga selalu terpenuhi dengan baik", true},
		{"Lara — Susah", "Pasangan sering ditimpa kesusahan dan penderitaan", false},
		{"Pati — Sengsara", "Pernikahan berpotensi membawa kesengsaraan bagi keduanya", false},
	}
	return entries[mod5]
}

// SaptawaraName returns the Saptawara name for a given birth date (for display).
func SaptawaraName(t time.Time) string {
	d := pawukon.DayOfCycle(t)
	return wewaran.SaptawaraNames[d%7]
}

// PancawaraName returns the Pancawara name for a given birth date (for display).
func PancawaraName(t time.Time) string {
	d := pawukon.DayOfCycle(t)
	return wewaran.PancawaraNames[d%5]
}

// UripTotal returns uripPancawara + uripSaptawara for a date (used for display).
func UripTotal(t time.Time) int {
	_, uripSapta, _, uripPanca, _ := getUrips(t)
	return uripSapta + uripPanca
}
