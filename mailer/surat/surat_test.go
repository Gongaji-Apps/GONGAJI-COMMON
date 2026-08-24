package surat

import (
	"strings"
	"testing"
)

func merekUji() Merek { return GoNgaji("https://learning.gongaji.id/") }

// Isi surel datang dari data (nama pengguna, judul pengumuman admin) — cetakan
// harus meng-escape-nya, bukan mempercayainya.
func TestHTMLMengEscapeIsi(t *testing.T) {
	s := Surat{Judul: `Halo <b>&`, Teks: "baris <script>satu</script>\n\nparagraf & dua"}
	h := s.HTML(merekUji())
	if strings.Contains(h, "<script>") {
		t.Fatalf("isi tidak di-escape: %q", h)
	}
	for _, mau := range []string{"Halo &lt;b&gt;&amp;", "baris &lt;script&gt;", "paragraf &amp; dua"} {
		if !strings.Contains(h, mau) {
			t.Fatalf("escape hilang, mencari %q", mau)
		}
	}
}

func TestHTMLTombolDanKode(t *testing.T) {
	s := Surat{Judul: "J", Teks: "isi", Kode: "AB12CD", CTALabel: "Buka", CTAURL: "https://x/y?a=1&b=2"}
	h := s.HTML(merekUji())
	if !strings.Contains(h, "AB12CD") {
		t.Fatal("kode tidak dirender")
	}
	// URL di atribut href harus ter-escape (& -> &amp;) dan tautan mentahnya ikut dicetak.
	if !strings.Contains(h, `href="https://x/y?a=1&amp;b=2"`) {
		t.Fatal("href CTA tidak ter-escape / hilang")
	}
	if strings.Count(h, "https://x/y") < 2 {
		t.Fatal("tautan mentah cadangan tidak dicetak di bawah tombol")
	}
}

// Tanpa CTA dan kode: tak boleh ada sisa markup tombol/kotak kosong.
func TestHTMLTanpaCTA(t *testing.T) {
	h := Surat{Judul: "J", Teks: "isi"}.HTML(merekUji())
	if strings.Contains(h, "Tombol tidak berfungsi") || strings.Contains(h, "Courier") {
		t.Fatal("markup tombol/kode muncul padahal tidak diminta")
	}
}

func TestTeksPolosMenurunkanParagrafHTML(t *testing.T) {
	s := Surat{Judul: "J", ParagrafHTML: []string{"<p>halo <strong>dunia</strong></p>", "baris &amp; dua"},
		CTALabel: "Buka", CTAURL: "https://x"}
	teks := s.TeksPolos()
	if strings.Contains(teks, "<") || !strings.Contains(teks, "halo dunia") || !strings.Contains(teks, "baris & dua") {
		t.Fatalf("teks polos salah: %q", teks)
	}
	if !strings.Contains(teks, "Buka: https://x") {
		t.Fatalf("CTA tak tercantum di teks polos: %q", teks)
	}
}

// Merek adalah parameter: warna & nama datang dari m, dan garis miring buntut
// basis URL dirapikan supaya tautan kaki tidak berujung "//faq".
func TestMerekDiterapkan(t *testing.T) {
	m := GoNgaji("https://contoh.id/")
	if m.BasisURL != "https://contoh.id" || m.TautanBantuan != "https://contoh.id/faq" {
		t.Fatalf("normalisasi basis URL salah: %+v", m)
	}
	m2 := Merek{Nama: "Toko <X>", WarnaPrimer: "#111111", WarnaPekat: "#222222", WarnaAksen: "#333333", BasisURL: "https://toko.id"}
	h := Surat{Judul: "J", Teks: "isi"}.HTML(m2)
	for _, mau := range []string{"#111111", "#222222", "#333333", "Toko &lt;X&gt;"} {
		if !strings.Contains(h, mau) {
			t.Fatalf("identitas merek tak diterapkan, mencari %q", mau)
		}
	}
	if strings.Contains(h, "Tanya Jawab") {
		t.Fatal("baris bantuan muncul padahal TautanBantuan kosong")
	}
}
