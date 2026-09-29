package utils

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
)

// OpsiOrder mengembalikan urutan opsi asli per posisi tampil untuk satu peserta+soal
// saat acak_opsi aktif: OpsiOrder(...)[i] adalah huruf asli yang tampil di posisi ke-i.
// Wajib identik dengan soalService.generateSeed + rand.Shuffle.
func OpsiOrder(pesertaID, soalID string) []string {
	hash := md5.Sum([]byte(pesertaID + "|" + soalID))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))
	rng := rand.New(rand.NewSource(seed))

	order := []string{"A", "B", "C", "D", "E"}
	rng.Shuffle(len(order), func(i, j int) {
		order[i], order[j] = order[j], order[i]
	})
	return order
}

// JawabanAsli mengonversi huruf yang dipilih peserta (versi acakan di layarnya)
// ke huruf opsi asli di bank soal, sehingga bisa dibandingkan dengan kunci.
func JawabanAsli(jawabanTampil, pesertaID, soalID string) string {
	if jawabanTampil == "" {
		return jawabanTampil
	}
	order := OpsiOrder(pesertaID, soalID)
	idx := int(jawabanTampil[0] - 'A')
	if idx < 0 || idx >= len(order) {
		return jawabanTampil
	}
	return order[idx]
}
