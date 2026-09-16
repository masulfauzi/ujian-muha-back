// Package assets menyimpan file statis yang di-embed langsung ke dalam binary
// (bukan dibaca dari disk saat runtime), supaya tidak tergantung lokasi/working
// directory saat deploy dan tetap tersedia meskipun aplikasi direstart.
package assets

import _ "embed"

// Logo dipakai sebagai logo di kartu ujian (lihat internal/modules/peserta/service).
// Ganti file internal/assets/logo.png dengan logo asli (format PNG) lalu build ulang
// aplikasi agar perubahan logo ikut ter-compile ke dalam binary.
//
//go:embed logo.png
var Logo []byte
