// Package surat merender surel HTML bermerek yang seragam lintas service
// GoNgaji: kepala wordmark, kartu isi, kotak kode, tombol CTA, catatan, dan
// kaki "surel otomatis" — beserta versi text/plain untuk multipart.
//
// HTML surel bukan HTML web: klien surel (terutama Outlook & Gmail) memangkas
// <style>, mengabaikan flex/grid, dan memblokir JavaScript sepenuhnya. Karena
// itu cetakannya tabel + CSS inline, dan "interaktif" dalam surel berarti
// tombol tautan (CTA) — bukan skrip.
//
// Identitas merek adalah parameter (bukan konstanta) supaya beberapa produk
// pada satu akun pengirim bisa berbagi struktur tanpa berbagi warna/nama.
// Package ini tidak membaca env — service pemanggil yang memutuskan mereknya.
//
// Contoh:
//
//	m := surat.GoNgaji("https://learning.gongaji.id")
//	s := surat.Surat{
//		Judul:    "Termin 2 jatuh tempo",
//		Teks:     "Assalamu'alaikum Budi,\n\nTermin 2 jatuh tempo 28 Agustus.",
//		CTALabel: "Lihat Tagihan",
//		CTAURL:   "https://learning.gongaji.id/tagihan",
//	}
//	msg := mailer.Message{To: []string{to}, Subject: s.Judul,
//		HTMLBody: s.HTML(m), TextBody: s.TeksPolos()}
package surat

import (
	"html"
	"strings"
)

// Merek adalah identitas visual pengirim: nama, warna, dan tautan kaki surel.
// Nol value TIDAK siap pakai — bangun lewat GoNgaji() atau isi lengkap sendiri.
type Merek struct {
	Nama    string // wordmark di pita kepala, mis. "Go Ngaji"
	Tagline string // teks kecil di samping wordmark, mis. "Belajar Ngaji"

	WarnaPrimer string // tombol CTA + tautan, mis. "#5145b9"
	WarnaPekat  string // pita kepala + teks kode, mis. "#3f3596"
	WarnaAksen  string // garis tipis di bawah kepala, mis. "#ffc748"

	// BasisURL menautkan wordmark di kaki surel; TautanBantuan (opsional)
	// menambah baris "Butuh bantuan?".
	BasisURL      string
	TautanBantuan string
}

// GoNgaji memulangkan merek bawaan Go Ngaji (indigo + emas, palet web santri).
// basisURL = asal web publik produk; tautan bantuannya <basisURL>/faq.
func GoNgaji(basisURL string) Merek {
	basisURL = strings.TrimRight(strings.TrimSpace(basisURL), "/")
	return Merek{
		Nama:          "Go Ngaji",
		Tagline:       "Belajar Ngaji",
		WarnaPrimer:   "#5145b9",
		WarnaPekat:    "#3f3596",
		WarnaAksen:    "#ffc748",
		BasisURL:      basisURL,
		TautanBantuan: basisURL + "/faq",
	}
}

// Surat adalah konten satu surel. Semua bidang opsional kecuali Judul; bagian
// yang kosong tidak dirender.
type Surat struct {
	// Judul tampil sebagai kepala kartu (subjek surel diatur pemanggil,
	// biasanya sama).
	Judul string
	// Teks = isi polos; baris kosong memisahkan paragraf. Di-escape otomatis —
	// pemanggil TIDAK perlu (dan tidak boleh) menyisipkan HTML di sini.
	Teks string
	// ParagrafHTML dipakai pemanggil yang butuh markup di isi (mis. <strong>).
	// Bila terisi, Teks hanya dipakai untuk versi teks-polos multipart; bila
	// Teks kosong, versi polos diturunkan dari ParagrafHTML dengan tag dilucuti.
	ParagrafHTML []string
	// Kode ditampilkan besar di kotak tersendiri (OTP, kode unduh, dsb).
	Kode string
	// CTALabel+CTAURL menjadi tombol; kosong = tanpa tombol. Tautan mentahnya
	// ikut dicetak di bawah tombol — sebagian klien surel menggagalkan tombol
	// bergaya tapi selalu merender tautan teks.
	CTALabel string
	CTAURL   string
	// Catatan tampil kecil & redup di bawah isi (masa berlaku kode, "abaikan
	// bila bukan kamu", dsb).
	Catatan string
}

// Warna netral struktur (kanvas, kartu, teks) bukan bagian identitas merek —
// konstanta, bukan bidang Merek, supaya dua produk beda warna tetap terasa
// satu keluarga.
const (
	warnaLatar   = "#f1f3f9"
	warnaKartu   = "#ffffff"
	warnaGaris   = "#e3e7f0"
	warnaTinta   = "#1f2a37"
	warnaIsi     = "#4a5568"
	warnaRedup   = "#8792a4"
	warnaKodeBg  = "#ecebfb"
	fontTumpukan = "-apple-system,'Segoe UI',Roboto,Helvetica,Arial,sans-serif"
)

// HTML merender surat menjadi dokumen surel utuh dengan identitas m.
func (s Surat) HTML(m Merek) string {
	var b strings.Builder
	b.Grow(4096)

	// Preheader: cuplikan yang ditampilkan klien surel di daftar masuk, di
	// samping subjek. Tanpa ini yang tampil adalah teks pertama dokumen —
	// yaitu nama merek di pita kepala, yang tak menceritakan apa-apa.
	pre := s.Teks
	if pre == "" && len(s.ParagrafHTML) > 0 {
		pre = lucutiTag(s.ParagrafHTML[0])
	}
	if len(pre) > 140 {
		pre = pre[:140]
	}

	b.WriteString(`<!doctype html><html lang="id"><head><meta charset="utf-8">` +
		`<meta name="viewport" content="width=device-width,initial-scale=1">` +
		`<meta name="color-scheme" content="light"><title>` + html.EscapeString(s.Judul) + `</title></head>` +
		`<body style="margin:0;padding:0;background:` + warnaLatar + `;">` +
		`<div style="display:none;max-height:0;overflow:hidden;">` + html.EscapeString(pre) + `</div>` +
		`<table role="presentation" width="100%" cellpadding="0" cellspacing="0" style="background:` + warnaLatar + `;">` +
		`<tr><td align="center" style="padding:32px 16px;">` +
		`<table role="presentation" width="600" cellpadding="0" cellspacing="0" style="max-width:600px;width:100%;">`)

	// Pita kepala: wordmark di warna pekat + garis aksen tipis di bawahnya.
	b.WriteString(`<tr><td style="background:` + m.WarnaPekat + `;border-radius:16px 16px 0 0;padding:22px 32px;">` +
		`<span style="font:800 20px ` + fontTumpukan + `;color:#ffffff;letter-spacing:.02em;">` + html.EscapeString(m.Nama) + `</span>`)
	if m.Tagline != "" {
		b.WriteString(`<span style="font:600 11px ` + fontTumpukan + `;color:#ffffffb3;letter-spacing:.14em;text-transform:uppercase;">&nbsp;&nbsp;·&nbsp;&nbsp;` +
			html.EscapeString(m.Tagline) + `</span>`)
	}
	b.WriteString(`</td></tr>` +
		`<tr><td style="background:` + m.WarnaAksen + `;height:4px;font-size:0;line-height:0;">&nbsp;</td></tr>`)

	// Kartu isi.
	b.WriteString(`<tr><td style="background:` + warnaKartu + `;border:1px solid ` + warnaGaris + `;border-top:0;border-radius:0 0 16px 16px;padding:32px;">`)
	b.WriteString(`<h1 style="margin:0 0 14px;font:700 21px/1.35 ` + fontTumpukan + `;color:` + warnaTinta + `;">` + html.EscapeString(s.Judul) + `</h1>`)

	for _, p := range s.paragrafSemua() {
		b.WriteString(`<p style="margin:0 0 12px;font:400 15px/1.65 ` + fontTumpukan + `;color:` + warnaIsi + `;">` + p + `</p>`)
	}

	if s.Kode != "" {
		b.WriteString(`<table role="presentation" width="100%" cellpadding="0" cellspacing="0"><tr>` +
			`<td align="center" style="background:` + warnaKodeBg + `;border-radius:12px;padding:18px;margin:0;">` +
			`<span style="font:700 28px 'Courier New',Courier,monospace;color:` + m.WarnaPekat + `;letter-spacing:.35em;">` +
			html.EscapeString(s.Kode) + `</span></td></tr></table>` +
			`<div style="height:16px;font-size:0;">&nbsp;</div>`)
	}

	if s.CTALabel != "" && s.CTAURL != "" {
		u := html.EscapeString(s.CTAURL)
		b.WriteString(`<table role="presentation" cellpadding="0" cellspacing="0" style="margin:8px 0 6px;"><tr>` +
			`<td style="background:` + m.WarnaPrimer + `;border-radius:12px;">` +
			`<a href="` + u + `" style="display:inline-block;padding:13px 30px;font:600 15px ` + fontTumpukan + `;color:#ffffff;text-decoration:none;border-radius:12px;">` +
			html.EscapeString(s.CTALabel) + `</a></td></tr></table>` +
			// Tautan mentah untuk klien yang mematikan gaya tombol.
			`<p style="margin:6px 0 0;font:400 12px/1.6 ` + fontTumpukan + `;color:` + warnaRedup + `;word-break:break-all;">` +
			`Tombol tidak berfungsi? Buka tautan ini: <a href="` + u + `" style="color:` + m.WarnaPrimer + `;">` + u + `</a></p>`)
	}

	if s.Catatan != "" {
		b.WriteString(`<p style="margin:18px 0 0;padding-top:14px;border-top:1px solid ` + warnaGaris + `;font:400 13px/1.6 ` + fontTumpukan + `;color:` + warnaRedup + `;">` +
			html.EscapeString(s.Catatan) + `</p>`)
	}

	b.WriteString(`</td></tr>`)

	// Kaki: identitas + pengingat surel otomatis.
	b.WriteString(`<tr><td style="padding:20px 12px 0;" align="center">` +
		`<p style="margin:0;font:400 12px/1.7 ` + fontTumpukan + `;color:` + warnaRedup + `;">` +
		`Surel otomatis dari <a href="` + html.EscapeString(m.BasisURL) + `" style="color:` + m.WarnaPrimer + `;text-decoration:none;font-weight:600;">` +
		html.EscapeString(m.Nama) + `</a> — mohon tidak membalas.`)
	if m.TautanBantuan != "" {
		b.WriteString(`<br>Butuh bantuan? Kunjungi <a href="` + html.EscapeString(m.TautanBantuan) + `" style="color:` + m.WarnaPrimer + `;">Tanya Jawab</a>.`)
	}
	b.WriteString(`</p></td></tr>`)

	b.WriteString(`</table></td></tr></table></body></html>`)
	return b.String()
}

// TeksPolos merender versi text/plain untuk bagian multipart — klien tanpa HTML
// dan penyaring spam sama-sama menghargainya.
func (s Surat) TeksPolos() string {
	var b strings.Builder
	if s.Teks != "" {
		b.WriteString(s.Teks)
	} else {
		for i, p := range s.ParagrafHTML {
			if i > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString(lucutiTag(p))
		}
	}
	if s.Kode != "" {
		b.WriteString("\n\nKode: " + s.Kode)
	}
	if s.CTALabel != "" && s.CTAURL != "" {
		b.WriteString("\n\n" + s.CTALabel + ": " + s.CTAURL)
	}
	if s.Catatan != "" {
		b.WriteString("\n\n" + s.Catatan)
	}
	return b.String()
}

// paragrafSemua: isi kartu sebagai potongan HTML per paragraf, apa pun bentuk
// masukannya (Teks polos di-escape + dipecah di baris kosong; ParagrafHTML
// dipakai apa adanya).
func (s Surat) paragrafSemua() []string {
	if len(s.ParagrafHTML) > 0 {
		return s.ParagrafHTML
	}
	potongan := strings.Split(strings.ReplaceAll(s.Teks, "\r\n", "\n"), "\n\n")
	out := make([]string, 0, len(potongan))
	for _, p := range potongan {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, strings.ReplaceAll(html.EscapeString(p), "\n", "<br>"))
	}
	return out
}

// lucutiTag membuang tag HTML sederhana untuk versi teks-polos/preheader.
// Bukan sanitiser umum — masukannya markup milik service sendiri.
func lucutiTag(s string) string {
	var b strings.Builder
	dalam := false
	for _, r := range s {
		switch {
		case r == '<':
			dalam = true
		case r == '>':
			dalam = false
		case !dalam:
			b.WriteRune(r)
		}
	}
	return html.UnescapeString(strings.TrimSpace(b.String()))
}
